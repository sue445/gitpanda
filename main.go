package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	// Version represents app version (injected from ldflags)
	Version string

	// Revision represents app revision (injected from ldflags)
	Revision string
)

var port int

var isPrintVersion bool

func init() {
	flag.BoolVar(&isPrintVersion, "version", false, "Whether showing version")

	flag.Parse()
}

func main() {
	if isPrintVersion {
		printVersion()
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("gitpanda started: port=%s\n", port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}

func printVersion() {
	fmt.Printf("gitpanda v%s, build %s\n", Version, Revision)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("It works"))

	case http.MethodPost:
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := strings.TrimSpace(buf.String())

		s := NewSlackWebhook(
			os.Getenv("SLACK_OAUTH_ACCESS_TOKEN"),
			&GitLabURLParserParams{
				APIEndpoint:  os.Getenv("GITLAB_API_ENDPOINT"),
				BaseURL:      os.Getenv("GITLAB_BASE_URL"),
				PrivateToken: os.Getenv("GITLAB_PRIVATE_TOKEN"),
			},
		)
		response, err := s.Request(
			body,
			false,
		)

		if err != nil {
			log.Printf("[ERROR] body=%s, response=%s, error=%v", body, response, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(response))
	}
}
