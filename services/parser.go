package services

import (
	"github.com/btcsuite/btcd/rpcclient"
	log "github.com/sirupsen/logrus"
	"ntc-services/config"
	"ntc-services/models"
	"time"
)

type Parser struct {
	BTCClient *rpcclient.Client
	TxRawC    chan ParserTxMsg
}

type ParserTxMsg struct {
	BlockHeight int64
	TxRaw       models.TxRaw
}

func NewParser() (*Parser, error) {
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

	return &Parser{
		BTCClient: client,
		TxRawC:    make(chan ParserTxMsg),
	}, nil
}

func (p *Parser) Run() {
	log.Infof("Running Parser")

	go p.Parse()
	go p.Scan()
}

func (p *Parser) Scan() {
LOOP:
	log.Infof("%+v", STATE)

	STATE, _ = models.GetState()
	latestBlockRaw, err := models.GetLatestBlock()
	if err != nil {
		log.Error(err)
		goto LOOP
	}

	if latestBlockRaw.Height == STATE.ParserBlockHeight {
		time.Sleep(time.Second)
		goto LOOP
	}

	blockRaw, err := models.GetBlockRawByHeight(STATE.ParserBlockHeight)
	if err != nil {
		log.Error(err)
		goto LOOP
	}

	txRaws, err := models.GetTxRawsByBlockID(blockRaw.ID.Hex())
	if err != nil {
		log.Error(err)
		goto LOOP
	}

	log.Info(len(txRaws))

	STATE.ParserBlockHeight = STATE.ParserBlockHeight + 1
	if err := STATE.Update(); err != nil {
		log.Error(err)
	}

	for _, txRaw := range txRaws {
		parserTxMsg := ParserTxMsg{
			BlockHeight: blockRaw.Height,
			TxRaw:       *txRaw,
		}
		p.TxRawC <- parserTxMsg
	}

	goto LOOP
}

func (p *Parser) Parse() {
	semaphore := NewSemaphore(12)
	for {
		select {
		case parserTxMsg := <-p.TxRawC:
			semaphore.Acquire()
			go func() {
				defer semaphore.Release()

				log.Infof("Parsing RawTx.ID: %+v", parserTxMsg.TxRaw.ID.Hex())

				tx := models.NewTx(parserTxMsg.TxRaw.BlockID, parserTxMsg.TxRaw.ID)
				if err := tx.Parse(parserTxMsg.BlockHeight, &parserTxMsg.TxRaw); err != nil {
					log.Error(err)
					return
				}
				if err := tx.Save(); err != nil {
					log.Error(err)
					return
				}

				for i, vinRaw := range parserTxMsg.TxRaw.TxRaw.Vin {
					vin := models.NewVin(tx.ID, parserTxMsg.TxRaw.BlockID, parserTxMsg.TxRaw.ID)
					if err := vin.Parse(int64(i), vinRaw); err != nil {
						log.Error(err)
						return
					}
					if err := vin.Save(); err != nil {
						log.Error(err)
						return
					}
				}

				for i, voutRaw := range parserTxMsg.TxRaw.TxRaw.Vout {
					vout := models.NewVout(tx.ID, parserTxMsg.TxRaw.BlockID, parserTxMsg.TxRaw.ID)
					if err := vout.Parse(int64(i), voutRaw); err != nil {
						log.Error(err)
						return
					}
					if err := vout.Save(); err != nil {
						log.Error(err)
						return
					}
				}
			}()
		}
	}
}
