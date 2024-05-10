package database

import (
	"context"
	"errors"
	"log"
	"webuyxch/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Profit(client *mongo.Client) (int16, error) {
	collection := client.Database("webuyxch").Collection("trades")

	opts := options.FindOne().SetSort(bson.D{{Key: "Ts", Value: -1}})
	var lastTrade models.Trade
	err := collection.FindOne(context.TODO(), bson.D{{}}, opts).Decode(&lastTrade)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err // Return an error if it occurs
	}

	pipeline := mongo.Pipeline{
		bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: nil}, // Grouping key '_id' set to nil for aggregating all documents
				{Key: "averagePx", Value: bson.D{{Key: "$avg", Value: "$Px"}}}, // Calculating average
			},
		}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	if len(results) > 0 {
		if avg, ok := results[0]["averagePx"].(float64); ok {
			averagePx := float64(avg)
			return int16((lastTrade.Px/averagePx)*100 - 100), nil
		} else {
			return 0, errors.New("can not convert AVG to float32: ")
		}
	} else {
		return 0, nil
	}
}
