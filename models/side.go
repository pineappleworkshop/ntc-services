package models

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/stores"
	"time"
)

type Side struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	WalletID           primitive.ObjectID `json:"wallet_id" bson:"wallet_id"`
	Wallet             *Wallet            `json:"wallet" bson:"-"`
	InscriptionNumbers []int64            `json:"inscription_numbers" bson:"inscription_numbers"`
	BTC                int64              `json:"btc" bson:"btc"`
	Inscriptions       []*Inscription     `json:"inscriptions" bson:"inscriptions"` // inscriptions in side
	UTXOs              []*UTXO            `json:"utxos" bson:"utxos"`               // maker will not select UTXOs
	CreatedAt          int64              `json:"created_at" bson:"created_at"`
	UpdatedAt          *int64             `json:"updated_at" bson:"updated_at"`
}

func NewSide(walletID primitive.ObjectID) *Side {
	return &Side{
		ID:        primitive.NewObjectID(),
		WalletID:  walletID,
		CreatedAt: time.Now().Unix(),
	}
}

func (s *Side) Create(c echo.Context) error {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_SIDES)
	if _, err := collection.InsertOne(context.TODO(), s); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

func (s *Side) Update(c echo.Context) error {
	now := time.Now().Unix()
	s.UpdatedAt = &now

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_SIDES)
	if _, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": s.ID}, s); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}
