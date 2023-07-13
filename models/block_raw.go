package models

import (
	"context"
	"github.com/btcsuite/btcd/btcjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ntc-services/stores"
	"time"
)

type BlockRaw struct {
	ID        primitive.ObjectID             `json:"id" bson:"_id"`
	Hash      string                         `json:"hash" bson:"hash"`
	Height    int64                          `json:"height" bson:"height"`
	Completed bool                           `json:"completed" bson:"completed"`
	CreatedAt time.Time                      `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time                     `json:"updated_at" bson:"updated_at"`
	Block     *btcjson.GetBlockVerboseResult `json:"block" bson:"block"`
}

func NewBlockRaw(block *btcjson.GetBlockVerboseResult) *BlockRaw {
	return &BlockRaw{
		ID:        primitive.NewObjectID(),
		Hash:      block.Hash,
		Height:    block.Height,
		Completed: false,
		CreatedAt: time.Now().UTC(),
		Block:     block,
	}
}

func (br *BlockRaw) Complete() error {
	br.Completed = true

	return nil
}

func (br *BlockRaw) Save() error {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_BLOCKS_RAW)
	if _, err := collection.InsertOne(context.TODO(), br, nil); err != nil {
		return err
	}

	return nil
}

func (br *BlockRaw) Update() error {
	now := time.Now().UTC()
	br.UpdatedAt = &now
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_BLOCKS_RAW)
	if _, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": br.ID}, br); err != nil {
		return err
	}

	return nil
}

func GetBlockRawByID(idStr string) (*BlockRaw, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": id}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_BLOCKS_RAW)

	var br *BlockRaw
	if err := collection.FindOne(context.TODO(), filter).Decode(&br); err != nil {
		return nil, err
	}

	return br, nil
}

func GetBlockRawByHeight(height int64) (*BlockRaw, error) {
	filter := bson.M{"height": height}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_BLOCKS_RAW)

	var br *BlockRaw
	if err := collection.FindOne(context.TODO(), filter).Decode(&br); err != nil {
		return nil, err
	}

	return br, nil
}

func GetBlockRawCompleted(idStr string) (bool, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": id}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_BLOCKS_RAW)

	var br *BlockRaw
	if err := collection.FindOne(context.TODO(), filter).Decode(&br); err != nil {
		return false, err
	}

	return br.Completed, nil
}

func GetBlockStarted(height int64) (bool, error) {
	filter := bson.M{"height": height}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_BLOCKS_RAW)

	var br *BlockRaw
	if err := collection.FindOne(context.TODO(), filter).Decode(&br); err != nil {
		if err.Error() == stores.MONGO_ERR_NOT_FOUND {
			return false, nil
		}
	}

	return true, nil
}

func GetLatestBlock() (*BlockRaw, error) {
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_BLOCKS_RAW)
	opts := options.FindOne().SetSort(bson.D{{"height", -1}})

	var br *BlockRaw
	if err := collection.FindOne(context.Background(), bson.D{{}}, opts).Decode(&br); err != nil {
		return nil, err
	}

	return br, nil
}

func GetUncompletedBlockRaws() ([]*BlockRaw, error) {
	filter := bson.M{"completed": false}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_BLOCKS_RAW)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var blockRaws []*BlockRaw
	if err := cursor.All(context.TODO(), &blockRaws); err != nil {
		return nil, err
	}

	return blockRaws, nil
}

// TODO: Save but check to pre existing record
