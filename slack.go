package floridaman

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func NewSlackRandomHandler(logger *log.Logger, ar ArticleReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ok := ValidateSlackRequest(r, logger, os.Getenv("SLACK_SIGNING_SECRET"))

		if !ok {
			logger.Printf("invalid slack request %v\n", r)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid slack request"})
			return
		}

		article, err := ar.Random()

		if err != nil {
			WriteInternalServerError(w)
			logger.Printf("%v\n", err)
			return
		}

		response := SlackResponse{
			Text:         fmt.Sprintf("%s (%s)", article.Title, article.Link),
			ResponseType: "in_channel",
		}
		json.NewEncoder(w).Encode(response)
	}
}

func NewSlackOAuthRedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("floridaman oauth success!"))
		w.WriteHeader(200)
	}
}

func ValidateSlackRequest(r *http.Request, logger *log.Logger, ssecret string) bool {
	ssig := r.Header.Get("X-Slack-Signature")
	t := r.Header.Get("X-Slack-Request-Timestamp")

	ts, err := strconv.ParseInt(t, 10, 64)

	if err != nil {
		logger.Printf("X-Slack-Request-Timestamp is not a parsable number, got \"%s\"\n", t)
		return false
	}

	tsu := time.Unix(ts, 0)

	if time.Now().Sub(tsu) > 5*time.Minute {
		logger.Println("timestamp difference is greater than 5 minutes")
		return false
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Printf("error reading body %v", err)
		return false
	}

	defer r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	msg := fmt.Sprintf("v0:%s:%s", t, body)

	sig := hashHMAC([]byte(msg), []byte(ssecret))
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
