package models

import (
	"context"
	"ntc-services/stores"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Offers struct {
	Page   int64    `json:"page"`
	Limit  int64    `json:"limit"`
	Total  int64    `json:"total"`
	Offers []*Offer `json:"offers"`
}

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
		ID:        primitive.NewObjectID(),
		TradeID:   tradeID,
		CreatedAt: time.Now().Unix(),
	}
}

func (o *Offer) Create(c echo.Context) error {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_OFFERS)
	if _, err := collection.InsertOne(context.TODO(), o); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

func (o *Offer) Update(c echo.Context) error {
	now := time.Now().Unix()
	o.UpdatedAt = &now

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_OFFERS)
	if _, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": o.ID}, o); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

func GetOffersByTradeID(c echo.Context) ([]*Offer, error) {
	tradeID := c.Param("id")
	idHex, err := primitive.ObjectIDFromHex(tradeID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"trade_id": idHex}

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_OFFERS)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}

	var offers []*Offer
	if err := cursor.All(context.TODO(), &offers); err != nil {
		c.Logger().Error(err)
		return nil, err
	}

	return offers, nil
}

func GetOffersPaginatedByTradeID(page, limit int64, tradeId string, status []string) ([]*Offer, int64, error) {
	opts := options.Find().SetLimit(limit).SetSkip(page - 1)
	var filter bson.M
	if len(status) > 0 && status[0] != "" {
		filter = bson.M{"status": bson.M{"$in": status}}
	} else {
		filter = bson.M{}
	}

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_OFFERS)

	total, err := collection.CountDocuments(context.TODO(), filter, nil)
	if err != nil {
		return nil, -1, err
	}

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, -1, err
	}
	var offers []*Offer
	if err := cursor.All(context.TODO(), &offers); err != nil {
		return nil, -1, err
	}

	return offers, total, nil
}
