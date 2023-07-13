package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/stores"
	"time"
)

type Order struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	MakerAddress    string             `json:"makerAddress" bson:"makerAddress"`
	TakerAddress    string             `json:"takerAddress" bson:"takerAddress"`
	Status          string             `json:"status" bson:"status"`
	MakerSelections []interface{}      `json:"makerSelections" bson:"makerSelections"`
	TakerSelections []interface{}      `json:"takerSelections" bson:"takerSelections"`
	FeeRate         interface{}        `json:"feeRate" bson:"feeRate"`
	TxID            string             `json:"txId" bson:"txId"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       *time.Time         `json:"updated_at" bson:"updated_at"`
}

func (o *Order) Update() error {
	now := time.Now().UTC()
	o.UpdatedAt = &now

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)
	if _, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": o.ID}, o); err != nil {
		return err
	}

	return nil
}
