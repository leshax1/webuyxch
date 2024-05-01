package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"webuyxch/database"
	"webuyxch/models"
	"webuyxch/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

// Response represents the top-level structure of the API response
type Response struct {
	Code string             `json:"code"`
	Msg  string             `json:"msg"`
	Data []models.TradeData `json:"data"`
}

func updateLastTrade(DB *mongo.Client) ([]byte, error) {
	method := "GET"

	endpoint := "/api/v5/market/history-trades?instId=XCH-USDT&limit=1"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")

	url := fmt.Sprintf("%s%s", os.Getenv("okxBaseUrl"), endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get the list of trades")
	}

	req.Header.Set("OK-ACCESS-KEY", os.Getenv("okxApiKey"))
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", os.Getenv("okxPassPhrase"))
	req.Header.Set("Content-Type", "application/json")

	if os.Getenv("okxSimulatedTrading") == "1" {
		fmt.Println("Simulation header added")
		req.Header.Set("x-simulated-trading", "1")
	}

	signature := utils.CalculateSignature(os.Getenv("okxApiSecret"), timestamp, method, endpoint, "")

	req.Header.Set("OK-ACCESS-SIGN", signature)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %v", err)
	}

	fmt.Println("Server Response:", string(body))

	var apiResponse Response
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	if len(apiResponse.Data) == 0 {
		return nil, fmt.Errorf("data array is empty: %v", apiResponse.Msg)
	}

	tradeData, err := database.InsertTradeData(DB, apiResponse.Data[0])
	if err != nil {
		return nil, err
	}

	fmt.Printf("Inserted tradeData with sz: %v , price: %v, time: %s \n", tradeData.Sz, tradeData.Px, tradeData.Ts)

	return body, nil
}
