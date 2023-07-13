package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/stores"
	"time"
)

type State struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	//Path               string             `json:"path" bson:"path"`
	ScannerBlockHeight int64      `json:"scanner_block_height" bson:"scanner_block_height"`
	ParserBlockHeight  int64      `json:"parser_block_height" bson:"parser_block_height"`
	UpdatedAt          *time.Time `json:"updated_at" bson:"updated_at"`
}

func NewState() *State {
	return &State{
		//Path:               path,
		ScannerBlockHeight: -1,
		ParserBlockHeight:  -1,
	}
}

func (s *State) Update() error {
	now := time.Now().UTC()
	s.UpdatedAt = &now

	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_STATE)
	if _, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": s.ID}, s); err != nil {
		return err
	}

	return nil
}

func GetState() (*State, error) {
	filter := bson.M{}
	collection := stores.DB.Mongo.Client.Database(stores.DB_NAME).Collection(stores.DB_COLLECTION_STATE)

	var s *State
	if err := collection.FindOne(context.TODO(), filter).Decode(&s); err != nil {
		return nil, err
	}

	return s, nil
}

//func (s *State) Read() error {
//	data, err := ioutil.ReadFile(s.Path)
//	if err != nil {
//		return fmt.Errorf("failed to read state: %w", err)
//	}
//
//	if err := json.Unmarshal(data, &s); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (s *State) Write() error {
//	data, err := json.Marshal(s)
//	if err != nil {
//		return err
//	}
//	if err := ioutil.WriteFile(s.Path, data, 0644); err != nil {
//		return fmt.Errorf("failed to write state: %w", err)
//	}
//
//	return nil
//}
