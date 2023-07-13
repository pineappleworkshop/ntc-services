package services

import (
	"ntc-services/models"
)

func BootstrapState() (*models.State, error) {
	//var path string
	//
	//// TODO: fix, not sensing dev in cluster
	//env := config.Conf.GetEnv()
	//if env != "dev" || env != "prod" {
	//	path = STATE_PATH
	//} else {
	//	path = STATE_PATH
	//}
	//state := models.NewState(path)
	//
	//if err := state.Read(); err != nil {
	//	var startingHeight int64
	//	latestBlock, err := models.GetLatestBlock()
	//	if err != nil {
	//
	//	}
	//	if latestBlock != nil {
	//		startingHeight = latestBlock.Height
	//	} else {
	//		startingHeight = startingHeight
	//	}
	//	state.ScannerBlockHeight = startingHeight
	//
	//	if err := state.Write(); err != nil {
	//		return nil, err
	//	}
	//}

	state, err := models.GetState()
	if err != nil {
		return nil, err
	}

	return state, nil
}
