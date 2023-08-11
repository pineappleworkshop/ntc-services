package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

/* Request Body
{
  "wallet_id": "bson_id"
}
*/

func PostTrades(c echo.Context) error {

	// TODO: find and verify wallet
	// TODO: create side & store
	// TODO: create trade & store

	return c.JSON(http.StatusCreated, nil) // TODO: return trade
}

/* Request Body
{
  "wallet_id": "bson_id",
  "inscriptions": [{models.inscription}],
  "utxos": [{models.inscription}]
}
*/

func PostTradesByIDMaker(c echo.Context) error {

	// TODO: find trade and ensure in correct state (CREATED)
	// TODO: find side and ensure requester is correct by wallet_id (and perhaps more)
	// TODO: update side (should already be referenced by trade)

	return c.JSON(http.StatusCreated, nil)
}
