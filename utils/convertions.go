package utils

import (
	"fmt"
	"strconv"
	"time"
)

func ConvertStringsToFloat32(strings []string) ([]float32, error) {
	// Create a slice to hold the converted float32 values
	floatValues := make([]float32, len(strings))

	// Iterate over the input slice of strings
	for i, s := range strings {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			// Return the error and where it occurred
			return nil, fmt.Errorf("error converting string '%s' to float32 at index %d: %v", s, i, err)
		}
		// Store the converted float32 value
		floatValues[i] = float32(f)
	}

	// Return the slice of float32 values and no error
	return floatValues, nil
}

func ConvertMilliStringToTime(millis string) (time.Time, error) {
	// Convert string to int64
	ms, err := strconv.ParseInt(millis, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("error converting string to int64: %v", err)
	}

	// Convert milliseconds to seconds
	seconds := ms / 1000
	// Convert the remainder to nanoseconds
	nanoseconds := (ms % 1000) * 1000000

	// Return the time corresponding to the Unix seconds and nanoseconds
	return time.Unix(seconds, nanoseconds), nil
}
