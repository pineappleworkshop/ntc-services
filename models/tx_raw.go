package models

import (
	"context"
	"github.com/btcsuite/btcd/btcjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/stores"
	"time"
)

type TxRaw struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id"`
	BlockID   primitive.ObjectID   `json:"block_id" bson:"block_id"`
	TxID      string               `json:"tx_id" bson:"tx_id"` // TODO: change to txid, audit data
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
		TxID:      txRaw.Txid,
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

func SaveRawTxs(txRaws []*TxRaw) error {
	documents := make([]interface{}, len(txRaws))
	for i, txRaw := range txRaws {
		documents[i] = txRaw
	}

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TXS_RAW)
	if _, err := collection.InsertMany(context.Background(), documents); err != nil {
		return err
	}

	return nil
}
func GetTxRawsByBlockID(idStr string) ([]*TxRaw, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"block_id": id}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TXS_RAW)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var blockRaws []*TxRaw
	if err := cursor.All(context.TODO(), &blockRaws); err != nil {
		return nil, err
	}

	return blockRaws, nil
}
