package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func VerifyInitData(initData, botToken string) (map[string]string, bool) {
	vals, err := url.ParseQuery(initData)
	if err != nil {
		return nil, false
	}

	// 1) collect everything except hash & signature
	dataMap := make(map[string]string, len(vals))
	var parts []string
	var providedHash string

	for k, vs := range vals {
		if k == "hash" || k == "signature" {
			if k == "hash" {
				providedHash = vs[0]
			}
			continue
		}
		dataMap[k] = vs[0]
		parts = append(parts, k+"="+vs[0])
	}

	// 2) sort and join with "\n"
	sort.Strings(parts)
	dataCheckString := strings.Join(parts, "\n")
	fmt.Printf(">> DATA CHECK STRING:\n%s\n<< END\n", dataCheckString)

	// 3) derive secret_key = HMAC(botToken, "WebAppData")
	h1 := hmac.New(sha256.New, []byte(botToken))
	h1.Write([]byte("WebAppData"))
	secretKey := h1.Sum(nil)

	// 4) compute expected hash = HMAC(dataCheckString, secretKey)
	h2 := hmac.New(sha256.New, secretKey)
	h2.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h2.Sum(nil))

	fmt.Println("Provided hash:", providedHash)
	fmt.Println("Expected hash:", expectedHash)

	// 5) constant‑time compare
	if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(providedHash)) != 1 {
		return nil, false
	}
	return dataMap, true
}
