package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"webuyxch/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type BalanceHandler struct {
	DB *mongo.Client
}

func (h *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint := "/api/v5/account/balance"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	url := fmt.Sprintf("%s%s", os.Getenv("okxBaseUrl"), endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("OK-ACCESS-KEY", os.Getenv("okxApiKey"))
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", os.Getenv("okxPassPhrase"))
	if os.Getenv("okxSimulatedTrading") == "1" {
		req.Header.Set("x-simulated-trading", "1")
	}

	signature := utils.CalculateSignature(os.Getenv("okxApiSecret"), timestamp, "GET", endpoint, "")
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
