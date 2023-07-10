package services

import (
	"ntc-services/config"
	"ntc-services/models"
)

func BootstrapState() (*models.State, error) {
	var path string
	env := config.Conf.GetEnv()
	if env != "dev" || env != "prod" {
		path = "./state"
	} else {
		path = "/state"
	}
	state := models.NewState(path)

	if err := state.Read(); err != nil {
		var startingHeight int64
		latestBlock, err := models.GetLatestBlock()
		if err != nil {

		}
		if latestBlock != nil {
			startingHeight = latestBlock.Height
		} else {
			startingHeight = startingHeight
		}
		state.ScannerBlockHeight = startingHeight

		if err := state.Write(); err != nil {
			return nil, err
		}
	}

	return state, nil
}
