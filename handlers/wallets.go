package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"ntc-services/models"
	"ntc-services/services"
	"ntc-services/stores"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

/* Request Body
{
  "wallet_type": "hiro" | "xverse" | "unisat",
  "taproot_addr": "addr",
  "taproot_public_key": "string",
  "segwit_addr": "addr | nil",
  "seqwit_public_key": "addr | nil"
}
The wallet types right now will be:
hiro:	Has segit & taproot addresses
xverse:	Has segit & taproot addresses
unisat:	Only has taproot address
*/

func PostWallets(c echo.Context) error {
	// Parse wallet and validate request fields
	wallet := models.NewWallet()
	if err := c.Bind(wallet); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if models.IsValidWalletType(wallet.Type) == false {
		err := errors.New(fmt.Sprintf("Invalid wallet address type: ", wallet.Type))
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := wallet.Validate(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Ensure wallet does not already exist
	walletExistingTapRoot, err := models.GetWalletByAddr(wallet.TapRootAddr, models.ADDRESS_TAPROOT)
	if err != nil {
		if err.Error() != stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	if walletExistingTapRoot != nil {
		c.Logger().Error("Wallet already exists in database")
		return c.JSON(http.StatusConflict, "Wallet already exists in database")
	}
	walletExistingSegwit, err := models.GetWalletByAddr(wallet.SegwitAddr, models.ADDRESS_SEGWIT)
	if err != nil {
		if err.Error() != stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	if walletExistingSegwit != nil {
		c.Logger().Error("Wallet already exists in database")
		return c.JSON(http.StatusConflict, "Wallet already exists in database")
	}

	// Store wallet
	if err := wallet.Save(); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, wallet)
}

func GetWallets(c echo.Context) error {
	wallets, err := models.GetWallets(c)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, wallets)
}

func GetWalletByID(c echo.Context) error {
	// Get wallet by id
	wallet, err := models.GetWalletByID(c.Param("id"))
	if err != nil {
		if err.Error() != stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Get all wallet inscriptions
	// TODO: cover wallets that have inscriptions greater then 100 (pagination)
	makerInscriptions, err := services.BESTINSLOT.GetInscriptionsByWalletAddr(
		c,
		wallet.TapRootAddr,
		100,
		1,
	)
	inscriptions := []*models.Inscription{}
	for _, makerInscription := range makerInscriptions.Data {
		inscription := models.ParseBISInscription(makerInscription)
		inscriptions = append(inscriptions, inscription)
	}
	wallet.Inscriptions = inscriptions

	// Get all wallet utxos
	// TODO: check blockinfo API response for pagination
	var makerPaymentAddr string
	if wallet.Type == "unisat" {
		makerPaymentAddr = wallet.TapRootAddr
	} else { // TODO: harden
		makerPaymentAddr = wallet.SegwitAddr
	}
	utxosResp, err := services.BLOCKCHAININFO.GetUTXOsForAddr(makerPaymentAddr)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, utxoI := range utxosResp["unspent_outputs"].([]interface{}) {
		utxo := new(models.UTXO)
		if err := utxo.Parse(utxoI.(map[string]interface{})); err != nil {
			err := errors.New(
				fmt.Sprintf("could not parse utxo from blockchain info in data schema"),
			)
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		for _, inscription := range wallet.Inscriptions {
			satpointS := strings.Split(inscription.Satpoint, ":")
			if satpointS[0] == utxo.TxHashBigEndian {
				index, err := strconv.Atoi(satpointS[1])
				if err != nil {
					c.Logger().Error(err)
					return c.JSON(http.StatusBadRequest, err.Error())
				}
				if utxo.TxOutputN == int64(index) {
					utxo.IsInscription = true
					utxo.InscriptionNum = &inscription.InscriptionNumber
				}
			}
		}
		wallet.UTXOs = append(wallet.UTXOs, utxo)
	}

	return c.JSON(http.StatusOK, wallet)
}

func GetWalletByAddr(c echo.Context) error {
	addr := c.Param("addr")
	addrType, err := models.GetAddressType(addr)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	wallet, err := models.GetWalletByAddr(addr, addrType)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if wallet == nil {
		wallet := models.NewWallet()
		wallet.Type = addrType
		if wallet.Type == models.ADDRESS_SEGWIT {
			wallet.SegwitAddr = addr
		} else if wallet.Type == models.ADDRESS_TAPROOT {
			wallet.TapRootAddr = addr
		} else {
			c.Logger().Error("Invalid Address")
			return c.JSON(http.StatusNotFound, "Invalid Address")
		}
		if err := wallet.Save(); err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, wallet)
}

func PostWalletsConnected(c echo.Context) error {
	// todo: define what is going to be passed here (the wallet_id?) so we can find the wallet in DB and update:
	// `last_connected_at` & `last_connected_block`

	return c.JSON(http.StatusCreated, nil)
}
