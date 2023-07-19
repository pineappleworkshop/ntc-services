package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
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

	url := fmt.Sprintf("https://blockchain.info/unspent?active=%+v", inscription.InscriberAddress)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var blockchainResponse models.BlockchainInfoResponse
	if err := json.Unmarshal(body, &blockchainResponse); err != nil {
		log.Fatalln(err)
	}

	data := services.InscriptionData{
		ContentType: inscription.Data.ContentType,
		Body:        []byte(inscription.Data.Body),
		Destination: inscription.Data.Destination,
	}

	var outpountList []*wire.OutPoint
	for _, op := range blockchainResponse.UnspentOutputs {
		hash, _ := chainhash.NewHashFromStr(op.TxHashBigEndian)
		outpoint := wire.OutPoint{
			Hash:  *hash,
			Index: uint32(op.TxOutputN),
		}
		outpountList = append(outpountList, &outpoint)
	}

	request := &services.InscriptionRequest{
		CommitTxOutPointList: outpountList,
		CommitFeeRate:        12,
		FeeRate:              12,
		DataList:             []services.InscriptionData{data},
		SingleRevealTxOnly:   true,
		RevealOutValue:       inscription.RevealOutValue,
	}

	inscriber, err := services.NewInscriber()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	_, reveal, _, err := inscriber.Inscribe(request)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ret := reveal[0].TxHash()

	return c.JSON(http.StatusOK, ret)
}
