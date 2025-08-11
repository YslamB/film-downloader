package myhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateExpirableHash(duration time.Duration, secretKey, salt string) (string, int64) {

	expiresAt := time.Now().Add(duration).Unix()

	dataToHash := fmt.Sprintf("%s%s%d", secretKey, salt, expiresAt)

	hasher := sha256.New()
	hasher.Write([]byte(dataToHash))

	return hex.EncodeToString(hasher.Sum(nil)), expiresAt
}

func VerifyHash(inputHash, secretKey, salt string, expiresAt int64) bool {

	if time.Now().Unix() > expiresAt {
		fmt.Println("Error: Hash has expired.")
		return false
	}

	dataToHash := fmt.Sprintf("%s%s%d", secretKey, salt, expiresAt)

	hasher := sha256.New()
	hasher.Write([]byte(dataToHash))
	recreatedHash := hex.EncodeToString(hasher.Sum(nil))

	if recreatedHash != inputHash {
		fmt.Println("Error: Hash is invalid.")
		return false
	}

	return true
}
