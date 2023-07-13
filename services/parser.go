package services

import (
	"github.com/btcsuite/btcd/rpcclient"
	log "github.com/sirupsen/logrus"
	"ntc-services/config"
)

type Parser struct {
	BTCClient *rpcclient.Client
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
	}, nil
}

func (p *Parser) Run() {
	log.Infof("Running Parser")
}
