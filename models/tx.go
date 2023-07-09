package models

import (
	"context"
	"github.com/btcsuite/btcd/btcjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/stores"
	"time"
)

type TxRaw struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id"`
	BlockID   primitive.ObjectID   `json:"block_id" bson:"block_id"`
	Hash      string               `json:"hash" bson:"hash"`
	Height    int64                `json:"height" bson:"height"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time           `json:"updated_at" bson:"updated_at"`
	TxRaw     *btcjson.TxRawResult `json:"tx_raw" bson:"tx_raw"`
	// TODO: Store and index all data points needed for certain operations
}

func NewTxRaw(blockID primitive.ObjectID, height int64, txRaw *btcjson.TxRawResult) *TxRaw {
	return &TxRaw{
		ID:        primitive.NewObjectID(),
		BlockID:   blockID,
		Hash:      txRaw.Hash,
		Height:    height,
		CreatedAt: time.Now().UTC(),
		TxRaw:     txRaw,
	}
}

func (tr *TxRaw) Save() error {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TXS_RAW)
	if _, err := collection.InsertOne(context.TODO(), tr, nil); err != nil {
		return err
	}

	return nil
}

// TODO: Save but check to prexisiting record
//
//
//type Tx struct {
//	Hex           string  `json:"hex" bson:"hex"`
//	TXID          string  `json:"txid" bson:"txid"`
//	Hash          string  `json:"hash" bson:"hash"`
//	Size          int64   `json:"size" bson:"size"`
//	VSize         int64   `json:"vsize" bson:"vsize"`
//	Version       int64   `json:"version" bson:"version"`
//	LockTime      int64   `json:"lockTime" bson:"lockTime"`
//	VIn           []*VIn  `json:"vin" bson:"vin"`
//	VOut          []*VOut `json:"vout" bson:"vout"`
//	BlockHash     string  `json:"blockhash" bson:"blockhash"`
//	Confirmations int64   `json:"confirmations" bson:"confirmations"`
//	Time          int64   `json:"time" bson:"time"`
//	BlockTime     int64   `json:"blocktime" bson:"blocktime"`
//}
//
//type VIn struct {
//	Sequence    int64      `json:"sequence" bson:"sequence"`
//	Coinbase    *bool      `json:"coinbase" bson:"coinbase"`
//	TXID        string     `json:"txid" bson:"txid"`
//	VOut        int64      `json:"vout" bson:"vout"`
//	ScriptSig   *ScriptSig `json:"script_sig" bson:"script_sig"`
//	TXINWitness []string   `json:"txinwitness" bson:"txinwitness"`
//}
//
//type ScriptSig struct {
//	ASM string `json:"asm" bson:"asm"`
//	Hex string `json:"hex" bson:"hex"`
//}
//
//type VOut struct {
//	Value        float64       `json:"value" bson:"value"`
//	N            int64         `json:"n" bson:"n"`
//	ScriptPubKey *ScriptPubKey `json:"script_pub_key" bson:"script_pub_key"`
//}
//
//type ScriptPubKey struct {
//	ASM       string   `json:"asm" bson:"asm"`
//	Hex       string   `json:"hex" bson:"hex"`
//	ReqSigs   *int32   `json:"req_sigs" bson:"req_sigs"`
//	Type      string   `json:"type" bson:"type"`
//	Addresses []string `json:"addresses" bson:"addresses"`
//	Address   string   `json:"address" json:"address"`
//}
