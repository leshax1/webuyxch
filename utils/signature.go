package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func CalculateSignature(apiSecret string, timestamp string, method string, requestPath string, body string) string {
	prehashString := fmt.Sprintf("%s%s%s%s", timestamp, method, requestPath, body)

	// Prepare the SecretKey
	key := []byte(apiSecret)

	// Sign the prehash string with HMAC SHA256
	h := hmac.New(sha256.New, key)
	h.Write([]byte(prehashString))
	signature := h.Sum(nil)

	// Encode the signature in Base64 format
	return base64.StdEncoding.EncodeToString(signature)
}
