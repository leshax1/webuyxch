package database

import (
	"context"
	"fmt"
	"webuyxch/models"
	"webuyxch/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertTradeData(client *mongo.Client, tradeData models.TradeData) (models.TradeData, error) {
	collection := client.Database("webuyxch").Collection("trades")

	strings := []string{tradeData.Sz, tradeData.Px}
	floats, stringConvertionError := utils.ConvertStringsToFloat32(strings)

	if stringConvertionError != nil {
		return tradeData, fmt.Errorf("error on string to float converstion  %w", stringConvertionError)
	}

	tradeTime, timeConvertionError := utils.ConvertMilliStringToTime(tradeData.Ts)
	if timeConvertionError != nil {
		return tradeData, fmt.Errorf("Failed to convert timestamp  %w", timeConvertionError)
	}

	trade := bson.D{
		{Key: "Sz", Value: floats[0]},
		{Key: "Px", Value: floats[1]},
		{Key: "Ts", Value: tradeTime},
	}

	_, err := collection.InsertOne(context.TODO(), trade)
	if err != nil {
		return tradeData, fmt.Errorf("unable to insert document %w", err)
	}

	return tradeData, nil
}
