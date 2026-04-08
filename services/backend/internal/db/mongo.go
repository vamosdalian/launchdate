package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client       *mongo.Client
	Database     *mongo.Database
	DatabaseName string
}

func NewMongoDB(uri, database string) (*MongoDB, func(), error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	return &MongoDB{
			Client:       client,
			Database:     client.Database(database),
			DatabaseName: database,
		}, func() {
			client.Disconnect(context.Background())
		}, nil
}

func (db *MongoDB) Collection(name string) *mongo.Collection {
	return db.Database.Collection(name)
}
