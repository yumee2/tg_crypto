package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func VerifyInitData(initData, botToken string) (map[string]string, bool) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, false
	}

	dataMap := map[string]string{}
	var pairs []string
	var providedHash string

	for k, v := range values {
		if k == "hash" {
			providedHash = v[0]
			continue
		}
		dataMap[k] = v[0]
		pairs = append(pairs, k+"="+v[0])
	}

	sort.Strings(pairs)
	dataCheckString := strings.Join(pairs, "\n")

	secretKey := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Provided hash:", providedHash)
	fmt.Println("Expected hash:", expectedHash)
	if expectedHash != providedHash {
		return nil, false
	}
	return dataMap, true
}
