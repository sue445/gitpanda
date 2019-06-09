package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"log"
	"net/http"
	"os"
	"strings"
)

func startLambda() {
	checkEnv("GITLAB_API_ENDPOINT_KEY")
	checkEnv("GITLAB_BASE_URL_KEY")
	checkEnv("GITLAB_PRIVATE_TOKEN_KEY")
	checkEnv("SLACK_OAUTH_ACCESS_TOKEN_KEY")

	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := strings.TrimSpace(request.Body)

	d, err := NewSsmDecrypter()

	if err != nil {
		fmt.Printf("Failed: NewSsmDecrypter, error=%v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed: NewSsmDecrypter",
		}, err
	}

	slackOAuthAccessToken, err := d.Decrypt(os.Getenv("SLACK_OAUTH_ACCESS_TOKEN_KEY"))
	if err != nil {
		fmt.Printf("Failed: Decrypt SLACK_OAUTH_ACCESS_TOKEN_KEY, error=%v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed: Decrypt SLACK_OAUTH_ACCESS_TOKEN_KEY",
		}, err
	}

	gitlabAPIEndpoint, err := d.Decrypt(os.Getenv("GITLAB_API_ENDPOINT_KEY"))
	if err != nil {
		fmt.Printf("Failed: Decrypt GITLAB_API_ENDPOINT_KEY, error=%v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed: Decrypt GITLAB_API_ENDPOINT_KEY",
		}, err
	}

	gitlabBaseURL, err := d.Decrypt(os.Getenv("GITLAB_BASE_URL_KEY"))
	if err != nil {
		fmt.Printf("Failed: Decrypt GITLAB_BASE_URL_KEY, error=%v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed: Decrypt GITLAB_BASE_URL_KEY",
		}, err
	}

	gitlabPrivateToken, err := d.Decrypt(os.Getenv("GITLAB_API_ENDPOINT_KEY"))
	if err != nil {
		fmt.Printf("Failed: Decrypt GITLAB_API_ENDPOINT_KEY, error=%v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed: Decrypt GITLAB_API_ENDPOINT_KEY",
		}, err
	}

	s := NewSlackWebhook(
		slackOAuthAccessToken,
		&GitLabURLParserParams{
			APIEndpoint:  gitlabAPIEndpoint,
			BaseURL:      gitlabBaseURL,
			PrivateToken: gitlabPrivateToken,
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

// SsmDecrypter stores the AWS Session used for SSM decrypter.
type SsmDecrypter struct {
	sess *session.Session
	svc  ssmiface.SSMAPI
}

// NewSsmDecrypter returns a new SsmDecrypter.
func NewSsmDecrypter() (*SsmDecrypter, error) {
	sess, err := session.NewSession()

	if err != nil {
		return nil, err
	}

	svc := ssm.New(sess)
	return &SsmDecrypter{sess, svc}, nil
}

// Decrypt decrypts string.
func (d *SsmDecrypter) Decrypt(encrypted string) (string, error) {
	params := &ssm.GetParameterInput{
		Name:           aws.String(encrypted),
		WithDecryption: aws.Bool(true),
	}
	resp, err := d.svc.GetParameter(params)
	if err != nil {
		return "", err
	}
	return *resp.Parameter.Value, nil
}
