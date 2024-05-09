package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"webuyxch/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type RequestBody struct {
	InstID  string `json:"instId"`
	TDMode  string `json:"tdMode"`
	Side    string `json:"side"`
	OrdType string `json:"ordType"`
	Sz      string `json:"sz"`
	TgtCcy  string `json:"tgtCcy"`
}

type BuyRequest struct {
	Quantity float32 `json:"quantity,string"`
}

type BuyHandler struct {
	DB *mongo.Client
}

func (buy *BuyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint := "/api/v5/trade/order"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	url := fmt.Sprintf("%s%s", os.Getenv("okxBaseUrl"), endpoint)

	var buyRequest BuyRequest
	err := json.NewDecoder(r.Body).Decode(&buyRequest)

	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	log.Println(buyRequest.Quantity)

	if !(buyRequest.Quantity >= 0.1 && buyRequest.Quantity < 100) {
		http.Error(w, "Quantity of xch is incorrent, should be between 0.1 and 100", http.StatusBadRequest)
		return
	}

	log.Println("Buying... " + fmt.Sprintf("%f", buyRequest.Quantity) + " xch")

	requestData := RequestBody{
		InstID:  "XCH-USDT",
		TDMode:  "cash",
		Side:    "buy",
		OrdType: "market",
		TgtCcy:  "base_ccy",
		Sz:      fmt.Sprintf("%f", buyRequest.Quantity),
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		http.Error(w, "Unable to parse ", http.StatusBadRequest)
		log.Println("Error marshaling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("OK-ACCESS-KEY", os.Getenv("okxApiKey"))
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", os.Getenv("okxPassPhrase"))
	req.Header.Set("Content-Type", "application/json")

	if os.Getenv("okxSimulatedTrading") == "1" {
		log.Println("Simulation header added")
		req.Header.Set("x-simulated-trading", "1")
	}

	signature := utils.CalculateSignature(os.Getenv("okxApiSecret"), timestamp, "POST", endpoint, string(requestBody))
	req.Header.Set("OK-ACCESS-SIGN", signature)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	log.Println(string(body))

	time.Sleep(5 * time.Second)

	orderDetails, err := updateLastTrade(buy.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(orderDetails)

}
