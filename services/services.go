package services

// var STATE *models.State
var (
	BESTINSLOT     *BestInSlot
	BLOCKCHAIN     *BlockChain
	MEMPOOL        *Mempool
	ORDEX          *Ordex
	BLOCKCHAININFO *BlockChainInfo
	NTCPSBT        *NtcPSBT
)

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

	if err := initClients(); err != nil {
		panic(err)
	}

	return nil
}

func initClients() (err error) {
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
	NTCPSBT, err = NewNtcPSBT()
	if err != nil {
		panic(err)
	}

	return nil
}
