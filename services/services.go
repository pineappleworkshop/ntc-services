package services

// var STATE *models.State
var BESTINSLOT *BestInSlot
var MEMPOOL *Mempool

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

	BESTINSLOT, err = NewBestInSlot()
	if err != nil {
		panic(err)
	}

	MEMPOOL, err = NewMempool()
	if err != nil {
		panic(err)
	}

	return nil
}
