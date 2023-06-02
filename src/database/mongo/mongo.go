package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Close(
	client *mongo.Client,
	ctx context.Context,
	cancel context.CancelFunc,
) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	return client, ctx, cancel, err
}

func Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Connected succesfully.")
	return nil
}

func insertOne(client *mongo.Client, ctx context.Context, database, col string, doc interface{}) (*mongo.InsertOneResult, error) {
	collection := client.Database(database).Collection(col)

	result, err := collection.InsertOne(ctx, doc)
	return result, err
}
