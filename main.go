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
		panic("Define OS Variables")
	}

	if err := connectToMongoDB(); err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/", fileServer)

	mux.HandleFunc("POST /buy/", middleware.SecretKeyMiddleware(handlers.BuyHandler))
	mux.HandleFunc("GET /balance/", middleware.SecretKeyMiddleware(handlers.BalanceHandler))

	fmt.Println("starting server on :4003")
	err := http.ListenAndServe("localhost:4003", mux)
	fmt.Println(err)

}

func connectToMongoDB() error {
	clientOptions := options.Client().ApplyURI(os.Getenv("okxMongoConnectionString"))

	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB")
	return nil
}
