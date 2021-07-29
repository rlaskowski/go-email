package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var hashNum = []byte{
	2, 36, 73, 28,
	33, 29, 27, 59,
	21, 69, 18, 35,
	6, 0, 2, 62,
	17, 57, 32, 18,
	0, 1, 33, 37,
	51, 59, 34, 3,
	35, 26, 55, 70,
	28, 37, 56, 19,
	41, 20, 72, 68,
}

func computeKey() []byte {
	key := make([]byte, 0)

	for _, v := range hashNum {
		key = append(key, v+((2<<1)^(3|4))<<2<<2)
	}

	return key
}

func ComputeHash(val string) (string, error) {
	key := computeKey()

	hash := hmac.New(sha256.New, key)

	if _, err := fmt.Fprintf(hash, "%s", val); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
