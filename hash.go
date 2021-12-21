package floridaman

import (
	"crypto/sha1"
	"encoding/hex"
)

func SHA1String(text string) string {
	h := sha1.New()
	h.Write([]byte(text))

	return hex.EncodeToString(h.Sum(nil))
}
