package services

import "ntc-services/models"

var STATE *models.State
var SCANNER *Scanner

func StartServices() (err error) {
	STATE, err = BootstrapState()
	if err != nil {
		panic(err)
	}

	SCANNER, err = NewScanner()
	if err != nil {
		return err
	}
	SCANNER.Run()

	return nil
}
