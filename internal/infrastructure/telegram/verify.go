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
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, false
	}

	dataMap := map[string]string{}
	var pairs []string
	var providedHash string

	for k, v := range values {
		switch k {
		case "hash":
			providedHash = v[0]
			continue
		case "signature":
			continue
		}

		dataMap[k] = v[0]
		pairs = append(pairs, k+"="+v[0])
	}

	sort.Strings(pairs)
	dataCheckString := strings.Join(pairs, "\n")

	// derive the secret key from your bot token
	secretKey := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	fmt.Println("Provided hash:", providedHash)
	fmt.Println("Expected hash:", expectedHash)

	if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(providedHash)) != 1 {
		return nil, false
	}
	return dataMap, true
}
