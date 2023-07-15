package services

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ntc-services/config"
	"ntc-services/models"
	"ntc-services/stores"
	"time"
)

/*
order statuses
	export const PENDING = "PENDING"; 		// Not used, intention is to have a draft order state for the maker
	export const CREATED = "CREATED"; 		// Maker has created and signed a trade (PSBT), and is ready to be reviewed,
												rejected, or confirmed by the taker
	export const SUBMITTED = "SUBMITTED";	// Taker has accepted, signed, and broadcasted the trade
												(has been sent and is in the mempool)
	export const REJECTED = "REJECTED";		// The taker rejects the trade from the maker
	export const FAILED = "FAILED";			// The tx has left the mempool but failed to write to the blockchain
	export const CONFIRMED = "CONFIRMED"	// The tx has left the mempool and has been successful facilitated
*/

type TxWatcher struct {
	BTCClient *rpcclient.Client
	Orders    map[string]models.Order
}

func NewTxWatcher() (*TxWatcher, error) {
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

	return &TxWatcher{
		BTCClient: client,
		Orders:    make(map[string]models.Order),
	}, nil
}

func (tw *TxWatcher) Run() {
	go tw.Watch()

	for {
		for txId, order := range tw.Orders {
			fmt.Println(fmt.Sprintf("Trades: %+v", tw.Orders))

			mempoolEntry, err := tw.BTCClient.GetMempoolEntry(txId)
			if err != nil {
				log.Error(err)
			}

			if mempoolEntry != nil {
				if order.Status == "SUBMITTED" {
					order.Status = "MEMPOOL"
					if err := order.Update(); err != nil {
						log.Error(err)
					}
				}
			} else {
				if order.Status == "MEMPOOL" {
					hash, err := chainhash.NewHashFromStr(txId)
					if err != nil {
						log.Error(err)
					}
					tx, err := tw.BTCClient.GetTransaction(hash)
					if err != nil {
						log.Error(err)
					}

					if tx.Confirmations >= 10 {
						order.Status = "CONFIRMED"
						if err := order.Update(); err != nil {
							log.Error(err)
						}
						delete(tw.Orders, txId)
					} else {
						// TODO: TX did not get written to block
					}
				}
			}

			//hash, err := chainhash.NewHashFromStr(order.TxID)
			//tx, err := tw.BTCClient.GetRawTransactionVerbose(hash)
			//if err != nil {
			//}
			//fmt.Println(fmt.Sprintf("Tx: %+v", tx))
		}

		time.Sleep(time.Second * 5)
	}
}

func (tw *TxWatcher) Watch() {
	log.Infof("Tx Watcher Running")

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"fullDocument.txId", bson.D{{"$exists", true}}},
				{"operationType", bson.M{"$in": bson.A{"insert", "update", "replace", "upsert"}}},
			}},
		},
	}
	changeStream, err := collection.Watch(context.Background(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(context.Background())

	for changeStream.Next(context.Background()) {
		var changeDoc struct {
			FullDocument models.Order `bson:"fullDocument"`
		}
		if err := changeStream.Decode(&changeDoc); err != nil {
			log.Println("Error decoding change document:", err)
			continue
		}
		tw.Orders[changeDoc.FullDocument.TxID] = changeDoc.FullDocument
	}
}
