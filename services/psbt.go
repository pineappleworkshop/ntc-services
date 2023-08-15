package services

import (
	"errors"
	"github.com/btcsuite/btcd/btcutil/psbt"
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
	MakerChange           int64
	TakerChange           int64
	PlatformFee           int64
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
	}
}

func CreatePSBT() (*psbt.Packet, error) {
	// TODO: create all inputs
	// TODO: create all outputs
	// TODO: assemble PSBT

	return nil, nil
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
	for _, inscription := range p.Trade.Maker.Inscriptions {
		if inscription.FloorPrice != nil {
			totalInscriptionFloorPrices = totalInscriptionFloorPrices + *inscription.FloorPrice
		} else {
			totalInscriptionFloorPrices = totalInscriptionFloorPrices + int64(10000)
		}
	}
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
	// Iterate through payment UTXOs and use UTXOs until payment is greater than needed
	makerPayment := int64(0) + p.PlatformFee/2
	for i, utxo := range p.MakerUTXOsAll {
		makerPayment = makerPayment + utxo.Value
		p.MakerPaymentUTXOs = append(p.MakerPaymentUTXOs, utxo)
		p.MakerUTXOsAll = append(p.MakerUTXOsAll[:i], p.MakerUTXOsAll[i+1:]...)

		// Check to see if payment is greater than needed
		if makerPayment > p.Trade.Maker.BTC {
			break
		}
	}
	// calculate maker change
	p.MakerChange = makerPayment - p.Trade.Maker.BTC

	// TAKER: Select the payment UTXOs for the psbt
	// Iterate through payment UTXOs and use UTXOs until payment is greater than needed
	takerPayment := int64(0) + p.PlatformFee/2
	for i, utxo := range p.TakerUTXOsAll {
		takerPayment = takerPayment + utxo.Value
		p.TakerPaymentUTXOs = append(p.TakerPaymentUTXOs, utxo)
		p.TakerUTXOsAll = append(p.TakerUTXOsAll[:i], p.TakerUTXOsAll[i+1:]...)

		// Check to see if payment is greater than needed
		if takerPayment > p.Trade.Taker.BTC {
			break
		}
	}
	// calculate maker change
	p.TakerChange = takerPayment - p.Trade.Taker.BTC

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

	// TODO: TAKER: Calculate the change to return

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

// isValidTxID checks if the given txID adheres to the expected format of a Bitcoin transaction ID.
func isValidTxID(txID string) bool {
	// Regular expression to match a 64-character hexadecimal string
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
