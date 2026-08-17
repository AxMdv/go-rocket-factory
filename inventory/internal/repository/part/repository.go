package part

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(ctx context.Context, mongoURI, dbName, collectionName string) (*repository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return nil, err
	}
	db := client.Database(dbName)

	collection := db.Collection(collectionName)
	repo := &repository{collection: collection}

	empty, err := isCollectionEmpty(ctx, collection)
	if err != nil {
		return nil, err
	}
	log.Println("collection empty:", empty)
	if !empty {
		return repo, nil
	}

	parts := createRepoParts(20)
	if err = repo.insertParts(ctx, parts); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *repository) Close(ctx context.Context) error {
	client := r.collection.Database().Client()
	cerr := client.Disconnect(ctx)
	if cerr != nil {
		log.Printf("failed to disconnect: %v\n", cerr)
		return cerr
	}
	return nil
}

func isCollectionEmpty(ctx context.Context, collection *mongo.Collection) (bool, error) {
	err := collection.FindOne(
		ctx,
		bson.D{},
		options.FindOne().SetProjection(bson.D{
			{Key: "_id", Value: 1},
		})).Err()
	switch {
	case err == nil:
		return false, nil

	case errors.Is(err, mongo.ErrNoDocuments):
		return true, nil

	default:
		return false, fmt.Errorf(
			"check whether collection %q is empty: %w",
			collection.Name(),
			err,
		)
	}
}
