package services

import (
	"github.com/btcsuite/btcd/rpcclient"
	"ntc-services/config"
)

func NewBitcoinClient() *rpcclient.Client {
	host, _ := config.GetBTCRPCHost()
	//if err != nil {
	//	return nil, err
	//}
	user, _ := config.GetBTCRPCUser()
	//if err != nil {
	//	return nil, err
	//}
	password, _ := config.GetBTCRPCPassword()
	//if err != nil {
	//	return nil, err
	//}

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         *host,
		User:         *user,
		Pass:         *password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	client, _ := rpcclient.New(connCfg, nil)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}

	return client
}
