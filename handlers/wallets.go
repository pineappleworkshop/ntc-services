package handlers

import (
	"fmt"
	"net/http"
	"ntc-services/models"

	"github.com/labstack/echo/v4"
)

/* Request Body
{
  "wallet_type": "some supported wallet type",
  "cardinal_addr": "addr | nil",
  "taproot_addr": "addr | nil",
  "segwit_addr": "addr | nil"
}
*/

func PostWallets(c echo.Context) error {
	wallet := models.NewWallet()
	if err := c.Bind(wallet); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var walletAddr string
	switch wallet.Type {
	case "cardinal":
		walletAddr = wallet.CardinalAddr
	case "taproot":
		walletAddr = wallet.TapRootAddr
	case "segwit":
		walletAddr = wallet.SegwitAddr
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("Invalid Wallet Address Type: %s", wallet.Type)})
	}
	walletExisting, err := models.GetWalletByAddr(walletAddr, wallet.Type)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if walletExisting != nil {
		c.Logger().Error("Wallet already exists in database")
		return c.JSON(http.StatusConflict, "Wallet already exists in database")
	}
	if err := wallet.Save(); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, wallet)
}
