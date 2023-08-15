package services

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"ntc-services/config"
	"ntc-services/models"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	log "github.com/sirupsen/logrus"
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

type TradeWatcher struct {
	BTCClient *rpcclient.Client
	Trades    map[string]*models.Trade
}

func NewTradeWatcher() (*TradeWatcher, error) {
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

	// TODO: Remove commented code
	txHash, err := chainhash.NewHashFromStr("57786ad9ae59f3afa6be5603248fc22c88b9d5a53685a72ec84750b47eacba5e")
	if err != nil {
		panic(err)
	}

	tx, err := client.GetRawTransaction(txHash)
	if err != nil {
		panic(err)
	}
	packet, _, _, err := psbt.NewFromSignedTx(tx.MsgTx())
	if err != nil {
		panic(err)
	}
	//txFee, err := packet.GetTxFee()
	//if err != nil {
	//	panic(err)
	//}
	txFee, _ := calculateMinerFeeForPSBT(tx.MsgTx(), 10.7)

	txJSON, err := json.MarshalIndent(packet, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("************************")
	fmt.Printf("%+v \n", string(txJSON))
	fmt.Printf("Tx Fee: %+v \n", txFee)
	fmt.Println("************************")

	return &TradeWatcher{
		BTCClient: client,
		Trades:    make(map[string]*models.Trade),
	}, nil
}

func (tw *TradeWatcher) Run() {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	go tw.Poll()

	for {
		for txId, t := range tw.Trades {
			log.Infof("Watching Trade: %+v", txId)

			trade, err := models.GetTradeByID(c, t.ID.Hex())
			if err != nil {
				log.Error(err)
				continue
			}

			if trade.Status == "SUBMITTED" {
				mempoolEntry, err := tw.BTCClient.GetMempoolEntry(txId)
				if err != nil {
					log.Error(err)
				}
				if mempoolEntry != nil {
					log.Infof("Trade in Mempool: %+v", txId)
					trade.Status = "MEMPOOL"
					if err := trade.Update(c); err != nil {
						log.Error(err)
					}
					continue
				}
			}

			if trade.Status == "MEMPOOL" {
				if _, err := tw.BTCClient.GetMempoolEntry(txId); err != nil {
					if err.Error() == ERR_NOT_IN_MEMPOOL {
						hash, err := chainhash.NewHashFromStr(txId)
						if err != nil {
							log.Error(err)
						}
						tx, err := tw.BTCClient.GetRawTransactionVerbose(hash)
						if err != nil {
							log.Error(err)
						}

						if tx != nil {
							if tx.Confirmations >= 1 {
								log.Infof("Trade Confirmed: %+v", txId)

								trade.Status = "CONFIRMED"
								if err := trade.Update(c); err != nil {
									log.Error(err)
								}
								delete(tw.Trades, txId)
							} else {
								log.Infof("Trade Failed: %+v", txId)

								trade.Status = "FAILED"
								if err := trade.Update(c); err != nil {
									log.Error(err)
								}
								delete(tw.Trades, txId)
							}
						} else {
							log.Infof("Trade Failed: %+v", txId)

							trade.Status = "FAILED"
							if err := trade.Update(c); err != nil {
								log.Error(err)
							}
							delete(tw.Trades, txId)
						}
					}
				}
			}
		}

		time.Sleep(time.Second * 5)
	}
}

func (tw *TradeWatcher) Poll() {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for {
		trades, err := models.GetTradesByStatus(c, "SUBMITTED")
		if err != nil {
			log.Error(err)
		}
		for _, trade := range trades {
			if trade.TxID != nil {
				tw.Trades[*trade.TxID] = trade
			}
		}

		time.Sleep(time.Second * 5)
	}
}

//func (tw *TradeWatcher) Watch() {
//	log.Infof("Tx Watcher Running")
//
//	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)
//	pipeline := mongo.Pipeline{
//		{
//			{"$match", bson.D{
//				{"fullDocument.txId", bson.D{{"$exists", true}}},
//				{"operationType", bson.M{"$in": bson.A{"insert", "update", "replace", "upsert"}}},
//			}},
//		},
//	}
//	changeStream, err := collection.Watch(context.Background(), pipeline)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer changeStream.Close(context.Background())
//
//	for changeStream.Next(context.Background()) {
//
//		fmt.Println("===================")
//		fmt.Println("Detected")
//		fmt.Println("===================")
//
//		var changeDoc struct {
//			FullDocument models.Trade `bson:"fullDocument"`
//		}
//		if err := changeStream.Decode(&changeDoc); err != nil {
//			log.Println("Error decoding change document:", err)
//			continue
//		}
//		tw.Trades[changeDoc.FullDocument.TxID] = changeDoc.FullDocument
//	}
//}
