package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ntc-services/models"
	"ntc-services/services"
	"ntc-services/stores"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if trade.Status != "CREATED" {
		c.Logger().Error("Invalid Status: ", trade.Status)
		return c.JSON(http.StatusNotFound, "Invalid Status")
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
	// look for another endpoint called get inscriptions by id (multiple same time)
	response, err := services.ORDEX.GetInscriptionsByIds(tradeMakerReqBody.InscriptionNumbers)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	// fmt.Printf("inscriptions: %+v\n", inscriptions)
	// Format response as readable JSON
	formattedJSON, err := formatJSON(response)
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	fmt.Println("Formatted JSON:")
	fmt.Println(formattedJSON)

	// for _, value := range tradeMakerReqBody.InscriptionNumbers {
	// 	inscription, err := services.BESTINSLOT.GetInscriptionById(c, value)
	// 	if err != nil {
	// 		c.Logger().Error(err)
	// 		return c.JSON(http.StatusInternalServerError, err)
	// 	}
	// 	fmt.Printf("inscription: %+v\n", inscription)
	// }

	// TODO: validate that assets still belong to maker wallet
	// TODO: update side
	maker.BTC = tradeMakerReqBody.Btc
	if err := maker.Update(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, maker)
}

/*
	Query Params

?status={enum,csv}
*/
func formatJSON(data interface{}) (string, error) {
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

func GetTrades(c echo.Context) error {
	// TODO: get trades by query
	status := c.QueryParam("status")
	statusValues := strings.Split(status, ",")
	// TODO: paginated response
	page, limit, err := parsePagination(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	trades, total, err := models.GetTradesPaginatedByStatus(page, limit, statusValues)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	resp := models.Trades{
		Page:   page,
		Limit:  limit,
		Total:  total,
		Trades: trades,
	}

	return c.JSON(http.StatusOK, resp)
}

/* Request Body
{
  "wallet_id": "bson_id",
  "btc": 1000,
  "inscription_numbers": [1234, 2344]
}
*/

func PostOfferByTradeID(c echo.Context) error {
	tradeID := c.Param("id")

	// TODO: find & verify wallet
	tradeMakerReqBody := models.NewTradeMakerReqBody()
	if err := c.Bind(tradeMakerReqBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	_, err := models.GetWalletByID(tradeMakerReqBody.WalletID)
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// TODO: find & verify trade is in correct status
	trade, err := models.GetTradeByID(c, tradeID)
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if trade.Status != "CREATED" {
		c.Logger().Error("Invalid Status: ", trade.Status)
		return c.JSON(http.StatusNotFound, "Invalid Status")
	}

	// TODO: validate that assets still belong to maker wallet
	// TODO: create offer for trade
	offer := models.NewOffer(trade.ID)
	walletIDHex, err := primitive.ObjectIDFromHex(tradeMakerReqBody.WalletID)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	offer.MakerID = walletIDHex
	if err := offer.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, offer)
}

/* Query Params
?status={enum,csv}
*/

func GetOffersByTradeID(c echo.Context) error {
	// TODO: get trades for trade by query
	// TODO: paginated response
	tradeID := c.Param("id")
	status := c.QueryParam("status")
	statusValues := strings.Split(status, ",")
	page, limit, err := parsePagination(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	offers, total, err := models.GetOffersPaginatedByTradeID(page, limit, tradeID, statusValues)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	resp := models.Offers{
		Page:   page,
		Limit:  limit,
		Total:  total,
		Offers: offers,
	}

	return c.JSON(http.StatusOK, resp)

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
