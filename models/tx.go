package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tx struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	BlockRawID primitive.ObjectID `json:"block_raw_id" bson:"block_raw_id"`
	TxRawID    primitive.ObjectID `json:"tx_raw_id" bson:"tx_raw_id"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time         `json:"updated_at" bson:"updated_at"`
}

type Vin struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	TxID      primitive.ObjectID `json:"tx_id" bson:"tx_id"`
	Index     int64              `json:"index" bson:"index"`
	Coinbase  string             `json:"coinbase" bson:"coinbase"`
	Txid      string             `json:"txid" bson:"txid"`
	Vout      uint32             `json:"vout" bson:"vout"`
	ScriptSig *ScriptSig         `json:"script_sig" bson:"script_sig"`
	Sequence  uint32             `json:"sequence" bson:"sequence"`
	Witness   []string           `json:"txinwitness" bson:"txinwitness"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at"`
}

type Vout struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	TxID         primitive.ObjectID `json:"tx_id" bson:"tx_id"`
	Index        int64              `json:"index" bson:"index"`
	Value        float64            `json:"value" bson:"value"`
	N            uint32             `json:"n" bson:"n"`
	ScriptPubKey ScriptPubKeyResult `json:"script_pub_key" bson:"script_pub_key"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    *time.Time         `json:"updated_at" bson:"updated_at"`
}

type ScriptSig struct {
	Asm string `json:"asm" bson:"asm"`
	Hex string `json:"hex" bson:"hex"`
}

type ScriptPubKeyResult struct {
	Asm       string   `json:"asm" bson:"asm"`
	Hex       string   `json:"hex" bson:"hex"`
	ReqSigs   int32    `json:"req_sigs" bson:"req_sigs"`
	Type      string   `json:"type" bson:"type"`
	Addresses []string `json:"addresses" bson:"addresses"`
}
