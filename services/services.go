package services

import (
	"ntc-services/models"
)

var STATE *models.State

func StartServices() (err error) {
	STATE, err = BootstrapState()
	if err != nil {
		panic(err)
	}

	TxWatcher, err := NewTxWatcher()
	if err != nil {
		panic(err)
	}
	go TxWatcher.Run()

	return nil
}
