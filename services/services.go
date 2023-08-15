package services

// var STATE *models.State
var BESTINSLOT *BestInSlot
var BLOCKCHAIN *BlockChain
var MEMPOOL *Mempool
var ORDEX *Ordex
var BLOCKCHAININFO *BlockChainInfo

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

	BLOCKCHAIN, err = NewBlockChain()
	if err != nil {
		panic(err)
	}

	MEMPOOL, err = NewMempool()
	if err != nil {
		panic(err)
	}

	ORDEX, err = NewOrdex()
	if err != nil {
		panic(err)
	}

	BLOCKCHAININFO, err = NewBlockChainInfo()
	if err != nil {
		panic(err)
	}

	return nil
}
