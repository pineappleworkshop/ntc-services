package services

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"github.com/labstack/echo/v4"
	"math"
	"ntc-services/models"
	"regexp"
	"strconv"
	"strings"
)

type PSBTReqBody struct {
	Inputs  map[int]*models.Input  `json:"inputs"`
	Outputs map[int]*models.Output `json:"outputs"`
}

type PSBTRespBody struct {
	//Bytes  []byte `json:"bytes"`
	Base64 string `json:"base64"`
}

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
	MinerFee              int64
	Inputs                map[int]*models.Input
	Outputs               map[int]*models.Output
	PreMinerFeePSBT       *models.PSBTSerialized
	PostMinerFeePSBT      *models.PSBTSerialized
	FinalizedPSBT         *models.PSBTSerialized
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
		Inputs:                map[int]*models.Input{},
		Outputs:               map[int]*models.Output{},
	}
}

func (p *PSBT) Create(c echo.Context) error {
	if err := p.selectInscriptionsUTXOs(c); err != nil {
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

func (p *PSBT) selectInscriptionsUTXOs(c echo.Context) error {
	// MAKER: parse inscription UTXOs from all UTXOs
	for i, utxo := range p.MakerUTXOsAll {
		// Find matching UTXOs for inscriptions in trade and add
		for _, inscription := range p.Trade.Maker.Inscriptions {
			// Ensure inscriptionID has txID & index
			inscriptionIdS := strings.Split(inscription.Satpoint, ":")
			if len(inscriptionIdS) != 3 {
				err := errors.New("maker Inscription.Satpoint is not in the right format")
				c.Logger().Error(err)
				return err
			}
			// Ensure txID is valid
			if valid := isValidTxID(inscriptionIdS[0]); !valid {
				err := errors.New("maker Inscription.TxID is not valid")
				c.Logger().Error(err)
				return err
			}
			// Ensure index in int
			if _, err := strconv.Atoi(inscriptionIdS[1]); err != nil {
				err := errors.New("maker Inscription.Index is not an integer")
				c.Logger().Error(err)
				return err
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
	}
	// MAKER: remove all inscription UTXOs from all UTXOs
	for _, inscription := range p.MakerInscriptionsAll {
		// Ensure inscriptionID has txID & index
		inscriptionIdS := strings.Split(inscription.Satpoint, ":")
		if len(inscriptionIdS) != 3 {
			err := errors.New("maker Inscription.Satpoint is not in the right format")
			c.Logger().Error(err)
			return err
		}
		// Ensure txID is valid
		if valid := isValidTxID(inscriptionIdS[0]); !valid {
			err := errors.New("maker Inscription.TxID is not valid")
			c.Logger().Error(err)
			return err
		}
		// Ensure index in int
		index, err := strconv.Atoi(inscriptionIdS[1])
		if err != nil {
			err := errors.New("maker Inscription.Index is not an integer")
			c.Logger().Error(err)
			return err
		}
		for i, utxo := range p.MakerUTXOsAll {
			// Remove inscription for original list
			if utxo.TxHashBigEndian == inscriptionIdS[0] {
				if utxo.TxOutputN == int64(index) {
					p.MakerUTXOsAll = append(p.MakerUTXOsAll[:i], p.MakerUTXOsAll[i+1:]...)
					break
				}
			}
		}
	}

	// TAKER: parse inscription UTXOs from all UTXOs
	for i, utxo := range p.TakerUTXOsAll {
		// Find matching UTXOs for inscriptions in trade and add
		for _, inscription := range p.Trade.Taker.Inscriptions {
			// Ensure inscriptionID has txID & index
			inscriptionIdS := strings.Split(inscription.Satpoint, ":")
			if len(inscriptionIdS) != 3 {
				err := errors.New("taker Inscription.Satpoint is not in the right format")
				c.Logger().Error(err)
				return err
			}
			// Ensure txID is valid
			if valid := isValidTxID(inscriptionIdS[0]); !valid {
				err := errors.New("maker Inscription.TxID is not valid")
				c.Logger().Error(err)
				return err
			}
			// Ensure index in int
			if _, err := strconv.Atoi(inscriptionIdS[1]); err != nil {
				err := errors.New("maker Inscription.Index is not an integer")
				c.Logger().Error(err)
				return err
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
	}
	// Taker: remove all inscription UTXOs from all UTXOs
	for _, inscription := range p.TakerInscriptionsAll {
		// Ensure inscriptionID has txID & index
		inscriptionIdS := strings.Split(inscription.Satpoint, ":")
		if len(inscriptionIdS) != 3 {
			err := errors.New("taker Inscription.Satpoint is not in the right format")
			c.Logger().Error(err)
			return err
		}

		// Ensure txID is valid
		if valid := isValidTxID(inscriptionIdS[0]); !valid {
			err := errors.New("taker Inscription.TxID is not valid")
			c.Logger().Error(err)
			return err
		}

		// Ensure index in int
		index, err := strconv.Atoi(inscriptionIdS[1])
		if err != nil {
			err := errors.New("taker Inscription.Index is not an integer")
			c.Logger().Error(err)
			return err
		}

		for i, utxo := range p.TakerUTXOsAll {
			// Remove inscription for original list
			if utxo.TxHashBigEndian == inscriptionIdS[0] {
				if utxo.TxOutputN == int64(index) {
					p.TakerUTXOsAll = append(p.TakerUTXOsAll[:i], p.TakerUTXOsAll[i+1:]...)
					break
				}
			}
		}
	}

	return nil
}

// TODO: revisit and handle beta testing with minimal fees
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
			totalInscriptionFloorPrices = totalInscriptionFloorPrices + int64(12000)
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
		if len(p.MakerUTXOsAll) > 1 {
			p.MakerUTXOsAll = append(p.MakerUTXOsAll[:i], p.MakerUTXOsAll[i+1:]...)
		} else {
			p.MakerUTXOsAll = []*models.UTXO{}
		}

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
		if len(p.TakerUTXOsAll) > 1 {
			p.TakerUTXOsAll = append(p.TakerUTXOsAll[:i], p.TakerUTXOsAll[i+1:]...)
		} else {
			p.TakerUTXOsAll = []*models.UTXO{}
		}

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
		input := new(models.Input)
		input.Addr = p.Trade.Maker.Wallet.TapRootAddr
		if p.Trade.Maker.Wallet.Type == "unisat" {
			input.PublicKey = p.Trade.Maker.Wallet.TapRootPublicKey[2:]
		} else { // xverse or hiro
			input.PublicKey = p.Trade.Maker.Wallet.TapRootPublicKey
		}
		input.Type = "p2tr"
		input.TxID = utxo.TxHashBigEndian
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
		input := new(models.Input)
		input.Addr = p.Trade.Taker.Wallet.TapRootAddr
		if p.Trade.Taker.Wallet.Type == "unisat" {
			input.PublicKey = p.Trade.Taker.Wallet.TapRootPublicKey[2:]
		} else { // xverse or hiro
			input.PublicKey = p.Trade.Taker.Wallet.TapRootPublicKey
		}
		input.Type = "p2tr"
		input.TxID = utxo.TxHashBigEndian
		input.Index = utxo.TxOutputN
		input.WitnessUTXO.Amount = utxo.Value
		input.WitnessUTXO.Script = utxo.Script

		// Append to psbt inputs
		p.Inputs[i] = input
	}

	return nil
}

func (p *PSBT) createInscriptionOutputs() error {
	// MAKER: Create the inscription outputs
	// Iterate through the maker inscriptions UTXOs and create outputs to the recipient
	for i, utxo := range p.MakerInscriptionUTXOs {
		// Create output
		output := new(models.Output)
		output.Addr = p.Trade.Taker.Wallet.TapRootAddr
		output.Value = utxo.Value
		output.IsInscription = true

		// append to outputs
		p.Outputs[i] = output
	}

	// TAKER: Create the inscription outputs
	// Iterate through the taker inscriptions UTXOs and create outputs to the recipient
	for i, utxo := range p.TakerInscriptionUTXOs {
		// Create output
		output := new(models.Output)
		output.Addr = p.Trade.Maker.Wallet.TapRootAddr
		output.Value = utxo.Value
		output.IsInscription = true

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
		input := new(models.Input)
		if p.Trade.Maker.Wallet.Type == "unisat" {
			input.Addr = p.Trade.Maker.Wallet.TapRootAddr
			input.PublicKey = p.Trade.Maker.Wallet.TapRootPublicKey[2:]
			input.Type = "p2tr"
		} else if p.Trade.Maker.Wallet.Type == "xverse" || p.Trade.Maker.Wallet.Type == "hiro" {
			input.Addr = p.Trade.Maker.Wallet.SegwitAddr
			input.PublicKey = p.Trade.Maker.Wallet.SegwitPublicKey
			input.Type = "p2wpkh"
		} else {
			// TODO: all logging
			return errors.New("Invalid maker wallet type")
		}
		input.TxID = utxo.TxHashBigEndian
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
		input := new(models.Input)
		if p.Trade.Taker.Wallet.Type == "unisat" {
			input.Addr = p.Trade.Taker.Wallet.TapRootAddr
			input.PublicKey = p.Trade.Taker.Wallet.TapRootPublicKey[2:]
			input.Type = "p2tr"
		} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
			input.Addr = p.Trade.Taker.Wallet.SegwitAddr
			input.PublicKey = p.Trade.Taker.Wallet.SegwitPublicKey
			input.Type = "p2wpkh"
		} else {
			// TODO: all logging
			return errors.New("Invalid taker wallet type")
		}
		input.TxID = utxo.TxHashBigEndian
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
	if p.Trade.Maker.BTC >= 546 {
		makerPaymentOutput := new(models.Output)
		if p.Trade.Taker.Wallet.Type == "unisat" {
			makerPaymentOutput.Addr = p.Trade.Taker.Wallet.TapRootAddr
		} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
			makerPaymentOutput.Addr = p.Trade.Taker.Wallet.SegwitAddr
		} else {
			// TODO: all logging
			return errors.New("Invalid taker wallet type")
		}
		makerPaymentOutput.Value = p.Trade.Maker.BTC
		makerPaymentOutput.IsPayment = true
		p.Outputs[len(p.Outputs)] = makerPaymentOutput
	}
	// Create output for maker change
	makerChangeOutput := new(models.Output)
	if p.Trade.Maker.Wallet.Type == "unisat" {
		makerChangeOutput.Addr = p.Trade.Maker.Wallet.TapRootAddr
	} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
		makerChangeOutput.Addr = p.Trade.Maker.Wallet.SegwitAddr
	} else {
		// TODO: all logging
		return errors.New("Invalid maker wallet type")
	}
	makerChangeOutput.Value = p.MakerChange
	makerChangeOutput.IsChange = true
	p.Outputs[len(p.Outputs)] = makerChangeOutput
	// Create output for maker platform fee
	makerPlatformFeeOutput := new(models.Output)
	makerPlatformFeeOutput.Addr = "3C7trrWesxpM5aHPTCPMeeBG418C5LvXbc"
	makerPlatformFeeOutput.Value = p.PlatformFee / 2
	if makerPlatformFeeOutput.Value < 600 {
		makerPlatformFeeOutput.Value = 600
	}
	makerPlatformFeeOutput.IsFee = true
	p.Outputs[len(p.Outputs)] = makerPlatformFeeOutput

	// TAKER: Create the payments outputs
	// Create output for payment from taker to maker
	if p.Trade.Taker.BTC >= 546 {
		takerPaymentOutput := new(models.Output)
		if p.Trade.Maker.Wallet.Type == "unisat" {
			takerPaymentOutput.Addr = p.Trade.Maker.Wallet.TapRootAddr
		} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
			takerPaymentOutput.Addr = p.Trade.Maker.Wallet.SegwitAddr
		} else {
			// TODO: all logging
			return errors.New("Invalid taker wallet type")
		}
		takerPaymentOutput.Value = p.Trade.Taker.BTC
		takerPaymentOutput.IsPayment = true
		p.Outputs[len(p.Outputs)] = takerPaymentOutput
	}
	// Create output for taker change
	takerChangeOutput := new(models.Output)
	if p.Trade.Taker.Wallet.Type == "unisat" {
		takerChangeOutput.Addr = p.Trade.Taker.Wallet.TapRootAddr
	} else if p.Trade.Taker.Wallet.Type == "xverse" || p.Trade.Taker.Wallet.Type == "hiro" {
		takerChangeOutput.Addr = p.Trade.Taker.Wallet.SegwitAddr
	} else {
		// TODO: all logging
		return errors.New("Invalid maker wallet type")
	}
	takerChangeOutput.Value = p.TakerChange
	takerChangeOutput.IsChange = true
	p.Outputs[len(p.Outputs)] = takerChangeOutput
	// Create output for maker platform fee
	takerPlatformFeeOutput := new(models.Output)
	takerPlatformFeeOutput.Addr = "3C7trrWesxpM5aHPTCPMeeBG418C5LvXbc"
	takerPlatformFeeOutput.Value = p.PlatformFee / 2
	if takerPlatformFeeOutput.Value < 600 {
		takerPlatformFeeOutput.Value = 600
	}
	takerPlatformFeeOutput.IsFee = true
	p.Outputs[len(p.Outputs)] = takerPlatformFeeOutput

	return nil
}

func (p *PSBT) adjustForMinerFee(minerFee int64) error {
	feePerSide := minerFee / 2
	p.MakerChange = p.MakerChange - feePerSide
	p.TakerChange = p.TakerChange - feePerSide

	for _, output := range p.Outputs {
		if output.IsChange {
			output.Value = output.Value - feePerSide
		}
	}

	return nil
}

func (p *PSBT) GeneratePSBT(c echo.Context) (*models.PSBT, error) {
	// Send psbt to psbt service to create initial psbt without gas
	preMinerFeePSBTResp, err := NTCPSBT.PostPSBT(p.ToReq())
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	p.PreMinerFeePSBT = preMinerFeePSBTResp
	preMinerFeePSBT, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(preMinerFeePSBTResp.Base64)), true)
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}

	// Measure gas and adjust the maker/taker change outputs
	minerFee, err := p.calculateMinerFeeForPSBT(preMinerFeePSBT.UnsignedTx, float64(p.Trade.FeeRate))
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}

	p.MinerFee = minerFee
	if err := p.adjustForMinerFee(minerFee); err != nil {
		c.Logger().Error(err)
		return nil, err
	}

	// Send finalized Tx to PSBT service
	postMinerFeePSBTResp, err := NTCPSBT.PostPSBT(p.ToReq())
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	p.PostMinerFeePSBT = postMinerFeePSBTResp

	psbtModel := &models.PSBT{
		MakerPayment:     p.MakerPayment,
		TakerPayment:     p.TakerPayment,
		MakerChange:      p.MakerChange,
		TakerChange:      p.TakerChange,
		PlatformFee:      p.PlatformFee,
		MinerFee:         p.MinerFee,
		Inputs:           p.Inputs,
		Outputs:          p.Outputs,
		PreMinerFeePSBT:  p.PreMinerFeePSBT,
		PostMinerFeePSBT: p.PostMinerFeePSBT,
		FinalizedPSBT:    nil,
	}

	return psbtModel, nil
}

// isValidTxID checks if the given txID adheres to the expected format of a Bitcoin transaction ID.
func isValidTxID(txID string) bool {
	pattern := `^[a-fA-F0-9]{64}$`
	match, _ := regexp.MatchString(pattern, txID)

	return match
}

func (p *PSBT) calculateMinerFeeForPSBT(tx *wire.MsgTx, feeRatePerVByte float64) (int64, error) {
	baseSize := tx.SerializeSizeStripped()
	totalSize := tx.SerializeSize()

	weight := baseSize*3 + totalSize
	virtualSize := int(math.Ceil(float64(weight) / 4.0))

	for _, input := range p.Inputs {
		if input.Type == "p2tr" {
			virtualSize = virtualSize + 57
		}
		if input.Type == "p2wpkh" {
			virtualSize = virtualSize + 102
		}
	}

	minerFee := float64(virtualSize) * feeRatePerVByte

	return int64(minerFee) + 50, nil
}

func extractTaprootPublicKeyHex(pkScript []byte) (string, error) {
	// Check that the script is the correct length for P2TR (1-byte version + 1-byte length + 32-byte key)
	if len(pkScript) != 34 {
		return "", fmt.Errorf("invalid pkScript length for P2TR")
	}

	// Check that the version byte is 0x51 (indicating P2TR)
	if pkScript[0] != 0x51 {
		return "", fmt.Errorf("invalid version byte for P2TR")
	}

	// Check that the length byte is 0x20 (indicating 32 bytes)
	if pkScript[1] != 0x20 {
		return "", fmt.Errorf("invalid length byte for P2TR public key")
	}

	// Convert the 32-byte taproot public key to hexadecimal format
	pubKeyHex := hex.EncodeToString(pkScript[2:])
	return pubKeyHex, nil
}

func extractTaprootInternalPublicKeyFromAddress(addr string) (string, error) {
	// TODO: check address first to see if taproot, beginning starts with `bc1p`
	decodedAddress, err := btcutil.DecodeAddress(addr, nil)
	if err != nil {
		return "", err
	}
	if wa, ok := decodedAddress.(*btcutil.AddressTaproot); ok {
		encodedKeyHex := hex.EncodeToString(wa.WitnessProgram())
		return encodedKeyHex, nil
	}

	return "", errors.New("address is not a taproot")
}

func detectAddrType(addr string) (string, error) {
	if strings.HasPrefix(addr, "3") {
		return "p2sh", nil
	}
	if strings.HasPrefix(addr, "bc1p") {
		return "p2tr", nil
	}

	return "", errors.New("address not supported")
}

//b9907521ddb85e0e6a37622b7c685efbdc8ae53a334928adbd12cf204ad4e717
//02818b7ff740a40f311d002123087053d5d9e0e1546674aedb10e15a5b57fd3985
//0223dfdbe72c5ee9e687946e9c17f68589d90552e37a6435da7c05c2f1fba21f15

//0223dfdbe72c5ee9e687946e9c17f68589d90552e37a6435da7c05c2f1fba21f15 unisat pk
//412a65ed2a0d2d5237f0a7ab338d0cc3b737daea28436138a68a40ab7f7ad4a7 xverse ordinal pk
//02015e26bf3ff8e7a470ec07b719a7414b6e6f8a67ef97302e28fd524b2d9b7ec2 xverse payment pk
