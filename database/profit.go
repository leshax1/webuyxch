package database

import (
	"context"
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

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: nil},
		{Key: "averagePx", Value: bson.D{{Key: "$avg", Value: "$Px"}}},
		{Key: "totalSz", Value: bson.D{{Key: "$sum", Value: "$Sz"}}},
	}}}
	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{groupStage})
	if err != nil {
		return models.Pnl{}, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		return models.Pnl{}, err
	}
	if len(results) == 0 {
		return models.Pnl{}, fmt.Errorf("no transactions found")
	}
	averageBuyPx := results[0]["averagePx"].(float64)
	totalBuySz := results[0]["totalSz"].(float64)

	var lastSellTransaction bson.M
	findOptions := options.FindOne().SetSort(bson.D{{Key: "Ts", Value: -1}})
	if err := collection.FindOne(context.TODO(), bson.D{}, findOptions).Decode(&lastSellTransaction); err != nil {
		return models.Pnl{}, err
	}
	lastSellPx := lastSellTransaction["Px"].(float64)

	fmt.Println(averageBuyPx)

	pnlDollars := (lastSellPx - averageBuyPx) * totalBuySz
	pnlPercent := (lastSellPx - averageBuyPx) / averageBuyPx * 100

	return models.Pnl{ProfitUsd: pnlDollars, Profit: pnlPercent}, nil
}
