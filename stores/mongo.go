package stores

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ntc-services/config"
	"time"
)

type Mongo struct {
	Hosts    []string `json:"hosts"`
	Client   *mongo.Client
	Username string
	Password string
}

func NewMongo(hosts []string) *Mongo {
	return &Mongo{
		Hosts: hosts,
	}
}

func (m *Mongo) Connect() error {
	fmt.Println("=============================")
	fmt.Println(fmt.Sprintf("%+v", m.Hosts))
	fmt.Println("=============================")

	t := 120 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), t)

	opts := &options.ClientOptions{
		ConnectTimeout: &t,
		Hosts:          m.Hosts,
	}
	if config.Conf.GetEnv() == config.STAGE || config.Conf.GetEnv() == config.PROD {
		user, err := config.GetMongoUser()
		if err != nil {
			return err
		}
		password, err := config.GetMongoPassword()
		if err != nil {
			return err
		}
		m.Username = *user
		m.Password = *password
		opts.Auth = &options.Credential{
			Username: m.Username,
			Password: m.Password,
		}
	}

	// IMPORTANT: the following block of code is to audit the production database and should
	// never be committed uncommented
	//if config.Conf.GetEnv() == config.WORKSTATION {
	//	user, err := config.GetMongoUser()
	//	if err != nil {
	//		return err
	//	}
	//	password, err := config.GetMongoPassword()
	//	if err != nil {
	//		return err
	//	}
	//	m.Username = *user
	//	m.Password = *password
	//	opts.Auth = &options.Credential{
	//		Username: m.Username,
	//		Password: m.Password,
	//	}
	//}

	if config.Conf.GetEnv() == config.WORKSTATION {
		direct := true
		opts.Direct = &direct
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return err
	}

	client.Connect(ctx)
	m.Client = client
	if err = m.Client.Ping(ctx, nil); err != nil {
		return err
	}

	if err := m.CreateIndexes(); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) CreateIndexes() error {
	if _, err := m.Client.Database(DB_NAME).Collection(DB_COLLECTION_BLOCKS_RAW).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "hash", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		return err
	}

	if _, err := m.Client.Database(DB_NAME).Collection(DB_COLLECTION_TXS_RAW).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "hash", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		return err
	}

	return nil
}
