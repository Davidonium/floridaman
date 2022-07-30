package server

import (
	"crypto/hmac"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/davidonium/floridaman/internal/floridaman"
	"github.com/davidonium/floridaman/internal/util"
)

type slackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

var ErrInvalidSlackRequest = errors.New("invalid slack request")

func (s *Server) slackRandomArticleHandler(
	logger *log.Logger,
	ar floridaman.ArticleReader,
	slackSecret string,
) APIHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Printf("error reading body %v", err)
			return ErrInvalidSlackRequest
		}

		defer func() {
			err := r.Body.Close()
			if err != nil {
				s.logger.Printf("failed to close response body: %v", err)
			}
		}()

		ok := ValidateSlackRequest(r, logger, slackSecret, body)

		if !ok {
			return ErrInvalidSlackRequest
		}

		article, err := ar.Random()
		if err != nil {
			return err
		}

		response := slackResponse{
			Text:         fmt.Sprintf("%s (%s)", article.Title, article.Link),
			ResponseType: "in_channel",
		}

		return s.writeJSON(w, response)
	}
}

func (s *Server) oauthSlackRedirectHandler() APIHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		_, err := w.Write([]byte("floridaman oauth success!"))
		if err != nil {
			return err
		}

		return nil
	}
}

func ValidateSlackRequest(r *http.Request, logger *log.Logger, ssecret string, response []byte) bool {
	ssig := r.Header.Get("X-Slack-Signature")
	t := r.Header.Get("X-Slack-Request-Timestamp")

	ts, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		logger.Printf("X-Slack-Request-Timestamp is not a parsable number, got \"%s\"\n", t)
		return false
	}

	tsu := time.Unix(ts, 0)

	if time.Since(tsu) > 5*time.Minute {
		logger.Println("timestamp difference is greater than 5 minutes")
		return false
	}
	msg := fmt.Sprintf("v0:%s:%s", t, response)

	sig := util.HMACString([]byte(msg), []byte(ssecret))
	ok := hmac.Equal([]byte(sig), []byte(ssig))

	if !ok {
		logger.Printf("error validating hmac signature from slack: %s, generated %s\n", ssig, sig)
		return false
	}

	return true
}
