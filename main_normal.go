// +build !linux !amd64

package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if isPrintVersion {
		printVersion()
		return
	}

	checkEnv("GITLAB_API_ENDPOINT")
	checkEnv("GITLAB_BASE_URL")
	checkEnv("GITLAB_PRIVATE_TOKEN")
	checkEnv("SLACK_OAUTH_ACCESS_TOKEN")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	sentryClient, close, err := NewSentryClient(sentryDsn, isDebugLogging)
	if err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	defer close()

	fmt.Printf("gitpanda started: port=%s\n", port)
	http.HandleFunc("/", sentryClient.HandleFunc(normalHandler))
	http.ListenAndServe(":"+port, nil)
}
