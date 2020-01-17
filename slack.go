package floridaman

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type SlackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func ValidateSlackRequest(r *http.Request, logger *log.Logger) bool {
	ssig := r.Header.Get("X-Slack-Signature")
	t := r.Header.Get("X-Slack-Request-Timestamp")

	ts, err := strconv.ParseInt(t, 10, 64)

	if err != nil {
		return false
	}

	tsu := time.Unix(ts, 0)

	if time.Now().Sub(tsu) > 5*time.Minute {
		logger.Println("timestamp difference is greater than 5 minutes")
		return false
	}

	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		return false
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	msg := fmt.Sprintf("v0:%s:%s", t, body)

	ss := os.Getenv("SLACK_SIGNING_SECRET")
	sig := hashHMAC([]byte(msg), []byte(ss))
	ok := hmac.Equal([]byte(sig), []byte(ssig))
	if !ok {
		logger.Printf("error validating hmac signature from slack: %s, generated %s\n", ssig, sig)
		return false
	}

	return true
}

func hashHMAC(msg, key []byte) string {
	hm := hmac.New(sha256.New, key)
	hm.Write(msg)
	finalHash := hm.Sum(nil)
	return fmt.Sprintf("v0=%s", hex.EncodeToString(finalHash))
}
