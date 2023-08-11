package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

/* Request Body
{
  "type": "some supported wallet type",
  "cardinal_addr": "addr | nil",
  "taproot_addr": "addr | nil",
  "segwit_addr": "addr | nil"
}
*/

func PostWallets(c echo.Context) error {
	//wallet := models.NewWallet("some_type")

	// TODO: find wallet by type, addr (of sorts, more than likely different logic per type)
	// TODO: if not found, store

	return c.JSON(http.StatusCreated, nil)
}
