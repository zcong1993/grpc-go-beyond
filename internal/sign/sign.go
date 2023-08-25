package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func B64(bt []byte) string {
	return base64.StdEncoding.EncodeToString(bt)
}

func Sign(key string, payload []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(payload)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
