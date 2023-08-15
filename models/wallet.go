package models

import (
	"context"
	"fmt"
	"ntc-services/stores"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Address struct {
	Addr     string `json:"addr" bson:"addr"`
	AddrType string `json:"addr_type" bson:"addr_type"`
}

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

func NewAddr() *Address {

	return &Address{}
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

func GetWalletBy(id string) (*Wallet, error) {
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

func GetWalletByAddr(addr, addrType string) (*Wallet, error) {
	filter := bson.M{}

	switch addrType {
	case ADDRESS_SEGWIT:
		filter["segwit_addr"] = addr
	case ADDRESS_TAPROOT:
		filter["tap_root_addr"] = addr
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

// func DetectAddressType(address string) string {
// 	if IsSegWitAddress(address) {
// 		return ADDRESS_SEGWIT
// 	} else if IsTaprootAddress(address) {
// 		return ADDRESS_TAPROOT
// 	} else {
// 		return ADDRESS_UNKNOWN
// 	}
// }

// func IsSegWitAddress(address string) bool {
// 	// SegWit addresses start with "bc1" for mainnet and "tb1" for testnet
// 	fmt.Println("IS SEG")
// 	return strings.HasPrefix(address, "bc1") || strings.HasPrefix(address, "tb1")
// }

// func IsTaprootAddress(address string) bool {
// 	fmt.Println("IS TAP")
// 	// Taproot addresses start with "bcrt1" for mainnet and "tbcrt1" for testnet
// 	return strings.HasPrefix(address, "bcrt1") || strings.HasPrefix(address, "tbcrt1")
// }

func GetAddressType(address string) (string, error) {
	decodedAddress, err := btcutil.DecodeAddress(address, nil)
	if err != nil {
		return "", err
	}

	switch decodedAddress.(type) {
	case *btcutil.AddressTaproot:
		return ADDRESS_TAPROOT, nil
	case *btcutil.AddressWitnessPubKeyHash:
		return ADDRESS_SEGWIT, nil
	default:
		return ADDRESS_UNKNOWN, nil
	}
}
