package services

// var STATE *models.State

func StartServices() (err error) {
	// STATE, err = BootstrapState()
	// if err != nil {
	// 	panic(err)
	// }

	TradeWatcher, err := NewTradeWatcher()
	if err != nil {
		panic(err)
	}
	go TradeWatcher.Run()

	return nil
}
