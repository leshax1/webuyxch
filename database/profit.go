package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"webuyxch/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Profit(client *mongo.Client) (models.Pnl, error) {
	pnl, error := GetPln(client)
	if error != nil {
		log.Printf("Error occurred: %v\n", error)
		return pnl, error
	}

	return pnl, nil

}

func GetPln(client *mongo.Client) (models.Pnl, error) {
	collection := client.Database("webuyxch").Collection("trades")

	pnl := new(models.Pnl)

	opts := options.FindOne().SetSort(bson.D{{Key: "Ts", Value: -1}})
	var lastTrade models.Trade
	err := collection.FindOne(context.TODO(), bson.D{{}}, opts).Decode(&lastTrade)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return *pnl, nil
		}
		return *pnl, err // Return an error if it occurs
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
		log.Printf("Error occurred: %v\n", err)
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Printf("Error occurred: %v\n", err)
	}

	////
	pipelineSum := mongo.Pipeline{
		bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: nil},                                   // Grouping key '_id' set to nil for aggregating all documents
				{Key: "sumPx", Value: bson.D{{Key: "$sum", Value: "$Px"}}}, // Calculating sum
			},
		}},
	}

	cursorSumPx, err := collection.Aggregate(context.TODO(), pipelineSum)
	if err != nil {
		log.Printf("Error occurred: %v\n", err)
	}
	defer cursorSumPx.Close(context.TODO())

	var resultsSumPx []bson.M
	if err = cursorSumPx.All(context.TODO(), &resultsSumPx); err != nil {
		log.Printf("Error occurred: %v\n", err)
	}

	///
	var averagePx float64
	if len(results) > 0 {
		if avg, ok := results[0]["averagePx"].(float64); ok {
			averagePx = float64(avg)
			pnl.Profit = (lastTrade.Px/averagePx)*100 - 100
		} else {
			return *pnl, errors.New("can not convert AVG to float64")
		}
	} else {
		return *pnl, nil
	}

	fmt.Println(resultsSumPx)
	if len(resultsSumPx) > 0 {
		if sum, ok := resultsSumPx[0]["sumPx"].(float64); ok {
			sumPx := float64(sum)
			fmt.Printf("(%f - %f) * %f", averagePx, lastTrade.Px, sumPx)
			pnl.ProfitUsd = (pnl.Profit / 100) * sumPx
		} else {
			return *pnl, errors.New("can not convert SumPx to float64")
		}
	} else {
		return *pnl, nil
	}

	return *pnl, nil
}
