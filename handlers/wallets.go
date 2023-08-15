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

	walletExistingTapRoot, err := models.GetWalletByAddr(wallet.TapRootAddr, models.ADDRESS_TAPROOT)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if walletExistingTapRoot != nil {
		c.Logger().Error("Wallet already exists in database")
		return c.JSON(http.StatusConflict, "Wallet already exists in database")
	}
	walletExistingSegwit, err := models.GetWalletByAddr(wallet.SegwitAddr, models.ADDRESS_SEGWIT)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if walletExistingSegwit != nil {
		c.Logger().Error("Wallet already exists in database")
		return c.JSON(http.StatusConflict, "Wallet already exists in database")
	}

	if err := wallet.Save(); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, wallet)
}

func GetWalletsByAddr(c echo.Context) error {
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
