package main

import (
	"fmt"
	"math/big"
	"crypto/rand"
)


func generateCode() (string, error) {
	const charSet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		result[i] = charSet[n.Int64()]
	}
	return string(result), nil
}

func isUnique(code string) bool {	
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	_, ok := UrlMappings[code]
	if !ok {
		return true
	}
	return false
}

func generateUniqueCode() string {
	for {
		code, err := generateCode()
		if err != nil {
			fmt.Println("Something went wrong")
		}
		if isUnique(code) {
			return code
		}
	}
	
}
