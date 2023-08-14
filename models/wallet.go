package models

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/stores"
	"time"
)

type Wallet struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Type         string             `json:"wallet_type" bson:"wallet_type"`
	CardinalAddr string             `json:"cardinal_addr" bson:"cardinal_addr"`
	TapRootAddr  string             `json:"tap_root_addr" bson:"tap_root_addr"`
	SegwitAddr   string             `json:"segwit_addr" bson:"segwit_addr"`
	CreatedAt    int64              `json:"created_at" bson:"created_at"`
	UpdatedAt    *int64             `json:"updated_at" bson:"updated_at"`
}

func NewWallet(walletType string) *Wallet {
	// TODO: validate wallet type against enum

	return &Wallet{
		ID:   primitive.NewObjectID(),
		Type: walletType,
		//CardinalAddr: cardinalAddr,
		//TapRootAddr:  tapRootAddr,
		//SegwitAddr:   "",
		CreatedAt: time.Now().Unix(),
	}
}

func (w *Wallet) Create(c echo.Context) error {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_WALLETS)
	if _, err := collection.InsertOne(context.TODO(), w); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}
