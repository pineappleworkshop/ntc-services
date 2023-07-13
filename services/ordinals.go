package services

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	log "github.com/sirupsen/logrus"
	"ntc-services/config"
	"ntc-services/models"
	"strings"
)

type Ordinal struct {
	BTCClient *rpcclient.Client
}

func NewOrdinal() (*Ordinal, error) {
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

	return &Ordinal{
		BTCClient: client,
	}, nil
}

func (o *Ordinal) Run() {
	log.Infof("Running Parser")
	blockCount, err := o.BTCClient.GetBlockCount()
	if err != nil {
		log.Error(err)
	}
	for i := FIRST_BLOCK_HEIGHT; i <= blockCount; i++ {
		log.Info("Get Block at height: %+v", i)
		blockRaw, err := models.GetBlockRawByHeight(i)
		if err != nil {
		}

		txRaws, err := models.GetTxRawsByBlockID(blockRaw.ID.Hex())
		if err != nil {
		}

		log.Info("Get TxRaws at height: %+v", blockRaw)
		for _, txRaw := range txRaws {
			if txRaw.TxRaw.Vin[0].HasWitness() {
				for _, witness := range txRaw.TxRaw.Vin[0].Witness {
					data, err := hex.DecodeString(witness)
					if err != nil {
						fmt.Println(err)
						return
					}

					decodedString := string(data)
					if strings.ContainsAny(decodedString, "ord") {
						fmt.Println("Contains 'ord'")
					}
				}
			}
		}
	}
}
