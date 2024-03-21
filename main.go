package main

import (
	"fmt"
	"net/http"
	"webuyxch/handlers"
	"webuyxch/utils"
)

func main() {
	if !utils.VariablesCheck() {
		panic("Define OS Variables")
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/", fileServer)

	mux.HandleFunc("GET /buy/", handlers.BuyHandler)
	mux.HandleFunc("GET /balance/", handlers.BalanceHandler)

	fmt.Println("starting server on :4000")
	err := http.ListenAndServe("localhost:4000", mux)
	fmt.Println(err)

}
