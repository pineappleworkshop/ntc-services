package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// TODO: Should we be using signatures to do RBAC
// TODO: other stuff

/* Request Body
{
  "wallet_id": "bson_id"
}
*/

func PostTrades(c echo.Context) error {

	// TODO: find & verify wallet
	// TODO: create side & store
	// TODO: create trade & store

	return c.JSON(http.StatusCreated, nil) // TODO: return trade
}

/* Request Body
{
  "wallet_id": "bson_id",
  "btc": 1000,
  "inscription_numbers": [1234, 2344]
}
*/

func PostMakerByTradeID(c echo.Context) error {

	// TODO: find & verify wallet
	// TODO: find trade and ensure in correct state (CREATED)
	// TODO: find side and ensure requester is correct by wallet_id (and perhaps more)
	// TODO: query ordex for extra inscription information (floor price, previous tx, more...)
	// TODO: validate that assets still belong to maker wallet
	// TODO: update side

	return c.JSON(http.StatusCreated, nil)
}

/* Query Params
?status={enum,csv}
*/

func GetTrades(c echo.Context) error {

	// TODO: get trades by query
	// TODO: paginated response

	return c.JSON(http.StatusOK, nil)
}

/* Request Body
{
  "wallet_id": "bson_id",
  "btc": 1000,
  "inscription_numbers": [1234, 2344]
}
*/

func PostOfferByTradeID(c echo.Context) error {

	// TODO: find & verify wallet
	// TODO: find & verify trade is in correct status
	// TODO: validate that assets still belong to maker wallet
	// TODO: create offer for trade

	return c.JSON(http.StatusCreated, nil)
}

/* Query Params
?status={enum,csv}
*/

func GetOffersByTradeID(c echo.Context) error {

	// TODO: get trades for trade by query
	// TODO: paginated response

	return c.JSON(http.StatusOK, nil)
}

/*
{}
*/

func PostAcceptOfferByTradeID(c echo.Context) error {

	// TODO: find & verify wallet
	// TODO: find & verify trade is in correct status
	// TODO: find & verify offer is in correct status
	// TODO: find & verify trade assets are all correct wallet
	// TODO: find & verify offer assets are all in the correct wallet
	// TODO: create PBST (this is for platform control)

	return c.JSON(http.StatusOK, nil)
}

/*
{}
*/

func PostSubmitTradeByID(c echo.Context) error {

	// TODO: Unwrap PBST to ensure its validity
	// TODO: Sign PBST and send to RPC gateway

	return c.JSON(http.StatusOK, nil)
}
