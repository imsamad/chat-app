package mongorm

import (
	"context"
	"fmt"
	"reflect"
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

func CountDocuments(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}) (int64, error) {
	coll := db.Collection(collectionName)

	count, err := coll.CountDocuments(ctx, filter)

	if err != nil {
		return 0, nil
	}

	return count, nil
}

func (m *Model) Create(ctx context.Context, db *mongo.Database, collectionName string, model interface{}) error {
	coll := db.Collection(collectionName)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := coll.InsertOne(ctx, model)

	if err != nil {
		return err
	}

	m.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (m *Model) Read(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, model interface{}, sort interface{}) error {
	coll := db.Collection(collectionName)

	// Prepare the options with sorting
	findOptions := options.FindOne()
	if sort != nil {
		findOptions.SetSort(sort)
	}

	err := coll.FindOne(ctx, filter, findOptions).Decode(model)

	return err
}

func ReadAll(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, result interface{}, sort interface{}) error {
	coll := db.Collection(collectionName)

	// Prepare the options with sorting
	findOptions := options.Find()
	if sort != nil {
		findOptions.SetSort(sort)
	}

	cursor, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		return err
	}

	// Ensure `result` is a pointer to a slice
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr || resultValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("result must be a pointer to a slice")
	}
	sliceValue := resultValue.Elem()

	// Iterate over the cursor and decode documents
	for cursor.Next(ctx) {
		// Create a new instance of the element type of the slice
		elemType := sliceValue.Type().Elem()
		elem := reflect.New(elemType).Elem()

		// Decode into the new instance
		if err := cursor.Decode(elem.Addr().Interface()); err != nil {
			return err
		}

		// Append the new instance to the slice
		sliceValue.Set(reflect.Append(sliceValue, elem))
	}

	// Check for errors encountered during iteration
	if err = cursor.Err(); err != nil {
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

func Delete(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}) error {
	coll := db.Collection(collectionName)
	_, err := coll.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	return nil
}
