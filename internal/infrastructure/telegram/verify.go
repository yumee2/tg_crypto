package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

func VerifyInitData(initData, botToken string) (map[string]string, bool) {
	// 1) Parse the raw query string
	vals, err := url.ParseQuery(initData)
	if err != nil {
		return nil, false
	}

	// 2) Collect key=value pairs, skipping "hash" and "signature"
	dataMap := make(map[string]string, len(vals))
	var parts []string
	var providedHash string

	for key, vs := range vals {
		if key == "hash" {
			providedHash = vs[0]
			continue
		}
		if key == "signature" {
			continue
		}
		dataMap[key] = vs[0]
		parts = append(parts, key+"="+vs[0])
	}

	// 3) Sort parts alphabetically and join with "\n"
	sort.Strings(parts)
	dataCheckString := strings.Join(parts, "\n")

	// 4) Derive secret_key = HMAC_SHA256(botToken, "WebAppData")
	h1 := hmac.New(sha256.New, []byte(botToken))
	h1.Write([]byte("WebAppData"))
	secretKey := h1.Sum(nil)

	// 5) Compute expectedHash = HMAC_SHA256(dataCheckString, secretKey)
	h2 := hmac.New(sha256.New, secretKey)
	h2.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h2.Sum(nil))

	// 6) Constant‑time compare
	if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(providedHash)) != 1 {
		return nil, false
	}

	return dataMap, true
}
