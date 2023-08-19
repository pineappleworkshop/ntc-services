package handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/models"
	"ntc-services/services"
	"ntc-services/stores"
	"strconv"
	"strings"
)

/* Request Body
{
  "wallet_type": "hiro" | "xverse" | "unisat",
  "taproot_addr": "addr | nil",
  "segwit_addr": "addr | nil"
}
*/

func PostTrades(c echo.Context) error {
	// Parse json body
	tradeReqBody := new(models.TradeReqBody)
	if err := c.Bind(tradeReqBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Find & verify wallet
	wallet, err := models.GetWalletByAddrAndWalletType(tradeReqBody.TapRootAddr, tradeReqBody.SegwitAddr, tradeReqBody.WalletType)
	if err != nil {
		if err.Error() != stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
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

	// Create side & store
	maker := models.NewSide(wallet.ID)
	if err := maker.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Create trade & store
	trade := models.NewTrade(maker.ID)
	trade.Maker = maker
	if err := trade.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, trade)
}

/* Request Body
{
  "wallet_id": "bson_id",
  "btc": 1000,
  "inscription_numbers": [1234, 2344]
}
*/

func PostMakerByTradeID(c echo.Context) error {
	// Find & verify wallet
	makerReqBody := new(models.SideReqBody)
	if err := c.Bind(makerReqBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	wallet, err := models.GetWalletByID(makerReqBody.WalletID)
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Find trade and ensure in correct state (CREATED)
	trade, err := models.GetTradeByID(c, c.Param("id"))
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if trade.Status != "CREATED" {
		err := errors.New(fmt.Sprintf("trade.status is not CREATED, status is: %v", trade.Status))
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Find side and ensure requester is correct by wallet_id (and perhaps more)
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
		return c.JSON(http.StatusBadRequest, "Maker Wallet does not match Wallet ID")
	}

	// Get inscriptions by inscription number and create a list of inscriptionIDs
	//for _, inscriptionNum := range tradeMakerReqBody.InscriptionNumbers {
	//	inscription, err := services.ORDEX.GetInscriptionByNumber(inscriptionNum)
	//	if err != nil {
	//		c.Logger().Error(err)
	//		return c.JSON(http.StatusInternalServerError, err.Error())
	//	}
	//	fmt.Printf("%+v \n", inscription)
	//}

	// Get maker inscriptions for trade, ensure maker owns those inscriptions, & append to maker side
	if err := parseMakerAssets(c, trade, maker, makerReqBody); err != nil {
		c.Logger().Error(err)
		return err
	}

	// update side
	if err := maker.Update(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := trade.SetStatus("OPEN"); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := trade.Update(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	trade.Maker = maker

	return c.JSON(http.StatusOK, trade)
}

/* Request Body
{
  "wallet_id": "bson_id",
  "btc": 1000,
  "inscription_numbers": [1234, 2344]
}
*/

func PostOfferByTradeID(c echo.Context) error {
	// Find & verify wallet
	// TODO: fix semantics
	makerReqBody := new(models.SideReqBody)
	if err := c.Bind(makerReqBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	wallet, err := models.GetWalletByID(makerReqBody.WalletID)
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Find & verify trade is in correct status
	trade, err := models.GetTradeByID(c, c.Param("id"))
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if trade.Status != "OPEN" {
		err := errors.New(fmt.Sprintf("trade.status is not OPEN, status is: %v", trade.Status))
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Create side & get maker inscriptions for offer, ensure maker owns those inscriptions, & append to maker side
	maker := models.NewSide(wallet.ID)
	maker.Wallet = wallet
	if err := parseMakerAssets(c, trade, maker, makerReqBody); err != nil {
		c.Logger().Error(err)
		return err
	}

	// Store offer make as a side
	if err := maker.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Create offer for trade
	offer := models.NewOffer(trade.ID)
	offer.MakerID = maker.ID
	if err := offer.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	offer.Maker = maker

	return c.JSON(http.StatusOK, offer)
}

/*
?status={enum,csv}
*/

func GetTrades(c echo.Context) error {
	// Get trades by query and return paginated response
	status := c.QueryParam("status")
	statusValues := strings.Split(status, ",")
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

func GetTradeByID(c echo.Context) error {
	// Get trade by tradeID and return response
	trade, err := models.GetTradeByID(c, c.Param("id"))
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, trade)
}

/* Query Params
?status={enum,csv}
*/

func GetOffersByTradeID(c echo.Context) error {
	// Get trade by ID to ensure trade exists
	trade, err := models.GetTradeByID(c, c.Param("id"))
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Get offers by Trade.ID and return response
	status := c.QueryParam("status")
	statusValues := strings.Split(status, ",")
	page, limit, err := parsePagination(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	offers, total, err := models.GetOffersPaginatedByTradeID(page, limit, trade.ID.Hex(), statusValues)
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

func parseMakerAssets(
	c echo.Context,
	trade *models.Trade,
	maker *models.Side,
	makerReqBody *models.SideReqBody,
) error {
	// TODO: cover wallets that have inscriptions greater then 100 (pagination)
	makerInscriptions, err := services.BESTINSLOT.GetInscriptionsByWalletAddr(
		c,
		maker.Wallet.TapRootAddr,
		100,
		1,
	)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	found := map[int64]bool{}
	for _, inscriptionNum := range makerReqBody.InscriptionNumbers {
		found[inscriptionNum] = false
		for _, makerInscription := range makerInscriptions.Data {
			if makerInscription.InscriptionNumber == inscriptionNum {
				found[inscriptionNum] = true
			}
		}
	}
	for k, v := range found {
		if v == false {
			err := errors.New(
				fmt.Sprintf("Inscription Number: %v not owned by Wallet: %v", k, trade.MakerID.Hex()),
			)
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}
	inscriptions := []*models.Inscription{}
	for _, makerInscription := range makerInscriptions.Data {
		inscription := models.ParseBISInscription(makerInscription)
		for _, incriptionNum := range makerReqBody.InscriptionNumbers {
			if incriptionNum == inscription.InscriptionNumber {
				inscriptions = append(inscriptions, inscription)
			}
		}
	}
	maker.Inscriptions = inscriptions
	maker.InscriptionNumbers = makerReqBody.InscriptionNumbers

	// Ensure maker has enough BTC for the offer
	// TODO: revisit to harden logic everywhere
	var makerPaymentAddr string
	if trade.Maker.Wallet.Type == "unisat" {
		makerPaymentAddr = trade.Maker.Wallet.TapRootAddr
	} else { // TODO: harden
		makerPaymentAddr = trade.Maker.Wallet.SegwitAddr
	}
	makerUTXOs, err := services.BLOCKCHAININFO.GetUTXOsForAddr(makerPaymentAddr)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var makerPaymentUTXOs []*models.UTXO
	for _, inscription := range maker.Inscriptions {
		inscriptionIDS := strings.Split(inscription.Satpoint, ":")
		if len(inscriptionIDS) != 3 {
			err := errors.New(
				fmt.Sprintf("error parsing paymentUTXOs for maker"),
			)
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		for _, utxoI := range makerUTXOs["unspent_outputs"].([]interface{}) {
			utxo := new(models.UTXO)
			if err := utxo.Parse(utxoI.(map[string]interface{})); err != nil {
				err := errors.New(
					fmt.Sprintf("could not parse utxo from blockchain info in data schema"),
				)
				c.Logger().Error(err)
				return c.JSON(http.StatusBadRequest, err.Error())
			}

			found := false
			if utxo.TxHashBigEndian == inscriptionIDS[0] {
				inscriptionIndex, err := strconv.Atoi(inscriptionIDS[1])
				if err != nil {
					err := errors.New(
						fmt.Sprintf("could not parse inscription index for maker"),
					)
					c.Logger().Error(err)
					return c.JSON(http.StatusBadRequest, err.Error())
				}
				if utxo.TxOutputN == int64(inscriptionIndex) {
					found = true
				}
			}
			if !found {
				makerPaymentUTXOs = append(makerPaymentUTXOs, utxo)
			}
		}
	}
	makerAvailableBTC := int64(0)
	for _, utxo := range makerPaymentUTXOs {
		makerAvailableBTC = makerAvailableBTC + utxo.Value
	}
	if makerAvailableBTC < makerReqBody.BTC {
		err := errors.New(fmt.Sprintf("maker does not have enough available BTC for trade"))
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	maker.BTC = makerReqBody.BTC

	//// TODO: query ordex for extra inscription information (floor price, previous tx, more...)
	//// look for another endpoint called get inscriptions by id (multiple same time)
	//response, err := services.ORDEX.GetInscriptionsByIds(tradeMakerReqBody.InscriptionNumbers)
	//if err != nil {
	//	c.Logger().Error(err)
	//	return c.JSON(http.StatusInternalServerError, err)
	//}
	//// fmt.Printf("inscriptions: %+v\n", inscriptions)
	//// Format response as readable JSON
	//formattedJSON, err := formatJSON(response)
	//if err != nil {
	//	fmt.Println("Error formatting JSON:", err)
	//	return c.JSON(http.StatusInternalServerError, err)
	//}
	//
	//fmt.Println("Formatted JSON:")
	//fmt.Println(formattedJSON)

	// for _, value := range tradeMakerReqBody.InscriptionNumbers {
	// 	inscription, err := services.BESTINSLOT.GetInscriptionById(c, value)
	// 	if err != nil {
	// 		c.Logger().Error(err)
	// 		return c.JSON(http.StatusInternalServerError, err)
	// 	}
	// 	fmt.Printf("inscription: %+v\n", inscription)
	// }

	return nil
}
