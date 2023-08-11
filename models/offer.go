package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Offer struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	TradeID         primitive.ObjectID `json:"trade_id" bson:"trade_id"`
	MakerID         primitive.ObjectID `json:"maker_id" bson:"maker_id"`
	Maker           *Side              `json:"maker" bson:"maker"`
	Status          string             `json:"status" bson:"status"`
	StatusChangedAt *int64             `json:"status_changed_at" bson:"status_changed_at"`
	CreatedAt       int64              `json:"created_at" bson:"created_at"`
	UpdatedAt       *int64             `json:"updated_at" bson:"updated_at"`
}

func NewOffer(tradeID primitive.ObjectID) *Offer {
	return &Offer{
		ID:        primitive.ObjectID{},
		TradeID:   tradeID,
		CreatedAt: time.Now().Unix(),
	}
}
