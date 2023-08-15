package handlers

import (
	"fmt"
	"net/http"
	"ntc-services/models"
	"ntc-services/services"
	"ntc-services/stores"

	"github.com/labstack/echo/v4"
)

// TODO: Should we be using signatures to do RBAC
// TODO: other stuff

/* Request Body
{
  "wallet_type": "hiro" | "xverse" | "unisat",
  "taproot_addr": "addr | nil",
  "segwit_addr": "addr | nil"
}
*/

func PostTrades(c echo.Context) error {
	// TODO: find & verify wallet
	tradeReqBody := models.NewTradeReqBody()
	if err := c.Bind(tradeReqBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	wallet, err := models.GetWalletByAddrAndWalletType(tradeReqBody.TapRootAddr, tradeReqBody.SegwitAddr, tradeReqBody.WalletType)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if wallet == nil {
		wallet = models.NewWallet()
		wallet.TapRootAddr = tradeReqBody.TapRootAddr
		wallet.SegwitAddr = tradeReqBody.SegwitAddr
		if err := wallet.Save(); err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	// TODO: create side & store
	side := models.NewSide(wallet.ID)
	if err := side.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// TODO: create trade & store
	trade := models.NewTrade(side.ID)
	trade.Maker = side
	if err := trade.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, trade)
}

/* Request Body
{
  "wallet_id": "bson_id",
  "btc": 1000,
  "inscription_numbers": [1234, 2344]
}
*/

func PostMakerByTradeID(c echo.Context) error {
	tradeID := c.Param("id")

	// TODO: find & verify wallet
	tradeMakerReqBody := models.NewTradeMakerReqBody()
	if err := c.Bind(tradeMakerReqBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	wallet, err := models.GetWalletByID(tradeMakerReqBody.WalletID)
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// TODO: find trade and ensure in correct state (CREATED)
	trade, err := models.GetTradeByID(c, tradeID)
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// TODO: find side and ensure requester is correct by wallet_id (and perhaps more)
	maker, err := models.GetSideByID(trade.MakerID.Hex())
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if maker.WalletID != wallet.ID {
		c.Logger().Error("Maker Wallet does not match Wallet ID")
		return c.JSON(http.StatusConflict, "Maker Wallet does not match Wallet ID")
	}
	// TODO: query ordex for extra inscription information (floor price, previous tx, more...)
	for _, value := range tradeMakerReqBody.InscriptionNumbers {
		inscription, err := services.BESTINSLOT.GetInscriptionById(c, value)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		fmt.Printf("inscription: %+v\n", inscription)
	}

	// TODO: validate that assets still belong to maker wallet
	// TODO: update side
	maker.BTC = tradeMakerReqBody.Btc
	if err := maker.Update(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, maker)
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
