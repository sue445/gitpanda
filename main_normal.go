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

	fmt.Printf("gitpanda started: port=%s\n", port)
	http.HandleFunc("/", normalHandler)
	http.ListenAndServe(":"+port, nil)
}
