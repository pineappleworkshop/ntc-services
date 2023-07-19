package handlers

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/models"
	"ntc-services/services"
)

func Inscribe(c echo.Context) error {
	var inscriptionI interface{}
	if err := c.Bind(&inscriptionI); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	inscription := models.NewInscription()
	if err := inscription.Parse(inscriptionI.(map[string]interface{})); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	hash, err := chainhash.NewHashFromStr(inscription.CommitTxOutPoint.Hash)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	outpoint := &wire.OutPoint{
		Hash:  *hash,
		Index: inscription.CommitTxOutPoint.Index,
	}

	data := services.InscriptionData{
		ContentType: inscription.Data.ContentType,
		Body:        []byte(inscription.Data.Body),
		Destination: inscription.Data.Destination,
	}

	request := &services.InscriptionRequest{
		CommitTxOutPointList:   []*wire.OutPoint{outpoint},
		CommitTxPrivateKeyList: nil,
		CommitFeeRate:          inscription.CommitFeeRate,
		FeeRate:                inscription.FeeRate,
		DataList:               []services.InscriptionData{data},
		SingleRevealTxOnly:     true,
		RevealOutValue:         inscription.RevealOutValue,
	}

	inscriber, err := services.NewInscriber()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	commit, reveals, fees, err := inscriber.Inscribe(request)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := models.InscriptionResp{
		CommitTxHash:     commit,
		RevealTxHashList: reveals,
		Fees:             fees,
	}

	return c.JSON(http.StatusNotImplemented, resp)
}
