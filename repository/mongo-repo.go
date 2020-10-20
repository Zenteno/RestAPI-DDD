package repository

import (
	"context"
	"api-ddd/entity"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

const (
	reportedAt       = "reported_at"
	sessionUUID      = "session_uuid"
	shopperUUID      = "shopper_uuid"
	orderDesc        = -1
	greaterEqualThan = "$gte"
)

func newMongoClient(mongoHost string, mongoPort int, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", mongoHost, mongoPort)))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepository(mongoHost string, mongoPort int, mongoDB string, mongoTimeout int) SessionRepository {
	client, _ := newMongoClient(mongoHost, mongoPort, mongoTimeout)
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
		client:   client,
	}
	return repo
}

func (r *mongoRepository) GetSessionHistory(sessionUUIDVal string) ([]*entity.ShopperHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("sessions")
	var results []*entity.ShopperHistory
	findOpts := options.Find()
	findOpts.SetSort(bson.D{primitive.E{Key: reportedAt, Value: orderDesc}})
	cur, err := collection.Find(ctx, bson.D{primitive.E{Key: sessionUUID, Value: sessionUUIDVal}}, findOpts)
	if err != nil {
		return results, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		s := &entity.ShopperHistory{}
		err := cur.Decode(&s)
		if err != nil {
			return results, err
		}
		results = append(results, s)
	}
	return results, nil
}

func (r *mongoRepository) GetShopperLocation(shoperUUIDVal string) (*entity.ShopperHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("sessions")
	now := time.Now()
	then := now.Add(-10 * time.Minute)
	findOpts := options.FindOne()
	findOpts.SetSort(bson.D{primitive.E{Key: reportedAt, Value: orderDesc}})
	result := &entity.ShopperHistory{}
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: shopperUUID, Value: shoperUUIDVal},
		primitive.E{Key: reportedAt, Value: bson.D{{Key: greaterEqualThan, Value: then.UTC()}}},
	}, findOpts).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *mongoRepository) Insert(session *entity.ShopperHistory) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("sessions")
	_, err := collection.InsertOne(
		ctx,
		session,
	)
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
