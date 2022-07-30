package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func SHA1String(text string) string {
	h := sha1.New()
	h.Write([]byte(text))

	return hex.EncodeToString(h.Sum(nil))
}

func HMACString(msg, key []byte) string {
	hm := hmac.New(sha256.New, key)
	hm.Write(msg)

	finalHash := hm.Sum(nil)

	return fmt.Sprintf("v0=%s", hex.EncodeToString(finalHash))
}
