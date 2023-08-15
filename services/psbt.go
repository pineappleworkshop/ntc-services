package services

import (
	"errors"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"ntc-services/models"
	"regexp"
	"strconv"
	"strings"
)

type PSBT struct {
	Trade                 *models.Trade
	MakerUTXOs            []*models.UTXO // All maker utxos
	TakerUTXOs            []*models.UTXO // All taker utxos
	MakerInscriptionUTXOs []*models.UTXO // Maker inscription utxos for trade
	MakerOtherUTXOs       []*models.UTXO // Maker other utcos for trade
	TakerInscriptionUTXOs []*models.UTXO // Taker inscription utxos for trade
	TakerOtherUTXOs       []*models.UTXO // Taker other utcos for trade
	MakerChange           int64
	TakerChange           int64
	PlatformFee           int64
}

func NewPBST(trade *models.Trade, makerUTXOs, takerUTXOs []*models.UTXO) *PSBT {
	return &PSBT{
		Trade:                 trade,
		MakerUTXOs:            makerUTXOs,
		TakerUTXOs:            takerUTXOs,
		MakerInscriptionUTXOs: []*models.UTXO{},
		MakerOtherUTXOs:       []*models.UTXO{},
		TakerInscriptionUTXOs: []*models.UTXO{},
		TakerOtherUTXOs:       []*models.UTXO{},
	}
}

func CreatePSBT() (*psbt.Packet, error) {

	// TODO: create all inputs
	// TODO: create all outputs
	// TODO: assemble PSBT

	return nil, nil
}

func (p *PSBT) selectInscriptionsUTXOs() error {
	// MAKER: parse inscription UTXOs from other UTXOs
	for _, utxo := range p.MakerUTXOs {
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

			// Parse inscription utxos and other utxos
			if utxo.TxHashBigEndian == inscriptionIdS[0] {
				p.MakerInscriptionUTXOs = append(p.MakerInscriptionUTXOs, utxo)
			}
		}
	}

	// TAKER: parse inscription UTXOs from other UTXOs
	for _, utxo := range p.TakerUTXOs {
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

			// Add inscription to psbt inscription utxos
			if utxo.TxHashBigEndian == inscriptionIdS[0] {
				p.TakerInscriptionUTXOs = append(p.TakerInscriptionUTXOs, utxo)
			}
		}
	}

	return nil
}

func (p *PSBT) selectPaymentUTXOs() error {

	// TODO: MAKER: Select the proper payment UTXOs for the psbt

	// TODO: TAKER: Select the proper payment UTXOs for the psbt

	return nil
}

func (p *PSBT) createInscriptionInputs() error {

	// TODO: MAKER: Create the inscription Inputs

	// TODO: TAKER: Create the inscription Inputs

	return nil
}

func (p *PSBT) createPaymentInputs() error {

	// TODO: MAKER: Create the inscription Inputs

	// TODO: TAKER: Create the inscription Inputs

	return nil
}

func (p *PSBT) calculateChange() error {

	// TODO: MAKER: Calculate the change to return

	// TODO: MAKER: Calculate the change to return

	return nil
}

func (p *PSBT) createInscriptionOutputs() error {

	// TODO: MAKER: Create the inscription outputs

	// TODO: TAKER: Create the inscription outputs

	return nil
}

func (p *PSBT) createPaymentsOutputs() error {

	// TODO: MAKER: Create the payments outputs

	// TODO: TAKER: Create the payments outputs

	// TODO: FEE: Create the payments outputs

	// TODO: MINER FEE: Create the payments outputs

	return nil
}

func (p *PSBT) f() error {

	// TODO:

	return nil
}

// isValidTxID checks if the given txID adheres to the expected format of a Bitcoin transaction ID.
func isValidTxID(txID string) bool {
	// Regular expression to match a 64-character hexadecimal string
	pattern := `^[a-fA-F0-9]{64}$`
	match, _ := regexp.MatchString(pattern, txID)
	return match
}
