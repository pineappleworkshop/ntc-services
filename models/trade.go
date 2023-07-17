package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/stores"
	"time"
)

type Trade struct {
	ID                    primitive.ObjectID `json:"id" bson:"_id"`
	MakerAddress          string             `json:"makerAddress" bson:"makerAddress"`
	TakerAddress          string             `json:"takerAddress" bson:"takerAddress"`
	Status                string             `json:"status" bson:"status"`
	PBST                  string             `json:"pbst" bson:"pbst"`
	MakerSelections       []interface{}      `json:"makerSelections" bson:"makerSelections"`
	TakerSelections       []interface{}      `json:"takerSelections" bson:"takerSelections"`
	MakerUninscribedUtxos []interface{}      `json:"makerUninscribedUtxos" bson:"makerUninscribedUtxos"`
	TakerUninscribedUtxos []interface{}      `json:"takerUninscribedUtxos" bson:"takerUninscribedUtxos"`
	FeeRate               interface{}        `json:"feeRate" bson:"feeRate"`
	TxID                  *string            `json:"txId" bson:"txId"`
	CreatedAt             int64              `json:"createdAt" bson:"createdAt"`
	UpdatedAt             *int64             `json:"updatedAt" bson:"updatedAt"`
	//Confirmations   int64              `json:"confirmations" bson:"confirmations"`
}

func (t *Trade) Update() error {
	now := time.Now().Unix()
	t.UpdatedAt = &now

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)
	if _, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": t.ID}, t); err != nil {
		return err
	}

	return nil
}

func GetTradesByStatus(status string) ([]*Trade, error) {
	filter := bson.M{"status": status}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var trades []*Trade
	if err := cursor.All(context.TODO(), &trades); err != nil {
		return nil, err
	}

	return trades, nil
}

func GetTradeByID(idStr string) (*Trade, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": id}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)

	var trade *Trade
	if err := collection.FindOne(context.TODO(), filter).Decode(&trade); err != nil {
		return nil, err
	}

	return trade, nil
}
