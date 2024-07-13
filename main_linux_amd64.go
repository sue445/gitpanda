package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if isPrintVersion {
		printVersion()
		return
	}

	if os.Getenv("LAMBDA_TASK_ROOT") != "" && os.Getenv("LAMBDA_RUNTIME_DIR") != "" {
		// for AWS Lambda
		startLambda()
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

	sentryClient, releaseSentry, err := NewSentryClient(sentryDsn, isDebugLogging)
	if err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	defer releaseSentry()

	fmt.Printf("gitpanda started: port=%s\n", port)
	http.HandleFunc("/", sentryClient.HandleFunc(normalHandler))

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
