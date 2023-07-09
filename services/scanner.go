package services

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/config"
	"ntc-services/models"
	"sync"
	"time"
)

/*
We need to implement something like this:
https://unisat.io/wallet-api-v4/address/inscriptions
	?address=bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd&cursor=0&size=10

We need to essentially dump redb to mongodb.

{
  "list": [
    {
      "inscriptionId": "4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2fi0",
      "inscriptionNumber": 11822188,
      "address": "bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd",
      "outputValue": 546,
      "preview": "https://ordinals.com/preview/4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2fi0",
      "content": "https://ordinals.com/content/4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2fi0",
      "contentLength": 57,
      "contentType": "text/plain;charset=utf-8",
      "contentBody": "",
      "timestamp": 1686720452,
      "genesisTransaction": "4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2f",
      "location": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:8:0",
      "output": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:8",
      "offset": 0
    },
    {
      "inscriptionId": "c0e627075a7991e5049230c886e841c9eb82f2b1dc392ec86acd706c25d72afdi0",
      "inscriptionNumber": 14938723,
      "address": "bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd",
      "outputValue": 546,
      "preview": "https://ordinals.com/preview/c0e627075a7991e5049230c886e841c9eb82f2b1dc392ec86acd706c25d72afdi0",
      "content": "https://ordinals.com/content/c0e627075a7991e5049230c886e841c9eb82f2b1dc392ec86acd706c25d72afdi0",
      "contentLength": 4,
      "contentType": "text/plain;charset=utf-8",
      "contentBody": "",
      "timestamp": 1688653182,
      "genesisTransaction": "c0e627075a7991e5049230c886e841c9eb82f2b1dc392ec86acd706c25d72afd",
      "location": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:5:0",
      "output": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:5",
      "offset": 0
    },
    {
      "inscriptionId": "720dff8e5224c5f974918cab07fcbe1d820f6cca2cbc3d5bfff22cd7eb76eb7ci0",
      "inscriptionNumber": 14938719,
      "address": "bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd",
      "outputValue": 546,
      "preview": "https://ordinals.com/preview/720dff8e5224c5f974918cab07fcbe1d820f6cca2cbc3d5bfff22cd7eb76eb7ci0",
      "content": "https://ordinals.com/content/720dff8e5224c5f974918cab07fcbe1d820f6cca2cbc3d5bfff22cd7eb76eb7ci0",
      "contentLength": 18,
      "contentType": "text/plain;charset=utf-8",
      "contentBody": "",
      "timestamp": 1688653182,
      "genesisTransaction": "720dff8e5224c5f974918cab07fcbe1d820f6cca2cbc3d5bfff22cd7eb76eb7c",
      "location": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:6:0",
      "output": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:6",
      "offset": 0
    },
    {
      "inscriptionId": "88665e98d24676cb2268551bda756dbfe79c0bb0706812fc4c6ebb5cdf31cf1ai0",
      "inscriptionNumber": 14938714,
      "address": "bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd",
      "outputValue": 546,
      "preview": "https://ordinals.com/preview/88665e98d24676cb2268551bda756dbfe79c0bb0706812fc4c6ebb5cdf31cf1ai0",
      "content": "https://ordinals.com/content/88665e98d24676cb2268551bda756dbfe79c0bb0706812fc4c6ebb5cdf31cf1ai0",
      "contentLength": 36,
      "contentType": "text/plain;charset=utf-8",
      "contentBody": "",
      "timestamp": 1688653182,
      "genesisTransaction": "88665e98d24676cb2268551bda756dbfe79c0bb0706812fc4c6ebb5cdf31cf1a",
      "location": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:7:0",
      "output": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:7",
      "offset": 0
    }
  ],
  "total": 4
}
*/

type Scanner struct {
	BTCClient   *rpcclient.Client
	WaitGroup   *sync.WaitGroup
	BlockHeight int64
	Txs         chan TxMsg
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
		BTCClient:   client,
		BlockHeight: 1,
		Txs:         make(chan TxMsg),
	}, nil
}

func (s *Scanner) CheckBlocks() {
	uncompletedBlockRaws, err := models.GetUncompletedBlockRaws()
	if err != nil {
		log.Error(err)
	}

	for _, br := range uncompletedBlockRaws {
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

func (s *Scanner) ScanBlock() {
LOOP:
	time.Sleep(time.Second * 5)

	var startingHeight int64
	latestBlock, err := models.GetLatestBlock()
	if latestBlock != nil {
		startingHeight = latestBlock.Height
	} else {
		startingHeight = FIRST_BLOCK_HEIGHT
	}

	blockCount, err := s.BTCClient.GetBlockCount()
	if err != nil {
		log.Error(err)
	}

	if startingHeight == blockCount {
		goto LOOP
	}

	for height := startingHeight; height < blockCount; height++ {
		blockHash, err := s.BTCClient.GetBlockHash(height)
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
			log.Error(err)
		}

		go func() {
			for txHeight, tx := range blockVerbose.Tx {
				txMsg := TxMsg{
					TxID:       tx,
					BlockRawID: blockRaw.ID,
					Height:     int64(txHeight),
					LastTxID:   blockRaw.Block.Tx[len(blockRaw.Block.Tx)-1],
				}
				s.Txs <- txMsg
			}
		}()
		time.Sleep(time.Second * 20)
	}
	goto LOOP
}

func (s *Scanner) ScanTxs() {
	semaphore := NewSemaphore(12)
	for {
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
				if err := tx.Save(); err != nil {
					log.Error(err)
				}

				if tx.TxID == txMsg.LastTxID {
					blockRaw, err := models.GetBlockRawByID(txMsg.BlockRawID.Hex())
					if err != nil {
						log.Error(err)
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
