package handlers

import (
	"fmt"
	"net/http"
	"ntc-services/models"

	"github.com/labstack/echo/v4"
)

/* Request Body
{
  "wallet_type": "hiro" | "xverse" | "unisat",
  "cardinal_addr": "addr | nil",
  "taproot_addr": "addr | nil",
  "segwit_addr": "addr | nil"
}
The wallet types right now will be:
hiro:	Has segit & taproot addresses
xverse:	Has segit & taproot addresses
unisat:	Only has taproot address
*/

func PostWallets(c echo.Context) error {
	wallet := models.NewWallet()
	if err := c.Bind(wallet); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if models.IsValidWalletType(wallet.Type) == false {
		c.Logger().Error("Invalid Wallet Address Type: ", wallet.Type)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("Invalid Wallet Address Type: %s", wallet.Type)})
	}

	walletExisting, err := models.GetWalletByAddr(wallet.TapRootAddr, wallet.SegwitAddr, wallet.Type)
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

func PostWalletsConnected(c echo.Context) error {
	// todo: define what is going to be passed here (the wallet_id?) so we can find the wallet in DB and update:
	// `last_connected_at` & `last_connected_block`

	return c.JSON(http.StatusCreated, nil)
}
