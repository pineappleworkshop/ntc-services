package services

var SCANNER *Scanner

func StartServices() (err error) {
	SCANNER, err = NewScanner()
	if err != nil {
		return err
	}

	SCANNER.CheckBlocks()
	go SCANNER.ScanTxs()
	go SCANNER.ScanBlock()

	return nil
}
