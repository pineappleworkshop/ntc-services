package services

import (
	"ntc-services/models"
)

var STATE *models.State
var SCANNER *Scanner
var PARSER *Parser

func StartServices() (err error) {
	STATE, err = BootstrapState()
	if err != nil {
		panic(err)
	}

	//SCANNER, err = NewScanner()
	//if err != nil {
	//	return err
	//}
	//SCANNER.Run()

	//PARSER, err := NewParser()
	//if err != nil {
	//	return err
	//}
	//PARSER.Run()

	TxWatcher, err := NewTxWatcher()
	if err != nil {
		panic(err)
	}
	go TxWatcher.Run()

	//Inscribe()

	return nil
}
