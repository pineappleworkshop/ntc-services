package services

import (
	log "github.com/sirupsen/logrus"
	"ntc-services/models"
	"sync"
	"time"
)

type Parser struct {
	TxRawC chan ParserTxMsg
}

type ParserTxMsg struct {
	BlockHeight int64
	TxRaws      []*models.TxRaw
}

func NewParser() (*Parser, error) {
	return &Parser{
		TxRawC: make(chan ParserTxMsg),
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
		goto LOOP
	}

	//for _, txRaw := range txRaws {
	//	parserTxMsg := ParserTxMsg{
	//		BlockHeight: blockRaw.Height,
	//		TxRaw:       *txRaw,
	//	}
	//	p.TxRawC <- parserTxMsg
	//}

	parserTxMsg := ParserTxMsg{
		BlockHeight: blockRaw.Height,
		TxRaws:      txRaws,
	}
	p.TxRawC <- parserTxMsg

	goto LOOP
}

func (p *Parser) Parse() {
	semaphore := NewSemaphore(100)
	for {
		select {
		case parserTxMsg := <-p.TxRawC:
			var txs []models.Tx
			for _, txRaw := range parserTxMsg.TxRaws {
				semaphore.Acquire()
				go func() {
					log.Infof(
						"Parsing RawTx: %+v w/ BlockHeight %+v",
						txRaw,
						parserTxMsg.BlockHeight,
					)

					defer semaphore.Release()

					tx := models.NewTx(txRaw.BlockID, txRaw.ID)
					if err := tx.Parse(parserTxMsg.BlockHeight, txRaw); err != nil {
						log.Error(err)
						return
					}
					txs = append(txs, *tx)

					var vins []models.Vin
					var wgVin sync.WaitGroup
					wgVin.Add(len(txRaw.TxRaw.Vin))
					for i, vinRaw := range txRaw.TxRaw.Vin {
						go func() {
							defer wgVin.Done()
							vin := models.NewVin(tx.ID, txRaw.BlockID, txRaw.ID)
							if err := vin.Parse(int64(i), vinRaw); err != nil {
								log.Error(err)
								return
							}
							vins = append(vins, *vin)
						}()
					}
					wgVin.Wait()
					if err := models.SaveVins(vins); err != nil {
						log.Error(err)
						return
					}

					var vouts []models.Vout
					var wgVout sync.WaitGroup
					wgVout.Add(len(txRaw.TxRaw.Vout))
					for i, voutRaw := range txRaw.TxRaw.Vout {
						go func() {
							defer wgVout.Done()
							vout := models.NewVout(tx.ID, txRaw.BlockID, txRaw.ID)
							if err := vout.Parse(int64(i), voutRaw); err != nil {
								log.Error(err)
								return
							}
							vouts = append(vouts, *vout)
						}()
					}
					wgVout.Wait()
					if err := models.SaveVouts(vouts); err != nil {
						log.Error(err)
						return
					}
				}()
			}

			log.Infof(
				"Parsing RawTxs: %+v w/ BlockHeight %+v",
				parserTxMsg.TxRaws,
				parserTxMsg.BlockHeight,
			)

			if err := models.SaveTxs(txs); err != nil {
				log.Error(err)
				return
			}
		}
	}
}
