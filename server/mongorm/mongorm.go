package mongorm

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *Model) Create(ctx context.Context, db *mongo.Database, collectionName string, model interface{}) error {
	coll := db.Collection(collectionName)

	m.CreateAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := coll.InsertOne(ctx, model)

	if err != nil {
		return err
	}

	m.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (m *Model) Read(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, model interface{}) error {
	coll := db.Collection(collectionName)

	err := coll.FindOne(ctx, filter).Decode(model)

	if err != nil {
		return err
	}

	return nil
}

func (m *Model) Update(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, update interface{}) error {
	coll := db.Collection(collectionName)

	m.UpdatedAt = time.Now()

	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {

		return err
	}

	return nil
}

func (m *Model) Delete(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}) error {
	coll := db.Collection(collectionName)
	_, err := coll.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	return nil
}
