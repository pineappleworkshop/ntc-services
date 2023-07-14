package services

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	log "github.com/sirupsen/logrus"
	"ntc-services/config"
)

type Insciber struct {
	BTCClient *rpcclient.Client
}

func Inscribe() (interface{}, error) {
	netParams := &chaincfg.MainNetParams

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

	commitTxOutPointList := make([]*wire.OutPoint, 0)
	// you can get from `client.ListUnspent()`
	unspent, err := client.ListUnspent()
	if err != nil {

	}

	log.Infof("unspent: %+v", unspent)

	utxoAddress := "bc1pxy8gsmgu5zzv0ytj7ae4pgnqkcdwaqas7xmc4szcg70zqwsj4rxss2tvhu"
	address, err := btcutil.DecodeAddress(utxoAddress, netParams)
	if err != nil {
		log.Fatalf("decode address err %v", err)
	}
	unspentList, err := client.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{address})
	if err != nil {
		log.Fatalf("list err err %v", err)
	}

	for i := range unspentList {
		inTxid, err := chainhash.NewHashFromStr(unspentList[i].TxID)
		if err != nil {
			log.Fatalf("decode in hash err %v", err)
		}
		commitTxOutPointList = append(commitTxOutPointList, wire.NewOutPoint(inTxid, unspentList[i].Vout))
	}

	//wallet, err := client.CreateWallet("Inscriber Wallet")
	//if err != nil {
	//	// do something
	//}
	//
	//client.ListUnspentMinMaxAddresses()

	// TODO: add some bitcoin
	// TODO: inscribe something

	return nil, nil
}
