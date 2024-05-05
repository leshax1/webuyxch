package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"net/http"
	"webuyxch/handlers"
	"webuyxch/middleware"
	"webuyxch/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	if !utils.VariablesCheck() {
		panic("Define all required OS Variables")
	}

	client, err := connectToMongoDB()

	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		panic("Update mongo connection config")
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/", fileServer)
	mux.Handle("POST /buy/", middleware.SecretKeyMiddleware(&handlers.BuyHandler{DB: client}))
	mux.Handle("GET /balance/", middleware.SecretKeyMiddleware(&handlers.BalanceHandler{DB: client}))

	fmt.Println("Starting server on :3001...")
	err = http.ListenAndServe(":3001", mux)
	fmt.Println(err)

}

func connectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("okxMongoConnectionString"))

	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return client, nil
}
