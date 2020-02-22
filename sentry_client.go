package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"net/http"
	"time"
)

// SentryClient wraps sentry
type SentryClient struct {
	dsn     string
	handler *sentryhttp.Handler
}

// Close must be called after NewSentryClient
type Close func()

// NewSentryClient returns new SentryClient instance
func NewSentryClient(dsn string, debug bool) (*SentryClient, Close, error) {
	if dsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:     dsn,
			Release: fmt.Sprintf("gitpanda@%s", Version),
			Debug:   debug,
		})

		if err != nil {
			return nil, nil, err
		}
	}

	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	return &SentryClient{dsn: dsn, handler: sentryHandler}, releaseSentry, nil
}

func releaseSentry() {
	sentry.Flush(2 * time.Second)
}

// HandleFunc wraps handler if necessary
func (s *SentryClient) HandleFunc(handler http.HandlerFunc) http.HandlerFunc {
	if s.dsn == "" {
		return handler
	}

	return s.handler.HandleFunc(handler)
}

// CaptureException send error to sentry if necessary
func (s *SentryClient) CaptureException(err error) {
	if s.dsn != "" {
		sentry.CaptureException(err)
	}
}
