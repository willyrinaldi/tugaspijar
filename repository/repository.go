package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName     = "simple_db"
	simpleColl = "simple"
)

type SimpleEntity struct {
	ID      primitive.ObjectID `bson:"_id"`
	User    string             `bson:"user"`
	Name    string             `bson:"name"`
	Created time.Time          `bson:"created"`
	Updated time.Time          `bson:"updated"`
}

func (e *SimpleEntity) Timestamp() {
	t := time.Now()
	e.Updated = t
	if e.Created.IsZero() {
		e.Created = t
	}
}

type Repository interface {
	Simple() ([]SimpleEntity, error)
	Insert(user, name string) (SimpleEntity, error)
	Update(id string, ent SimpleEntity) (SimpleEntity, error)
	SimpleOne(id string) (SimpleEntity, error)
}

type repository struct{}

func NewRepository() Repository {
	return repository{}
}

func (r repository) Simple() ([]SimpleEntity, error) {
	var result []SimpleEntity

	client, err := mongoClient()
	defer client.Disconnect(context.Background())
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := client.
		Database(dbName).
		Collection(simpleColl).
		Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	err = cursor.All(ctx, &result)

	return result, err
}

func (r repository) Insert(user, name string) (SimpleEntity, error) {
	var result SimpleEntity
	client, err := mongoClient()
	defer client.Disconnect(context.Background())
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result = SimpleEntity{
		ID:   primitive.NewObjectID(),
		User: user,
		Name: name,
	}
	result.Timestamp()

	_, err = client.
		Database(dbName).
		Collection(simpleColl).
		InsertOne(ctx, result)
	return result, err
}

func (r repository) Update(id string, ent SimpleEntity) (SimpleEntity, error) {
	var result SimpleEntity
	client, err := mongoClient()
	defer client.Disconnect(context.Background())
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	err = client.
		Database(dbName).
		Collection(simpleColl).
		FindOne(ctx, bson.M{"_id": ID}).
		Decode(&result)
	if err != nil {
		fmt.Println("FindOne")
		fmt.Println(err)
		return result, err
	}

	result.User = ent.User
	result.Name = ent.Name
	result.Timestamp()

	_, err = client.
		Database(dbName).
		Collection(simpleColl).
		UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": result})
	if err != nil {
		fmt.Println("UpdateOne")
		fmt.Println(err)
	}

	return result, err
}

func (r repository) SimpleOne(id string) (SimpleEntity, error) {
	var result SimpleEntity

	client, err := mongoClient()
	defer client.Disconnect(context.Background())
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.
		Database(dbName).
		Collection(simpleColl).
		FindOne(ctx, bson.M{"_id": ID}).
		Decode(&result)

	return result, err
}

func mongoClient() (*mongo.Client, error) {
	var err error
	var client *mongo.Client
	uri := "mongodb://localhost:27017"
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetMaxPoolSize(10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if client, err = mongo.Connect(ctx, opts); err != nil {
		fmt.Println("connect...")
		fmt.Println(err.Error())
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("ping...")
		fmt.Println(err)
		return nil, err
	}

	return client, nil
}