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

type TradeReqBody struct {
	WalletType  string `json:"wallet_type" bson:"wallet_type"`
	TapRootAddr string `json:"tap_root_addr" bson:"tap_root_addr"`
	SegwitAddr  string `json:"segwit_addr" bson:"segwit_addr"`
}

type TradeMakerReqBody struct {
	WalletID           string   `json:"wallet_id" bson:"wallet_id"`
	Btc                int64    `json:"btc" bson:"btc"`
	InscriptionNumbers []string `json:"inscription_numbers" bson:"inscription_numbers"`
}

type Trades struct {
	Page   int64    `json:"page"`
	Limit  int64    `json:"limit"`
	Total  int64    `json:"total"`
	Trades []*Trade `json:"trades"`
}
type Trade struct {
	ID              primitive.ObjectID  `json:"id" bson:"_id"`
	MakerID         primitive.ObjectID  `json:"maker_id" bson:"maker_id"`
	Maker           *Side               `json:"maker" bson:"-"`
	TakerID         *primitive.ObjectID `json:"taker_id" bson:"taker_id"`
	Taker           *Side               `json:"taker" bson:"-"`
	PSBT            *string             `json:"psbt" bson:"psbt"`
	FeeRate         *FeeRate            `json:"fee_rate" bson:"fee_rate"`
	PlatformFee     *int64              `json:"platform_fee" bson:"platform_fee"`
	TxID            *string             `json:"txId" bson:"txId"`
	Status          string              `json:"status" bson:"status"`
	StatusChangedAt *int64              `json:"status_changed_at" bson:"status_changed_at"`
	CreatedAt       int64               `json:"created_at" bson:"created_at"`
	UpdatedAt       *int64              `json:"updated_at" bson:"updated_at"`
}

type FeeRate struct {
	Label string `json:"label" bson:"label"`
	Value int32  `json:"value" bson:"value"`
}

func NewTradeReqBody() *TradeReqBody {

	return &TradeReqBody{}
}
func NewTradeMakerReqBody() *TradeMakerReqBody {

	return &TradeMakerReqBody{}
}
func NewTrade(makerID primitive.ObjectID) *Trade {
	return &Trade{
		ID:        primitive.NewObjectID(),
		MakerID:   makerID,
		Status:    "CREATED",
		CreatedAt: time.Now().Unix(),
	}
}

func (t *Trade) Create(c echo.Context) error {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)
	if _, err := collection.InsertOne(context.TODO(), t); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

func (t *Trade) Update(c echo.Context) error {
	now := time.Now().Unix()
	t.UpdatedAt = &now

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)
	if _, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": t.ID}, t); err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

func GetTradesByStatus(c echo.Context, status string) ([]*Trade, error) {
	filter := bson.M{"status": status}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}

	var trades []*Trade
	if err := cursor.All(context.TODO(), &trades); err != nil {
		c.Logger().Error(err)
		return nil, err
	}

	return trades, nil
}

func GetTradesPaginatedByStatus(page, limit int64, status []string) ([]*Trade, int64, error) {
	opts := options.Find().SetLimit(limit).SetSkip(page - 1)

	var filter bson.M
	if len(status) > 0 && status[0] != "" {
		filter = bson.M{"status": bson.M{"$in": status}}
	} else {
		filter = bson.M{}
	}

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)

	total, err := collection.CountDocuments(context.TODO(), filter, nil)
	if err != nil {
		return nil, -1, err
	}

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, -1, err
	}
	var trades []*Trade
	if err := cursor.All(context.TODO(), &trades); err != nil {
		return nil, -1, err
	}

	return trades, total, nil
}

func GetTradeByID(c echo.Context, idStr string) (*Trade, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	filter := bson.M{"_id": id}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_TRADES)

	var trade *Trade
	if err := collection.FindOne(context.TODO(), filter).Decode(&trade); err != nil {
		c.Logger().Error(err)
		return nil, err
	}

	return trade, nil
}
