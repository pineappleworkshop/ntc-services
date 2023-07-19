package handlers

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/models"
	"ntc-services/pkg/btcapi/mempool"
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

	netParams := &chaincfg.MainNetParams
	btcApiClient := mempool.NewClient(netParams)

	utxoPrivateKeyHex := ""
	//destination := "bc1p7ncck66wthnjl2clcry46f2uxjcn8naw95e6r8ag0x9zremx00lqvf5wve"

	commitTxOutPointList := make([]*wire.OutPoint, 0)
	commitTxPrivateKeyList := make([]*btcec.PrivateKey, 0)

	{
		utxoPrivateKeyBytes, err := hex.DecodeString(utxoPrivateKeyHex)
		if err != nil {
			fmt.Println(err)
		}
		utxoPrivateKey, _ := btcec.PrivKeyFromBytes(utxoPrivateKeyBytes)

		utxoTaprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(utxoPrivateKey.PubKey())), netParams)
		if err != nil {
			fmt.Println(err)
		}

		unspentList, err := btcApiClient.ListUnspent(utxoTaprootAddress)
		if err != nil {
			fmt.Println(err)
		}

		for i := range unspentList {
			commitTxOutPointList = append(commitTxOutPointList, unspentList[i].Outpoint)
			commitTxPrivateKeyList = append(commitTxPrivateKeyList, utxoPrivateKey)
		}
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

	fmt.Println(commitTxOutPointList)

	request := &services.InscriptionRequest{
		CommitTxOutPointList:   commitTxOutPointList,
		CommitTxPrivateKeyList: commitTxPrivateKeyList,
		CommitFeeRate:          12,
		FeeRate:                12,
		DataList:               []services.InscriptionData{data},
		SingleRevealTxOnly:     true,
		RevealOutValue:         inscription.RevealOutValue,
	}

	fmt.Print(inscription.RevealOutValue)
	fmt.Print(data)
	fmt.Print(outpoint)

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

	return c.JSON(http.StatusOK, resp)
}
