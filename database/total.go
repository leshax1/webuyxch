package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Total(client *mongo.Client) (float64, error) {
	collection := client.Database("webuyxch").Collection("trades")

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "totalSz", Value: bson.D{{Key: "$sum", Value: "$Sz"}}},
			}},
		},
	}

	// Execute the aggregation query
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())

	// Parse the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return 0, err
	}

	// Check the result and extract the sum
	if len(results) > 0 {
		if totalSz, ok := results[0]["totalSz"].(float64); ok {
			return totalSz, nil
		}
		return 0, fmt.Errorf("could not find 'totalSz' in the results")
	}

	return 0, fmt.Errorf("no documents found in collection")
}
