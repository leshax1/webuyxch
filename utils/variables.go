package utils

import (
	"fmt"
	"os"
)

func VariablesCheck() bool {
	variables := []string{"okxApiKey", "okxPassPhrase", "okxApiSecret"}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			fmt.Printf("%s is not declared\n", variable)
			return false
		}
	}
	return true
}
