package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Side struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	WalletID     primitive.ObjectID `json:"wallet_id" bson:"wallet_id"`
	Wallet       *Wallet            `json:"wallet" bson:"-"`
	Inscriptions []*Inscription     `json:"inscriptions" bson:"inscriptions"`
	UTXOs        []*UTXO            `json:"utxos" bson:"utxos"`
	CreatedAt    int64              `json:"created_at" bson:"created_at"`
	UpdatedAt    *int64             `json:"updated_at" bson:"updated_at"`
}

func NewSide(walletID primitive.ObjectID) *Side {
	return &Side{
		ID:        primitive.NewObjectID(),
		WalletID:  walletID,
		CreatedAt: time.Now().Unix(),
	}
}
