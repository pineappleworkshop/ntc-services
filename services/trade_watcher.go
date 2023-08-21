package services

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil/base58"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http/httptest"
	"ntc-services/config"
	"ntc-services/models"
	"strings"
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
	txHash, err := chainhash.NewHashFromStr("c659c6211c20cfefb622527aca19c2a3b55897c914be7134fa3de757f40ef2dc")
	if err != nil {
		panic(err)
	}

	tx, err := client.GetRawTransactionVerbose(txHash)
	if err != nil {
		panic(err)
	}
	//packet, _, _, err := psbt.NewFromSignedTx(tx.MsgTx())
	//if err != nil {
	//	panic(err)
	//}
	//txFee, err := packet.GetTxFee()
	//if err != nil {
	//	panic(err)
	//}
	//txFee, _ := calculateMinerFeeForPSBT(tx.MsgTx(), 10.7)
	//txJSON, err := json.MarshalIndent(packet, "", "  ")
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("************************")
	fmt.Printf("TX: %+v \n", tx)
	//fmt.Printf("%+v \n", string(txJSON))
	//fmt.Printf("Tx Fee: %+v \n", txFee)
	fmt.Println("************************")

	//b64 := "cHNidP8BAPsCAAAAA9YejvLXC7L0sJf/Y9YDTvQS+G8EqGQM3rqPt/fPUlv8AQAAAAD/////HmRDGSqd+MBIPK4i4mkFBdHXcIJf7bn8jyqtLEJ9t58AAAAAAP////8aj+VYomLXUvBSGNJvHT1v+vvnmBeXoYNB9BikXl0bOAAAAAAA/////wMiAgAAAAAAACJRIDEOiG0coITHkXL3c1CiYLYa7oOw8beKwFhHniA6EqjN8EkCAAAAAAAiUSB66vCgZOWuc47uNVKRLUfEt3NMO6ahGJ1MIRShZtRwdlDDAAAAAAAAF6kUcmXCrsFBJ3nk1RX/qwbzMDvnDu+HAAAAAAABASAiAgAAAAAAABepFMK6ZqsRuDLgNS4CX38IOiqocMDYhwEDBIMAAAABBBYAFEmBMmp37+S2ulB2qjryM6yyy6hPAAEBIEANAwAAAAAAF6kUwrpmqxG4MuA1LgJffwg6KqhwwNiHAQMEgwAAAAEEFgAUSYEyanfv5La6UHaqOvIzrLLLqE8AAQEggMPJAQAAAAAXqRQd9xCNdV6gM+38KQoBI6yaE2J52ocBAwSDAAAAAQQWABQwvaBK8ySDUIiA61CvM32uJeMhIwAAAAA="
	//b64 := "cHNidP8BAPsCAAAAA9YejvLXC7L0sJf/Y9YDTvQS+G8EqGQM3rqPt/fPUlv8AQAAAAD/////HmRDGSqd+MBIPK4i4mkFBdHXcIJf7bn8jyqtLEJ9t58AAAAAAP////8aj+VYomLXUvBSGNJvHT1v+vvnmBeXoYNB9BikXl0bOAAAAAAA/////wMiAgAAAAAAACJRIDEOiG0coITHkXL3c1CiYLYa7oOw8beKwFhHniA6EqjN8EkCAAAAAAAiUSB66vCgZOWuc47uNVKRLUfEt3NMO6ahGJ1MIRShZtRwdlDDAAAAAAAAF6kUcmXCrsFBJ3nk1RX/qwbzMDvnDu+HAAAAAAABASsiAgAAAAAAACJRIFF4WLIywXR+AU/K9Yd1MAZVmt3LkOak9MmIwfVm5Q2YAQMEgwAAAAEXIHrq8KBk5a5zju41UpEtR8S3c0w7pqEYnUwhFKFm1HB2AAEBK0ANAwAAAAAAIlEgUXhYsjLBdH4BT8r1h3UwBlWa3cuQ5qT0yYjB9WblDZgBAwSDAAAAARcgeurwoGTlrnOO7jVSkS1HxLdzTDumoRidTCEUoWbUcHYAAQErgMPJAQAAAAAiUSBReFiyMsF0fgFPyvWHdTAGVZrdy5DmpPTJiMH1ZuUNmAEDBIMAAAABFyB66vCgZOWuc47uNVKRLUfEt3NMO6ahGJ1MIRShZtRwdgAAAAA="
	//b64 := "cHNidP8BAP1xAQIAAAAD1h6O8tcLsvSwl/9j1gNO9BL4bwSoZAzeuo+3989SW/wBAAAAAP////8eZEMZKp34wEg8riLiaQUF0ddwgl/tufyPKq0sQn23nwAAAAAA/////xqP5ViiYtdS8FIY0m8dPW/6++eYF5ehg0H0GKReXRs4AAAAAAD/////BiICAAAAAAAAIlEgMQ6IbRyghMeRcvdzUKJgthrug7Dxt4rAWEeeIDoSqM3wSQIAAAAAACJRIHrq8KBk5a5zju41UpEtR8S3c0w7pqEYnUwhFKFm1HB2UMMAAAAAAAAXqRRyZcKuwUEneeTVFf+rBvMwO+cO74eAlpgAAAAAACJRIHrq8KBk5a5zju41UpEtR8S3c0w7pqEYnUwhFKFm1HB2sGkwAQAAAAAiUSAxDohtHKCEx5Fy93NQomC2Gu6DsPG3isBYR54gOhKozVDDAAAAAAAAF6kUcmXCrsFBJ3nk1RX/qwbzMDvnDu+HAAAAAAABASsiAgAAAAAAACJRIHrq8KBk5a5zju41UpEtR8S3c0w7pqEYnUwhFKFm1HB2AQMEgwAAAAEXICPf2+csXunmh5RunBf2hYnZBVLjemQ12nwFwvH7oh8VAAEBK0ANAwAAAAAAIlEgeurwoGTlrnOO7jVSkS1HxLdzTDumoRidTCEUoWbUcHYBAwSDAAAAARcgI9/b5yxe6eaHlG6cF/aFidkFUuN6ZDXafAXC8fuiHxUAAQErgMPJAQAAAAAiUSB66vCgZOWuc47uNVKRLUfEt3NMO6ahGJ1MIRShZtRwdgEDBIMAAAABFyAj39vnLF7p5oeUbpwX9oWJ2QVS43pkNdp8BcLx+6IfFQAAAAAAAAA="
	b64 := "cHNidP8BAP1xAQIAAAADgomNzXcfxOoBIU/npIhCqhJGgNRI3DV0RWoLYy0gKFYAAAAAAP////+X/hEpdl+blv3XSM9YFOmbNB17VYXU883bXnjqftZZIgEAAAAA/////xqP5ViiYtdS8FIY0m8dPW/6++eYF5ehg0H0GKReXRs4AAAAAAD/////BiICAAAAAAAAIlEgMQ6IbRyghMeRcvdzUKJgthrug7Dxt4rAWEeeIDoSqM2SJAEAAAAAACJRIHrq8KBk5a5zju41UpEtR8S3c0w7pqEYnUwhFKFm1HB29AEAAAAAAAAXqRRyZcKuwUEneeTVFf+rBvMwO+cO74eghgEAAAAAACJRIHrq8KBk5a5zju41UpEtR8S3c0w7pqEYnUwhFKFm1HB2jQIDAAAAAAAiUSAxDohtHKCEx5Fy93NQomC2Gu6DsPG3isBYR54gOhKozfQBAAAAAAAAF6kUcmXCrsFBJ3nk1RX/qwbzMDvnDu+HAAAAAAABASsiAgAAAAAAACJRIHrq8KBk5a5zju41UpEtR8S3c0w7pqEYnUwhFKFm1HB2AQMEgwAAAAEXICPf2+csXunmh5RunBf2hYnZBVLjemQ12nwFwvH7oh8VAAEBK0UvAQAAAAAAIlEgMQ6IbRyghMeRcvdzUKJgthrug7Dxt4rAWEeeIDoSqM0BAwSDAAAAARcgaOsnERGZYk4tTzHE5D5qPFiVTDqjnilfs6PGOnn4u6QAAQEr4JMEAAAAAAAiUSAxDohtHKCEx5Fy93NQomC2Gu6DsPG3isBYR54gOhKozQEDBIMAAAABFyBo6ycREZliTi1PMcTkPmo8WJVMOqOeKV+zo8Y6efi7pAAAAAAAAAA="

	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(decoded)
	pack, err := psbt.NewFromRawBytes(reader, false)
	if err != nil {
		return nil, err
	}

	pk, err := extractTaprootPublicKeyHex(pack.UnsignedTx.TxOut[0].PkScript)
	if err != nil {
		return nil, err
	}

	//23dfdbe72c5ee9e687946e9c17f68589d90552e37a6435da7c05c2f1fba21f15
	//5120310e886d1ca084c79172f77350a260b61aee83b0f1b78ac058479e203a12a8cd

	fmt.Println("+++++++++++++++++++++++")
	fmt.Printf("%+v \n", pack.Inputs)
	fmt.Printf("%+v \n", pk)
	//fmt.Printf("%+v \n", pkHex)
	fmt.Println("+++++++++++++++++++++++")

	decodedAddress, err := btcutil.DecodeAddress("bc1p0t40pgryukh88rhwx4ffzt28cjmhxnpm56s3382vyy22zek5wpmq8rps2l", nil)
	if wa, ok := decodedAddress.(*btcutil.AddressTaproot); ok {
		encodedKey := base64.StdEncoding.EncodeToString(wa.WitnessProgram())
		encodedKeyB58 := base58.Encode(wa.WitnessProgram())
		encodedKeyHex := hex.EncodeToString(wa.WitnessProgram())
		fmt.Println(wa.WitnessVersion())
		fmt.Println(wa.WitnessProgram())
		fmt.Println(encodedKey)
		fmt.Println(encodedKeyB58)
		fmt.Println(encodedKeyHex)
	}

	fmt.Println("-------------------------------------")
	for _, in := range pack.UnsignedTx.TxIn {
		fmt.Printf("in: %+v \n", in)
		//fmt.Printf("in: %+v \n", hex.EncodeToString(in.))
	}
	fmt.Println("-------------------------------------")

	fmt.Println("-------------------------------------")
	for _, out := range pack.UnsignedTx.TxOut {
		fmt.Printf("out: %+v \n", out)
		fmt.Printf("out: %+v \n", hex.EncodeToString(out.PkScript))
	}
	fmt.Println("-------------------------------------")

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
