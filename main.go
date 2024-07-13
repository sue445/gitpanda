package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/sue445/gitpanda/gitlab"
	"github.com/sue445/gitpanda/webhook"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	// Version represents app version (injected from ldflags)
	Version string

	// Revision represents app revision (injected from ldflags)
	Revision string
)

var isPrintVersion bool

var isDebugLogging bool

var truncateLines = 0

var sentryDsn = ""

func init() {
	flag.BoolVar(&isPrintVersion, "version", false, "Whether showing version")

	flag.Parse()

	if os.Getenv("DEBUG_LOGGING") != "" {
		isDebugLogging = true
	}

	if os.Getenv("SENTRY_DSN") != "" {
		sentryDsn = os.Getenv("SENTRY_DSN")
	}

	if s := os.Getenv("TRUNCATE_LINES"); s != "" {
		lines, err := strconv.Atoi(s)

		if err != nil {
			fmt.Printf("%s is invalid number, error=%v\n", s, err)
			os.Exit(1)
		}

		if lines > 0 {
			truncateLines = lines
		}
	}
}

func checkEnv(name string) {
	if os.Getenv(name) != "" {
		return
	}

	fmt.Printf("[ERROR] %s is required\n", name)
	fmt.Println("")
	printUsage()
	os.Exit(1)
}

func printVersion() {
	fmt.Printf("gitpanda %s, build %s\n", Version, Revision)
}

func printUsage() {
	fmt.Println("[Usage]")
	fmt.Println("  PORT=8000 \\")
	fmt.Println("  GITLAB_API_ENDPOINT=https://your-gitlab.example.com/api/v4 \\")
	fmt.Println("  GITLAB_BASE_URL=https://your-gitlab.example.com \\")
	fmt.Println("  GITLAB_PRIVATE_TOKEN=xxxxxxxxxx \\")
	fmt.Println("  SLACK_OAUTH_ACCESS_TOKEN=xoxp-0000000000-0000000000-000000000000-00000000000000000000000000000000 \\")
	fmt.Println("  ./gitpanda")
}

func normalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")

	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("It works"))
		if err != nil {
			fmt.Printf("[ERROR] error=%v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error())) //nolint:errcheck
			panic(err)
		}

	case http.MethodPost:
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			fmt.Printf("[ERROR] error=%v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error())) //nolint:errcheck
			panic(err)
		}

		body := strings.TrimSpace(buf.String())

		if isDebugLogging {
			fmt.Printf("[DEBUG] normalHandler: body=%s\n", body)
		}

		s := webhook.NewSlackWebhook(
			os.Getenv("SLACK_OAUTH_ACCESS_TOKEN"),
			os.Getenv("SLACK_VERIFICATION_TOKEN"),
			&gitlab.URLParserParams{
				APIEndpoint:     os.Getenv("GITLAB_API_ENDPOINT"),
				BaseURL:         os.Getenv("GITLAB_BASE_URL"),
				PrivateToken:    os.Getenv("GITLAB_PRIVATE_TOKEN"),
				GitPandaVersion: Version,
				IsDebugLogging:  isDebugLogging,
			},
		)
		response, err := s.Request(body, truncateLines)

		if err != nil {
			fmt.Printf("[ERROR] body=%s, response=%s, error=%v\n", body, response, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error())) //nolint:errcheck
			panic(err)
		}

		_, err = w.Write([]byte(response))
		if err != nil {
			fmt.Printf("[ERROR] body=%s, response=%s, error=%v\n", body, response, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error())) //nolint:errcheck
			panic(err)
		}
	}
}
