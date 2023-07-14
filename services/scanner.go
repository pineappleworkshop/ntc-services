package services

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/config"
	"ntc-services/models"
	"ntc-services/stores"
	"time"
)

type Scanner struct {
	BTCClient      *rpcclient.Client
	BlockHeight    int64
	Txs            chan TxMsg
	CheckingBlocks bool
}

type TxMsg struct {
	TxID       string
	BlockRawID primitive.ObjectID
	Height     int64
	LastTxID   string
}

func NewScanner() (*Scanner, error) {
	host, err := config.GetBTCRPCHost()
	if err != nil {
		return nil, err
	}
	user, err := config.GetBTCRPCUser()
	if err != nil {
		return nil, err
	}
	password, err := config.GetBTCRPCPassword()
	if err != nil {
		return nil, err
	}

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         *host,
		User:         *user,
		Pass:         *password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Scanner{
		BTCClient:      client,
		BlockHeight:    1,
		Txs:            make(chan TxMsg),
		CheckingBlocks: false,
	}, nil
}

func (s *Scanner) Run() {
	//go SCANNER.CheckBlocks()
	go s.ScanTxs()
	go s.ScanBlock()
}

func (s *Scanner) CheckBlocks() {
	s.CheckingBlocks = true
	uncompletedBlockRaws, err := models.GetUncompletedBlockRaws()
	if err != nil {
		log.Error(err)
	}

	// TODO: fix this hack of forcing some concurrency
	scramble(uncompletedBlockRaws)

	for _, br := range uncompletedBlockRaws {
		if time.Now().UTC().Sub(br.CreatedAt) > 10*time.Minute {
			completed, err := models.GetBlockRawCompleted(br.ID.Hex())
			if err != nil {

			}
			if !completed {
				for txHeight, tx := range br.Block.Tx {
					txMsg := TxMsg{
						TxID:       tx,
						BlockRawID: br.ID,
						Height:     int64(txHeight),
						LastTxID:   br.Block.Tx[len(br.Block.Tx)-1],
					}
					s.Txs <- txMsg
				}
			}
		}
	}
	s.CheckingBlocks = false
}

func (s *Scanner) ScanBlock() {
LOOP:
	if s.CheckingBlocks {
		goto LOOP
	}

	STATE, _ = models.GetState()

	log.Infof("State: %+v", STATE)
	blockCount, err := s.BTCClient.GetBlockCount()
	if err != nil {
		log.Error(err)
	}
	if STATE.ScannerBlockHeight == blockCount {
		goto LOOP
	}

	//NEXT_BLOCK:
	started, err := models.GetBlockStarted(STATE.ScannerBlockHeight)
	if started {
		log.Infof("Next Block: %+v", STATE.ScannerBlockHeight)
		STATE.ScannerBlockHeight = STATE.ScannerBlockHeight + 1
		if err := STATE.Update(); err != nil {
			log.Error(err)
		}
		goto LOOP
	}

	blockHash, err := s.BTCClient.GetBlockHash(STATE.ScannerBlockHeight)
	if err != nil {
		log.Error(err)
	}

	blockVerbose, err := s.BTCClient.GetBlockVerbose(blockHash)
	if err != nil {
		log.Error(err)
	}

	log.Infof("Block Height: %+v", blockVerbose.Height)

	blockRaw := models.NewBlockRaw(blockVerbose)
	if err := blockRaw.Save(); err != nil {
		STATE.ScannerBlockHeight = STATE.ScannerBlockHeight + 1
		if err := STATE.Update(); err != nil {
			log.Error(err)
		}
		goto LOOP
	}

	STATE.ScannerBlockHeight = STATE.ScannerBlockHeight + 1
	if err := STATE.Update(); err != nil {
		log.Error(err)
	}

	for txHeight, tx := range blockVerbose.Tx {
		txMsg := TxMsg{
			TxID:       tx,
			BlockRawID: blockRaw.ID,
			Height:     int64(txHeight),
			LastTxID:   blockRaw.Block.Tx[len(blockRaw.Block.Tx)-1],
		}
		s.Txs <- txMsg
	}
	s.BTCClient.GetBalances()

	goto LOOP
}

func (s *Scanner) ScanTxs() {
	semaphore := NewSemaphore(12)
	for {
		var txRaws []*models.TxRaw

		select {
		case txMsg := <-s.Txs:
			semaphore.Acquire()
			go func() {
				defer semaphore.Release()
				txID, err := chainhash.NewHashFromStr(txMsg.TxID)
				if err != nil {
					log.Error(err)
				}

				txRaw, err := s.BTCClient.GetRawTransactionVerbose(txID)
				if err != nil {
					log.Error(err)
				}

				tx := models.NewTxRaw(txMsg.BlockRawID, txMsg.Height, txRaw)
				txRaws = append(txRaws, tx)
				if err := tx.Save(); err != nil {
					log.Error(err)
				}

				if tx.TxID == txMsg.LastTxID {
					blockRaw, err := models.GetBlockRawByID(txMsg.BlockRawID.Hex())
					if err != nil {
						log.Error(err)
						if err.Error() == stores.MONGO_ERR_NOT_FOUND {
							return
						}
					}
					if err := blockRaw.Complete(); err != nil {
						log.Error(err)
					}
					if err := blockRaw.Update(); err != nil {
						log.Error(err)
					}
				}

				log.Infof("Store TxRaw: %+v", tx)
			}()
		}
	}
}
