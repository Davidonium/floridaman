package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func SHA1String(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func HMAC256String(msg, key []byte) string {
	hm := hmac.New(sha256.New, key)
	hm.Write(msg)
	return hex.EncodeToString(hm.Sum(nil))
}
