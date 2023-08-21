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
  "wallet_id": "some_bson_id",
  "btc": 0,
  "inscription_numbers": [1234, 2344],
  "fee_rate"; 10
}
*/

func PostTrades(c echo.Context) error {
	// Verify request body
	tradeReqBody := new(models.TradeReqBody)
	if err := c.Bind(tradeReqBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Find & verify wallet
	wallet, err := models.GetWalletByID(tradeReqBody.WalletID)
	if err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Create side & store
	maker := models.NewSide(wallet.ID)
	if err := maker.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	maker.Wallet = wallet

	// Create trade & store
	trade := models.NewTrade(maker.ID, tradeReqBody.FeeRate)
	trade.Maker = maker
	if err := trade.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Get maker inscriptions for trade, ensure maker owns those inscriptions, & append to maker side
	if err := parseMakerAssets(c, trade, maker, tradeReqBody); err != nil {
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
	makerReqBody := new(models.TradeReqBody)
	if err := c.Bind(makerReqBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Find & verify wallet
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

	///////////////////////////////////////////////////////////////////////////////////

	// Get all inscriptions and utxos for trade.maker
	// TODO: check need to add checks against current trade
	//tradeMaker, err := models.GetSideByID(trade.MakerID.Hex())
	//if err != nil {
	//	if err.Error() == stores.MONGO_ERR_NOT_FOUND {
	//		c.Logger().Error(err)
	//		return c.JSON(http.StatusNotFound, err.Error())
	//	}
	//	c.Logger().Error(err)
	//	return c.JSON(http.StatusInternalServerError, err.Error())
	//}
	//allTradeMakerInscriptions, allTradeMakerUTXOs, err := parseMakerAssets(c, trade, tradeMaker, makerReqBody)
	//if err != nil {
	//	c.Logger().Error(err)
	//	return err
	//}

	//allTradeMakerInscriptionsBIS, err := services.BESTINSLOT.GetInscriptionsByWalletAddr(
	//	c,
	//	trade.Maker.Wallet.TapRootAddr,
	//	100,
	//	1,
	//)
	//if err != nil {
	//	c.Logger().Error(err)
	//	return c.JSON(http.StatusInternalServerError, err.Error())
	//}
	//var allTradeMakerInscriptions []*models.Inscription
	//for _, inscriptionBIS := range allTradeMakerInscriptionsBIS.Data {
	//	inscription := models.ParseBISInscription(inscriptionBIS)
	//	allTradeMakerInscriptions = append(allTradeMakerInscriptions, inscription)
	//}
	//
	//var makerPaymentAddr string
	//if trade.Maker.Wallet.Type == "unisat" {
	//	makerPaymentAddr = trade.Maker.Wallet.TapRootAddr
	//} else { // TODO: harden
	//	makerPaymentAddr = trade.Maker.Wallet.SegwitAddr
	//}
	//allTradeMakerUTXOIs, err := services.BLOCKCHAININFO.GetUTXOsForAddr(makerPaymentAddr)
	//if err != nil {
	//	c.Logger().Error(err)
	//	return c.JSON(http.StatusInternalServerError, err.Error())
	//}
	//var allTradeMakerUTXOs []*models.UTXO
	//for _, utxoI := range allTradeMakerUTXOIs["unspent_outputs"].([]interface{}) {
	//	utxo := new(models.UTXO)
	//	if err := utxo.Parse(utxoI.(map[string]interface{})); err != nil {
	//		err := errors.New(
	//			fmt.Sprintf("could not parse utxo from blockchain info in data schema"),
	//		)
	//		c.Logger().Error(err)
	//		return c.JSON(http.StatusInternalServerError, err.Error())
	//	}
	//	allTradeMakerUTXOs = append(allTradeMakerUTXOs, utxo)
	//}

	if err := getWalletAssetsByID(c, trade.Maker.Wallet); err != nil {
		c.Logger().Error(err)
		return c.JSON(c.Response().Status)
	}

	///////////////////////////////////////////////////////////////////////////////////

	// Create side & get maker inscriptions for offer, ensure maker owns those inscriptions, & append to maker side
	maker := models.NewSide(wallet.ID)
	maker.Wallet = wallet

	if err := parseMakerAssets(c, trade, maker, makerReqBody); err != nil {
		c.Logger().Error(err)
		return err
	}

	if err := getWalletAssetsByID(c, maker.Wallet); err != nil {
		c.Logger().Error(err)
		return err
	}

	// Store offer maker as a side
	if err := maker.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Create offer for trade and add to trade for parsing psbt
	offer := models.NewOffer(trade.ID)
	offer.MakerID = maker.ID
	trade.Taker = maker
	trade.Taker.Inscriptions = maker.Inscriptions

	fmt.Println("**************************")
	fmt.Printf("allTradeMakerUTXOs: %+v \n", trade.Maker.Wallet.UTXOs)
	fmt.Printf("allOfferMakerUTXOs: %+v \n", maker.Wallet.UTXOs)
	fmt.Printf("allTradeMakerInscriptions: %+v \n", trade.Maker.Wallet.Inscriptions)
	fmt.Printf("allOfferMakerInscriptions: %+v \n", maker.Wallet.Inscriptions)
	fmt.Println("**************************")

	// Create PSBT
	psbtService := services.NewPBST(
		trade,
		trade.Maker.Wallet.UTXOs,
		maker.Wallet.UTXOs,
		trade.Maker.Wallet.Inscriptions,
		maker.Wallet.Inscriptions,
	)
	if err := psbtService.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	psbt, err := psbtService.GeneratePSBT(c)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	offer.PSBT = psbt

	// Store offer
	if err := offer.Create(c); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	offer.Maker = maker

	return c.JSON(http.StatusOK, offer)
}

/* Query Params
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

/* Request Body
{
  "wallet_id": "bson_id",
  "btc": 1000,
  "inscription_numbers": [1234, 2344]
}
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

// TODO: need to fix some logical shit here
func parseMakerAssets(
	c echo.Context,
	trade *models.Trade,
	maker *models.Side,
	makerReqBody *models.TradeReqBody,
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

	var allInscriptions []*models.Inscription
	var sideInscriptions []*models.Inscription
	for _, makerInscription := range makerInscriptions.Data {
		inscription := models.ParseBISInscription(makerInscription)
		for _, inscriptionNum := range makerReqBody.InscriptionNumbers {
			if inscriptionNum == inscription.InscriptionNumber {
				sideInscriptions = append(sideInscriptions, inscription)
			}
			allInscriptions = append(allInscriptions, inscription)
		}
	}
	maker.Inscriptions = sideInscriptions
	maker.InscriptionNumbers = makerReqBody.InscriptionNumbers

	// Ensure maker has enough BTC for the offer
	// TODO: revisit to harden logic everywhere
	var makerPaymentAddr string
	if maker.Wallet.Type == "unisat" {
		makerPaymentAddr = maker.Wallet.TapRootAddr
	} else { // TODO: harden
		makerPaymentAddr = maker.Wallet.SegwitAddr
	}
	makerUTXOs, err := services.BLOCKCHAININFO.GetUTXOsForAddr(makerPaymentAddr)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// TODO: there's something fucked up with this logic
	var allMakerUTXOs []*models.UTXO
	var makerPaymentUTXOs []*models.UTXO
	for _, utxoI := range makerUTXOs["unspent_outputs"].([]interface{}) {
		utxo := new(models.UTXO)
		if err := utxo.Parse(utxoI.(map[string]interface{})); err != nil {
			err := errors.New(
				fmt.Sprintf("could not parse utxo from blockchain info in data schema"),
			)
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		allMakerUTXOs = append(allMakerUTXOs, utxo)

		found := false
		if len(maker.Inscriptions) > 0 {
			for _, inscription := range allInscriptions {
				inscriptionIDS := strings.Split(inscription.Satpoint, ":")
				if len(inscriptionIDS) != 3 {
					err := errors.New(
						fmt.Sprintf("error parsing paymentUTXOs for maker"),
					)
					c.Logger().Error(err)
					return c.JSON(http.StatusBadRequest, err.Error())
				}
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
		} else {
			makerPaymentUTXOs = append(makerPaymentUTXOs, utxo)
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
