package models

import (
	"context"
	"fmt"
	"ntc-services/stores"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Wallet struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	Type               string             `json:"wallet_type" bson:"wallet_type"`
	CardinalAddr       string             `json:"cardinal_addr" bson:"cardinal_addr"`
	TapRootAddr        string             `json:"tap_root_addr" bson:"tap_root_addr"`
	SegwitAddr         string             `json:"segwit_addr" bson:"segwit_addr"`
	CreatedAt          int64              `json:"created_at" bson:"created_at"`
	UpdatedAt          *int64             `json:"updated_at" bson:"updated_at"`
	LastConnectedAt    *int64             `json:"last_connected_at" bson:"last_connected_at"`
	LastConnectedBlock *int64             `json:"last_connected_block" bson:"last_connected_block"`
}

func NewWallet() *Wallet {

	return &Wallet{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now().Unix(),
	}
}

func IsValidWalletType(walletType string) bool {
	switch walletType {
	case "hiro", "xverse", "unisat":
		return true
	default:
		return false
	}
}

func (w *Wallet) Save() error {
	client := stores.DB.Mongo.Client
	collection := client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_WALLETS)
	if _, err := collection.InsertOne(context.TODO(), w); err != nil {
		return err
	}

	return nil
}

func (w *Wallet) Update() error {
	now := time.Now().UTC().Unix()

	w.UpdatedAt = &now
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_WALLETS)
	if _, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": w.ID}, w); err != nil {
		return err
	}

	return nil
}

func GetWalletByID(id string) (*Wallet, error) {
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_WALLETS)
	filter := bson.M{"_id": idHex}
	result := collection.FindOne(context.TODO(), filter)
	var wallet *Wallet
	if err := result.Decode(&wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func GetWalletByAddr(tapRootAddr, segwitAddr, addrType string) (*Wallet, error) {
	filter := bson.M{}
	switch addrType {
	case "hiro":
		filter["tap_root_addr"] = tapRootAddr
		filter["segwit_addr"] = segwitAddr
	case "xverse":
		filter["tap_root_addr"] = tapRootAddr
		filter["segwit_addr"] = segwitAddr
	case "unisat":
		filter["segwit_addr"] = segwitAddr
	default:
		return nil, fmt.Errorf("invalid address type: %s", addrType)
	}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_WALLETS)
	result := collection.FindOne(context.TODO(), filter)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil // No wallet found
	} else if result.Err() != nil {
		return nil, result.Err() // Error occurred
	}
	var wallet Wallet
	if err := result.Decode(&wallet); err != nil {
		return nil, err
	}

	return &wallet, nil
}
