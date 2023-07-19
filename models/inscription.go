package models

import (
	"github.com/btcsuite/btcd/wire"
)

type Inscription struct {
	//CommitTxOutPoint *OutPoint        `json:"commitTxOutPoint"`
	InscriberAddress string           `json:"inscriberAddress"`
	CommitFeeRate    int64            `json:"commitFeeRate"`
	FeeRate          int64            `json:"feeRate"`
	Data             *InscriptionData `json:"data"`
	RevealOutValue   int64            `json:"revealOutValue"`
}

type InscriptionData struct {
	ContentType string `json:"contentType"`
	Body        string `json:"body"`
	Destination string `json:"destination"`
}

type OutPoint struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}

type InscriptionResp struct {
	CommitTxHash     *wire.MsgTx   `json:"commitTxHash"`
	RevealTxHashList []*wire.MsgTx `json:"revealTxHashList"`
	Fees             int64         `json:"fees"`
}

func NewInscription() *Inscription {
	//commitTxOutPoint := OutPoint{}
	data := InscriptionData{}
	return &Inscription{
		//CommitTxOutPoint: &commitTxOutPoint,
		Data: &data,
	}
}

func (i *Inscription) Parse(body map[string]interface{}) error {
	//if body["commitTxOutPoint"] != nil {
	//	if body["commitTxOutPoint"].(map[string]interface{})["hash"] != nil {
	//		i.CommitTxOutPoint.Hash = body["commitTxOutPoint"].(map[string]interface{})["hash"].(string)
	//	}
	//	if body["commitTxOutPoint"].(map[string]interface{})["index"] != nil {
	//		i.CommitTxOutPoint.Index = uint32(body["commitTxOutPoint"].(map[string]interface{})["index"].(float64))
	//	}
	//}
	if body["commitFeeRate"] != nil {
		i.CommitFeeRate = int64(body["commitFeeRate"].(float64))
	}
	if body["feeRate"] != nil {
		i.CommitFeeRate = int64(body["feeRate"].(float64))
	}
	if body["data"] != nil {
		if body["data"].(map[string]interface{})["contentType"] != nil {
			i.Data.ContentType = body["data"].(map[string]interface{})["contentType"].(string)
		}
		if body["data"].(map[string]interface{})["body"] != nil {
			i.Data.Body = body["data"].(map[string]interface{})["body"].(string)
		}
		if body["data"].(map[string]interface{})["destination"] != nil {
			i.Data.Destination = body["data"].(map[string]interface{})["destination"].(string)
		}
	}
	if body["revealOutValue"] != nil {
		i.RevealOutValue = int64(body["revealOutValue"].(float64))
	}

	return nil
}

type UnspentOutput struct {
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxHash          string `json:"tx_hash"`
	TxOutputN       int    `json:"tx_output_n"`
	Script          string `json:"script"`
	Value           int    `json:"value"`
	ValueHex        string `json:"value_hex"`
	Confirmations   int    `json:"confirmations"`
	TxIndex         int64  `json:"tx_index"`
}

type BlockchainInfoResponse struct {
	Notice         string          `json:"notice"`
	UnspentOutputs []UnspentOutput `json:"unspent_outputs"`
}
