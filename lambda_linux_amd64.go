package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"log"
	"net/http"
	"os"
	"strings"
)

func startLambda() {
	checkEnv("GITLAB_API_ENDPOINT")
	checkEnv("GITLAB_BASE_URL")
	checkEnv("GITLAB_PRIVATE_TOKEN")
	checkEnv("SLACK_OAUTH_ACCESS_TOKEN")

	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := strings.TrimSpace(request.Body)

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
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       response,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       response,
	}, nil
}
