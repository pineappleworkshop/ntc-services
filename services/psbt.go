package services

import (
	"errors"
	"github.com/btcsuite/btcd/wire"
	"math"
	"ntc-services/models"
	"regexp"
	"strconv"
	"strings"
)

type PSBT struct {
	Trade                 *models.Trade
	MakerUTXOsAll         []*models.UTXO        // All maker UTXOs
	MakerInscriptionsAll  []*models.Inscription // All maker inscriptions
	TakerUTXOsAll         []*models.UTXO        // All taker UTXOs
	TakerInscriptionsAll  []*models.Inscription // All maker inscriptions
	MakerInscriptionUTXOs []*models.UTXO        // Maker inscription UTXOs for trade
	MakerPaymentUTXOs     []*models.UTXO        // Maker other UTXOs for trade
	TakerInscriptionUTXOs []*models.UTXO        // Taker inscription UTXOs for trade
	TakerPaymentUTXOs     []*models.UTXO        // Taker other UTXOs for trade
	MakerPayment          int64
	TakerPayment          int64
	MakerChange           int64
	TakerChange           int64
	PlatformFee           int64
	Inputs                map[int]*Input
	Outputs               map[int]*Output
}

type PSBTReqBody struct {
	Inputs  map[int]*Input  `json:"inputs"`
	Outputs map[int]*Output `json:"outputs"`
}

type Input struct {
	SenderAddr      string `json:"sender_addr"`
	SenderPublicKey string `json:"sender_public_key"`
	Type            string `json:"type"`
	TxID            string `json:"tx_id"`
	Index           int64  `json:"index"`
	WitnessUTXO     struct {
		Script string `json:"script"`
		Amount int64  `json:"amount"`
	} `json:"witness_utxo"`
}

type Output struct {
	RecipientAddr string `json:"recipient_addr"`
	Value         int64  `json:"value"`
}

func NewPBST(
	trade *models.Trade,
	makerUTXOsAll, takerUTXOsAll []*models.UTXO,
	makerInscriptionsAll, takerInscriptionsAll []*models.Inscription,
) *PSBT {
	return &PSBT{
		Trade:                 trade,
		MakerUTXOsAll:         makerUTXOsAll,
		MakerInscriptionsAll:  makerInscriptionsAll,
		TakerUTXOsAll:         takerUTXOsAll,
		TakerInscriptionsAll:  takerInscriptionsAll,
		MakerInscriptionUTXOs: []*models.UTXO{},
		MakerPaymentUTXOs:     []*models.UTXO{},
		TakerInscriptionUTXOs: []*models.UTXO{},
		TakerPaymentUTXOs:     []*models.UTXO{},
		Inputs:                map[int]*Input{},
		Outputs:               map[int]*Output{},
	}
}

func (p *PSBT) Create() error {
	if err := p.selectInscriptionsUTXOs(); err != nil {
		return err
	}
	if err := p.calculatePlatformFee(); err != nil {
		return err
	}
	if err := p.selectPaymentUTXOs(); err != nil {
		return err
	}
	if err := p.createInscriptionInputs(); err != nil {
		return err
	}
	if err := p.createInscriptionOutputs(); err != nil {
		return err
	}
	if err := p.createPaymentInputs(); err != nil {
		return err
	}
	if err := p.createPaymentsOutputs(); err != nil {
		return err
	}

	return nil
}

func (p *PSBT) ToReq() *PSBTReqBody {
	return &PSBTReqBody{
		Inputs:  p.Inputs,
		Outputs: p.Outputs,
	}
}

func (p *PSBT) selectInscriptionsUTXOs() error {
	// MAKER: parse inscription UTXOs from all UTXOs
	for i, utxo := range p.MakerUTXOsAll {
		// Find matching UTXOs for inscriptions in trade and add
		for _, inscription := range p.Trade.Maker.Inscriptions {
			// Ensure inscriptionID has txID & index
			inscriptionIdS := strings.Split(inscription.InscriptionID, "i")
			if len(inscriptionIdS) != 2 {
				// TODO: log
				return errors.New("InscriptionID is not in the right format")
			}

			// Ensure txID is valid
			if valid := isValidTxID(inscriptionIdS[0]); !valid {
				// TODO: log
				return errors.New("Inscription TxID is not valid")
			}

			// Ensure index in int
			if _, err := strconv.Atoi(inscriptionIdS[1]); err != nil {
				// TODO: log
				return errors.New("Inscription Index is not an integer")
			}

			// Parse inscription utxos and other utxos and remove for original list
			if utxo.TxHashBigEndian == inscriptionIdS[0] {
				p.MakerInscriptionUTXOs = append(p.MakerInscriptionUTXOs, utxo)
				p.MakerUTXOsAll = append(p.MakerUTXOsAll[:i], p.MakerUTXOsAll[i+1:]...)

				// Remove inscription from all inscriptions
				for ii, inscriptionX := range p.MakerInscriptionsAll {
					if inscription.InscriptionID == inscriptionX.InscriptionID {
						p.MakerInscriptionsAll = append(p.MakerInscriptionsAll[:ii], p.MakerInscriptionsAll[ii+1:]...)
					}
				}
			}
		}
		// Remove the rest of inscription utxos
		for _, inscription := range p.MakerInscriptionsAll {
			// Ensure inscriptionID has txID & index
			inscriptionIdS := strings.Split(inscription.InscriptionID, "i")
			if len(inscriptionIdS) != 2 {
				// TODO: log
				return errors.New("InscriptionID is not in the right format")
			}

			// Ensure txID is valid
			if valid := isValidTxID(inscriptionIdS[0]); !valid {
				// TODO: log
				return errors.New("Inscription TxID is not valid")
			}

			// Ensure index in int
			if _, err := strconv.Atoi(inscriptionIdS[1]); err != nil {
				// TODO: log
				return errors.New("Inscription Index is not an integer")
			}

			// Remove inscription for original list
			if utxo.TxHashBigEndian == inscriptionIdS[0] {
				p.MakerUTXOsAll = append(p.MakerUTXOsAll[:i], p.MakerUTXOsAll[i+1:]...)
			}
		}
	}

	// TAKER: parse inscription UTXOs from all UTXOs
	for i, utxo := range p.TakerUTXOsAll {
		// Find matching UTXOs for inscriptions in trade and add
		for _, inscription := range p.Trade.Taker.Inscriptions {
			// Ensure inscriptionID has txID & index
			inscriptionIdS := strings.Split(inscription.InscriptionID, "i")
			if len(inscriptionIdS) != 2 {
				// TODO: log
				return errors.New("InscriptionID is not in the right format")
			}

			// Ensure txID is valid
			if valid := isValidTxID(inscriptionIdS[0]); !valid {
				// TODO: log
				return errors.New("Inscription TxID is not valid")
			}

			// Ensure index in int
			if _, err := strconv.Atoi(inscriptionIdS[1]); err != nil {
				// TODO: log
				return errors.New("Inscription Index is not an integer")
			}

			// Add inscription to psbt inscription utxos and remove for original list
			if utxo.TxHashBigEndian == inscriptionIdS[0] {
				p.TakerInscriptionUTXOs = append(p.TakerInscriptionUTXOs, utxo)
				p.TakerUTXOsAll = append(p.TakerUTXOsAll[:i], p.TakerUTXOsAll[i+1:]...)

				// Remove inscription from all inscriptions
				for ii, inscriptionX := range p.TakerInscriptionsAll {
					if inscription.InscriptionID == inscriptionX.InscriptionID {
						p.TakerInscriptionsAll = append(p.TakerInscriptionsAll[:ii], p.TakerInscriptionsAll[ii+1:]...)
					}
				}
			}
		}
		// Remove the rest of inscription utxos
		for _, inscription := range p.TakerInscriptionsAll {
			// Ensure inscriptionID has txID & index
			inscriptionIdS := strings.Split(inscription.InscriptionID, "i")
			if len(inscriptionIdS) != 2 {
				// TODO: log
				return errors.New("InscriptionID is not in the right format")
			}

			// Ensure txID is valid
			if valid := isValidTxID(inscriptionIdS[0]); !valid {
				// TODO: log
				return errors.New("Inscription TxID is not valid")
			}

			// Ensure index in int
			if _, err := strconv.Atoi(inscriptionIdS[1]); err != nil {
				// TODO: log
				return errors.New("Inscription Index is not an integer")
			}

			// Remove inscription for original list
			if utxo.TxHashBigEndian == inscriptionIdS[0] {
				p.TakerUTXOsAll = append(p.TakerUTXOsAll[:i], p.TakerUTXOsAll[i+1:]...)
			}
		}
	}

	return nil
}

func (p *PSBT) calculatePlatformFee() error {
	// Calculate total fees and split between parties
	totalPayments := p.Trade.Maker.BTC + p.Trade.Taker.BTC

	// Calculate inscriptions total floor price
	totalInscriptionFloorPrices := int64(0)
	// TODO: think about, should maker floor prices be calculated or just the taker
	//for _, inscription := range p.Trade.Maker.Inscriptions {
	//	if inscription.FloorPrice != nil {
	//		totalInscriptionFloorPrices = totalInscriptionFloorPrices + *inscription.FloorPrice
	//	} else {
	//		totalInscriptionFloorPrices = totalInscriptionFloorPrices + int64(10000)
	//	}
	//}
	for _, inscription := range p.Trade.Taker.Inscriptions {
		if inscription.FloorPrice != nil {
			totalInscriptionFloorPrices = totalInscriptionFloorPrices + *inscription.FloorPrice
		} else {
			totalInscriptionFloorPrices = totalInscriptionFloorPrices + int64(10000)
		}
	}
	p.PlatformFee = totalPayments/100 + totalInscriptionFloorPrices/100

	return nil
}

func (p *PSBT) selectPaymentUTXOs() error {
	// MAKER: Select the payment UTXOs for the psbt
	// Calculate total payment and set control vars
	p.MakerPayment = p.Trade.Maker.BTC + p.PlatformFee/2
	makerUTXOsTotalValue := int64(0)
	makerEnoughToPay := false
	// Iterate through payment UTXOs and use UTXOs until payment is greater than needed
	for i, utxo := range p.MakerUTXOsAll {
		makerUTXOsTotalValue = makerUTXOsTotalValue + utxo.Value
		p.MakerPaymentUTXOs = append(p.MakerPaymentUTXOs, utxo)
		p.MakerUTXOsAll = append(p.MakerUTXOsAll[:i], p.MakerUTXOsAll[i+1:]...)

		// Check to see if payment is greater than needed
		if makerUTXOsTotalValue >= p.MakerPayment {
			makerEnoughToPay = true
			break
		}
	}
	if !makerEnoughToPay {
		return errors.New("Maker does not have enough BTC to pay")
	}
	// calculate maker change
	p.MakerChange = makerUTXOsTotalValue - p.MakerPayment

	// TAKER: Select the payment UTXOs for the psbt
	// Calculate total payment and set control vars
	p.TakerPayment = p.Trade.Taker.BTC + p.PlatformFee/2
	takerUTXOsTotalValue := int64(0)
	takerEnoughToPay := false
	// Iterate through payment UTXOs and use UTXOs until payment is greater than needed
	for i, utxo := range p.TakerUTXOsAll {
		takerUTXOsTotalValue = takerUTXOsTotalValue + utxo.Value
		p.TakerPaymentUTXOs = append(p.TakerPaymentUTXOs, utxo)
		p.TakerUTXOsAll = append(p.TakerUTXOsAll[:i], p.TakerUTXOsAll[i+1:]...)

		// Check to see if payment is greater than needed
		if takerUTXOsTotalValue >= p.TakerPayment {
			takerEnoughToPay = true
			break
		}
	}
	if !takerEnoughToPay {
		return errors.New("Taker does not have enough BTC to pay")
	}
	// calculate taker change
	p.TakerChange = takerUTXOsTotalValue - p.TakerPayment

	return nil
}

//type Input struct {
//	SenderAddr      string `json:"sender_addr"`
//	SenderPublicKey string `json:"sender_public_key"`
//	Type            string `json:"type"`
//	TxID            string `json:"tx_id"`
//	Index           int64  `json:"index"`
//	WitnessUTXO     struct {
//		Script string `json:"script"`
//		Amount int64  `json:"amount"`
//	} `json:"witness_utxo"`
//}

func (p *PSBT) createInscriptionInputs() error {
	// MAKER: Create the inscription Inputs
	// Iterate through the maker inscriptions and create the proper outputs
	for i, utxo := range p.MakerInscriptionUTXOs {
		// Create input
		input := new(Input)
		input.SenderAddr = p.Trade.Maker.Wallet.TapRootAddr
		input.SenderPublicKey = p.Trade.Maker.Wallet.TapRootPublicKey
		input.Type = "taproot" // TODO: detect from wallet type
		input.TxID = utxo.TxHash
		input.Index = utxo.TxOutputN
		input.WitnessUTXO.Amount = utxo.Value
		input.WitnessUTXO.Script = utxo.Script

		// Append to psbt inputs
		p.Inputs[i] = input
	}

	// TAKER: Create the inscription Inputs
	// Iterate through the taker inscriptions and create the proper outputs
	for i, utxo := range p.TakerInscriptionUTXOs {
		// Create input
		input := new(Input)
		input.SenderAddr = p.Trade.Taker.Wallet.TapRootAddr
		input.SenderPublicKey = p.Trade.Taker.Wallet.TapRootPublicKey
		input.Type = "taproot" // TODO: detect from wallet type
		input.TxID = utxo.TxHash
		input.Index = utxo.TxOutputN
		input.WitnessUTXO.Amount = utxo.Value
		input.WitnessUTXO.Script = utxo.Script

		// Append to psbt inputs
		p.Inputs[i] = input
	}

	return nil
}

//type Output struct {
//	RecipientAddr      string `json:"recipient_addr"`
//	Value              int64  `json:"value"`
//}

func (p *PSBT) createInscriptionOutputs() error {
	// MAKER: Create the inscription outputs
	// Iterate through the maker inscriptions UTXOs and create outputs to the recipient
	for i, utxo := range p.MakerInscriptionUTXOs {
		// Create output
		output := new(Output)
		output.RecipientAddr = p.Trade.Taker.Wallet.TapRootAddr
		output.Value = utxo.Value

		// append to outputs
		p.Outputs[i] = output
	}

	// TAKER: Create the inscription outputs
	// Iterate through the taker inscriptions UTXOs and create outputs to the recipient
	for i, utxo := range p.TakerInscriptionUTXOs {
		// Create output
		output := new(Output)
		output.RecipientAddr = p.Trade.Maker.Wallet.TapRootAddr
		output.Value = utxo.Value

		// append to outputs
		p.Outputs[i] = output
	}

	return nil
}

//type Input struct {
//	SenderAddr      string `json:"sender_addr"`
//	SenderPublicKey string `json:"sender_public_key"`
//	Type            string `json:"type"`
//	TxID            string `json:"tx_id"`
//	Index           int64  `json:"index"`
//	WitnessUTXO     struct {
//		Script string `json:"script"`
//		Amount int64  `json:"amount"`
//	} `json:"witness_utxo"`
//}

func (p *PSBT) createPaymentInputs() error {
	// MAKER: Create the inscription Inputs
	// Iterate through maker payment utxos and create input
	for i, utxo := range p.MakerPaymentUTXOs {
		// Create input
		input := new(Input)
		if p.Trade.Maker.Wallet.Type == "unisat" {
			input.SenderAddr = p.Trade.Maker.Wallet.TapRootAddr
			input.SenderPublicKey = p.Trade.Maker.Wallet.TapRootPublicKey
			input.Type = "taproot" // TODO: detect from wallet type
		} else if p.Trade.Maker.Wallet.Type == "xverse" || p.Trade.Maker.Wallet.Type == "hiro" {
			input.SenderAddr = p.Trade.Maker.Wallet.TapRootAddr
			input.SenderPublicKey = p.Trade.Maker.Wallet.TapRootPublicKey
			input.Type = "segwit"
		} else {
			// TODO: all logging
			return errors.New("Invalid maker wallet type")
		}
		input.TxID = utxo.TxHash
		input.Index = utxo.TxOutputN
		input.WitnessUTXO.Amount = utxo.Value
		input.WitnessUTXO.Script = utxo.Script

		// Append to psbt inputs
		p.Inputs[len(p.Inputs)+i] = input
	}

	// TAKER: Create the inscription Inputs
	// Iterate through taker payment utxos and create input
	for i, utxo := range p.TakerPaymentUTXOs {
		// Create input
		input := new(Input)
		if p.Trade.Taker.Wallet.Type == "unisat" {
			input.SenderAddr = p.Trade.Taker.Wallet.TapRootAddr
			input.SenderPublicKey = p.Trade.Taker.Wallet.TapRootPublicKey
			input.Type = "taproot" // TODO: detect from wallet type
		} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
			input.SenderAddr = p.Trade.Taker.Wallet.TapRootAddr
			input.SenderPublicKey = p.Trade.Taker.Wallet.TapRootPublicKey
			input.Type = "segwit"
		} else {
			// TODO: all logging
			return errors.New("Invalid taker wallet type")
		}
		input.TxID = utxo.TxHash
		input.Index = utxo.TxOutputN
		input.WitnessUTXO.Amount = utxo.Value
		input.WitnessUTXO.Script = utxo.Script

		// Append to psbt inputs
		p.Inputs[len(p.Inputs)+i] = input
	}

	return nil
}

//type Output struct {
//	RecipientAddr      string `json:"recipient_addr"`
//	Value              int64  `json:"value"`
//}

func (p *PSBT) createPaymentsOutputs() error {
	// MAKER: Create the payments outputs
	// Create output for payment from maker to taker
	if p.Trade.Maker.BTC > 546 {
		makerPaymentOutput := new(Output)
		if p.Trade.Taker.Wallet.Type == "unisat" {
			makerPaymentOutput.RecipientAddr = p.Trade.Taker.Wallet.TapRootAddr
		} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
			makerPaymentOutput.RecipientAddr = p.Trade.Taker.Wallet.SegwitAddr
		} else {
			// TODO: all logging
			return errors.New("Invalid taker wallet type")
		}
		makerPaymentOutput.Value = p.Trade.Maker.BTC
		p.Outputs[len(p.Outputs)] = makerPaymentOutput
	}
	// Create output for maker change
	makerChangeOutput := new(Output)
	if p.Trade.Maker.Wallet.Type == "unisat" {
		makerChangeOutput.RecipientAddr = p.Trade.Maker.Wallet.TapRootAddr
	} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
		makerChangeOutput.RecipientAddr = p.Trade.Maker.Wallet.SegwitAddr
	} else {
		// TODO: all logging
		return errors.New("Invalid maker wallet type")
	}
	makerChangeOutput.Value = p.MakerChange
	p.Outputs[len(p.Outputs)] = makerChangeOutput
	// Create output for maker platform fee
	makerPlatformFeeOutput := new(Output)
	makerPlatformFeeOutput.RecipientAddr = "3C7trrWesxpM5aHPTCPMeeBG418C5LvXbc"
	makerPlatformFeeOutput.Value = p.PlatformFee / 2
	p.Outputs[len(p.Outputs)] = makerPlatformFeeOutput

	// TAKER: Create the payments outputs
	// Create output for payment from taker to maker
	if p.Trade.Taker.BTC > 546 {
		takerPaymentOutput := new(Output)
		if p.Trade.Maker.Wallet.Type == "unisat" {
			takerPaymentOutput.RecipientAddr = p.Trade.Maker.Wallet.TapRootAddr
		} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
			takerPaymentOutput.RecipientAddr = p.Trade.Maker.Wallet.SegwitAddr
		} else {
			// TODO: all logging
			return errors.New("Invalid taker wallet type")
		}
		takerPaymentOutput.Value = p.Trade.Taker.BTC
		p.Outputs[len(p.Outputs)] = takerPaymentOutput
	}
	// Create output for taker change
	takerChangeOutput := new(Output)
	if p.Trade.Taker.Wallet.Type == "unisat" {
		takerChangeOutput.RecipientAddr = p.Trade.Taker.Wallet.TapRootAddr
	} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
		takerChangeOutput.RecipientAddr = p.Trade.Taker.Wallet.SegwitAddr
	} else {
		// TODO: all logging
		return errors.New("Invalid maker wallet type")
	}
	takerChangeOutput.Value = p.TakerChange
	p.Outputs[len(p.Outputs)] = takerChangeOutput
	// Create output for maker platform fee
	takerPlatformFeeOutput := new(Output)
	takerPlatformFeeOutput.RecipientAddr = "3C7trrWesxpM5aHPTCPMeeBG418C5LvXbc"
	takerPlatformFeeOutput.Value = p.PlatformFee / 2
	p.Outputs[len(p.Outputs)] = takerPlatformFeeOutput

	return nil
}

// isValidTxID checks if the given txID adheres to the expected format of a Bitcoin transaction ID.
func isValidTxID(txID string) bool {
	pattern := `^[a-fA-F0-9]{64}$`
	match, _ := regexp.MatchString(pattern, txID)

	return match
}

func calculateMinerFeeForPSBT(tx *wire.MsgTx, feeRatePerVByte float64) (int64, error) {
	baseSize := tx.SerializeSizeStripped()
	totalSize := tx.SerializeSize()

	weight := baseSize*3 + totalSize
	virtualSize := int(math.Ceil(float64(weight) / 4.0))
	minerFee := float64(virtualSize) * feeRatePerVByte

	return int64(minerFee) + 50, nil
}
