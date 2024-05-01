package database

import (
	"context"
	"fmt"

	"webuyxch/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertTradeData(client *mongo.Client, tradeData models.TradeData) (models.TradeData, error) {
	collection := client.Database("webuyxch").Collection("trades")

	trade := bson.D{
		{Key: "Sz", Value: tradeData.Sz},
		{Key: "Px", Value: tradeData.Px},
		{Key: "Ts", Value: tradeData.Ts},
	}

	_, err := collection.InsertOne(context.TODO(), trade)
	if err != nil {
		return tradeData, fmt.Errorf("Unable to insert document %v", err)
	}

	return tradeData, nil
}
