package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateContainerName(base string) string {
	if base == "" {
		base = "container"
	}

	randBytes := make([]byte, 2)
	if _, err := rand.Read(randBytes); err != nil {
		panic(fmt.Errorf("failed to generate random bytes: %v", err))
	}
	randomSuffix := hex.EncodeToString(randBytes)

	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s-%s-%s", base, timestamp, randomSuffix)
}
