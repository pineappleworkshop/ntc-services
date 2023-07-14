package models

import (
	"context"
	"github.com/btcsuite/btcd/btcjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/stores"
	"time"
)

type Tx struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	BlockRawID    primitive.ObjectID `json:"block_raw_id" bson:"block_raw_id"`
	TxRawID       primitive.ObjectID `json:"tx_raw_id" bson:"tx_raw_id"`
	BlockHeight   int64              `json:"block_height" bson:"block_height"`
	Height        int64              `json:"height" bson:"height"`
	Hex           string             `json:"hex" bson:"hex"`
	Txid          string             `json:"txid" bson:"txid"`
	Hash          string             `json:"hash" bson:"hash"`
	Size          int32              `json:"size" bson:"size"`
	Vsize         int32              `json:"vsize" bson:"vsize"`
	Weight        int32              `json:"weight" bson:"weight"`
	Version       uint32             `json:"version" bson:"version"`
	LockTime      uint32             `json:"locktime" bson:"locktime"`
	BlockHash     string             `json:"blockhash" bson:"blockhash"`
	Confirmations uint64             `json:"confirmations" bson:"confirmations"`
	Time          int64              `json:"time" bson:"time"`
	Blocktime     int64              `json:"blocktime" bson:"blocktime"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     *time.Time         `json:"updated_at" bson:"updated_at"`
	//Vin           []Vin  `json:"vin"`
	//Vout          []Vout `json:"vout"`
}

func NewTx(blockRawID, txRawID primitive.ObjectID) *Tx {
	return &Tx{
		ID:         primitive.NewObjectID(),
		BlockRawID: blockRawID,
		TxRawID:    txRawID,
		CreatedAt:  time.Now().UTC(),
	}
}

func (t *Tx) Parse(blockHeight int64, txRaw *TxRaw) error {
	t.BlockHeight = blockHeight
	t.Height = txRaw.Height
	t.Hex = txRaw.TxRaw.Hex
	t.Txid = txRaw.TxRaw.Txid
	t.Hash = txRaw.TxRaw.Hash
	t.Size = txRaw.TxRaw.Size
	t.Vsize = txRaw.TxRaw.Vsize
	t.Weight = txRaw.TxRaw.Weight
	t.Version = txRaw.TxRaw.Version
	t.LockTime = txRaw.TxRaw.LockTime
	t.BlockHash = txRaw.TxRaw.BlockHash
	t.Confirmations = txRaw.TxRaw.Confirmations
	t.Time = txRaw.TxRaw.Time
	t.Blocktime = txRaw.TxRaw.Blocktime

	return nil
}

func (t *Tx) Save() error {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TXS)
	if _, err := collection.InsertOne(context.TODO(), t, nil); err != nil {
		return err
	}

	return nil
}

type Vin struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	TxID       primitive.ObjectID `json:"tx_id" bson:"tx_id"`
	BlockRawID primitive.ObjectID `json:"block_raw_id" bson:"block_raw_id"`
	TxRawID    primitive.ObjectID `json:"tx_raw_id" bson:"tx_raw_id"`
	Index      int64              `json:"index" bson:"index"`
	Coinbase   string             `json:"coinbase" bson:"coinbase"`
	Txid       string             `json:"txid" bson:"txid"`
	Vout       uint32             `json:"vout" bson:"vout"`
	ScriptSig  *ScriptSig         `json:"script_sig" bson:"script_sig"`
	Sequence   uint32             `json:"sequence" bson:"sequence"`
	Witness    []string           `json:"txinwitness" bson:"txinwitness"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time         `json:"updated_at" bson:"updated_at"`
}

func NewVin(txID, blockRawID, txRawID primitive.ObjectID) *Vin {
	return &Vin{
		ID:         primitive.NewObjectID(),
		TxID:       txID,
		BlockRawID: blockRawID,
		TxRawID:    txRawID,
		CreatedAt:  time.Now().UTC(),
	}
}

func (v *Vin) Parse(index int64, vin btcjson.Vin) error {
	v.Index = index
	v.Coinbase = vin.Coinbase
	v.Txid = vin.Txid
	v.Vout = vin.Vout
	v.Sequence = vin.Sequence
	v.Witness = vin.Witness
	if vin.ScriptSig != nil {
		v.ScriptSig = &ScriptSig{
			Asm: vin.ScriptSig.Asm,
			Hex: vin.ScriptSig.Hex,
		}
	}

	return nil
}

func (v *Vin) Save() error {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_VINS)
	if _, err := collection.InsertOne(context.TODO(), v, nil); err != nil {
		return err
	}

	return nil
}

type Vout struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	TxID         primitive.ObjectID `json:"tx_id" bson:"tx_id"`
	BlockRawID   primitive.ObjectID `json:"block_raw_id" bson:"block_raw_id"`
	TxRawID      primitive.ObjectID `json:"tx_raw_id" bson:"tx_raw_id"`
	Index        int64              `json:"index" bson:"index"`
	Value        float64            `json:"value" bson:"value"`
	N            uint32             `json:"n" bson:"n"`
	ScriptPubKey ScriptPubKey       `json:"script_pub_key" bson:"script_pub_key"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    *time.Time         `json:"updated_at" bson:"updated_at"`
}

func NewVout(txID, blockRawID, txRawID primitive.ObjectID) *Vout {
	return &Vout{
		ID:         primitive.NewObjectID(),
		TxID:       txID,
		BlockRawID: blockRawID,
		TxRawID:    txRawID,
		CreatedAt:  time.Now().UTC(),
	}
}

func (v *Vout) Parse(index int64, vout btcjson.Vout) error {
	v.Index = index
	v.Value = vout.Value
	v.N = vout.N
	v.ScriptPubKey = ScriptPubKey{
		Asm:       vout.ScriptPubKey.Asm,
		Hex:       vout.ScriptPubKey.Hex,
		ReqSigs:   vout.ScriptPubKey.ReqSigs,
		Type:      vout.ScriptPubKey.Type,
		Addresses: vout.ScriptPubKey.Addresses,
	}

	return nil
}

func (v *Vout) Save() error {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_VOUTS)
	if _, err := collection.InsertOne(context.TODO(), v, nil); err != nil {
		return err
	}

	return nil
}

type ScriptSig struct {
	Asm string `json:"asm" bson:"asm"`
	Hex string `json:"hex" bson:"hex"`
}

type ScriptPubKey struct {
	Asm       string   `json:"asm" bson:"asm"`
	Hex       string   `json:"hex" bson:"hex"`
	ReqSigs   int32    `json:"req_sigs" bson:"req_sigs"`
	Type      string   `json:"type" bson:"type"`
	Addresses []string `json:"addresses" bson:"addresses"`
}
