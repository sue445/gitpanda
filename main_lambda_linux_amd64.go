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
	"github.com/sue445/gitpanda/gitlab"
	"github.com/sue445/gitpanda/webhook"
	"net/http"
	"os"
	"strings"
)

func startLambda() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := strings.TrimSpace(request.Body)

	if isDebugLogging {
		fmt.Printf("[DEBUG] lambdaHandler: body=%s\n", body)
	}

	response, err := lambdaMain(body)

	if err != nil {
		fmt.Printf("[ERROR] %s, error=%v\n", response, err)
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

func lambdaMain(body string) (string, error) {
	slackOAuthAccessToken, err := GetParameterStoreOrEnv("SLACK_OAUTH_ACCESS_TOKEN", os.Getenv("SLACK_OAUTH_ACCESS_TOKEN_KEY"))
	if err != nil {
		return "Failed: slackOAuthAccessToken", err
	}

	gitlabAPIEndpoint, err := GetParameterStoreOrEnv("GITLAB_API_ENDPOINT", os.Getenv("GITLAB_API_ENDPOINT_KEY"))
	if err != nil {
		return "Failed: gitlabAPIEndpoint", err
	}

	gitlabBaseURL, err := GetParameterStoreOrEnv("GITLAB_BASE_URL", os.Getenv("GITLAB_BASE_URL_KEY"))
	if err != nil {
		return "Failed: gitlabBaseURL", err
	}

	gitlabPrivateToken, err := GetParameterStoreOrEnv("GITLAB_PRIVATE_TOKEN", os.Getenv("GITLAB_PRIVATE_TOKEN_KEY"))
	if err != nil {
		return "Failed: gitlabPrivateToken", err
	}

	s := webhook.NewSlackWebhook(
		slackOAuthAccessToken,
		&gitlab.GitLabURLParserParams{
			APIEndpoint:     gitlabAPIEndpoint,
			BaseURL:         gitlabBaseURL,
			PrivateToken:    gitlabPrivateToken,
			GitPandaVersion: Version,
		},
	)
	response, err := s.Request(body, truncateLines, isDebugLogging)

	if err != nil {
		return response, err
	}

	return response, nil
}

// GetParameterStoreOrEnv returns Environment variable or Parameter Store variable
func GetParameterStoreOrEnv(envKey string, parameterStoreKey string) (string, error) {
	if parameterStoreKey != "" {
		// Get from Parameter Store
		d, err := NewSsmDecrypter()

		if err != nil {
			return "", err
		}

		decryptedValue, err := d.Decrypt(parameterStoreKey)
		if err != nil {
			return fmt.Sprintf("Failed: Decrypt %s", parameterStoreKey), err
		}

		return decryptedValue, nil
	}

	// Get from env
	if os.Getenv(envKey) != "" {
		return os.Getenv(envKey), nil
	}

	return "", fmt.Errorf("Either %s or %s is required", envKey, parameterStoreKey)
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
