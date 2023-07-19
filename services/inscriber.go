package services

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	log "github.com/sirupsen/logrus"
	"ntc-services/config"
)

type Inscriber struct {
	BTCClient *rpcclient.Client
}

func NewInscriber() (*Inscriber, error) {
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

	return &Inscriber{
		BTCClient: client,
	}, nil
}

func (i *Inscriber) Inscribe(request *InscriptionRequest) (*wire.MsgTx, []*wire.MsgTx, int64, error) {
	netParams := &chaincfg.MainNetParams

	host, err := config.GetBTCRPCHost()
	if err != nil {
		return nil, nil, -1, err
	}
	user, err := config.GetBTCRPCUser()
	if err != nil {
		return nil, nil, -1, err
	}
	password, err := config.GetBTCRPCPassword()
	if err != nil {
		return nil, nil, -1, err
	}

	// Connect to local bitcoin core RPC server using HTTP POST mode.z
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
		return nil, nil, -1, err
	}

	tool, err := NewInscriptionTool(netParams, client, request)
	if err != nil {
		log.Fatalf("Failed to create inscription tool: %v", err)
	}

	commitTx, revealTxs, fees, err := tool.BuildInscribe()
	if err != nil {
		return nil, nil, -1, err
	}

	return commitTx, revealTxs, fees, nil
}
