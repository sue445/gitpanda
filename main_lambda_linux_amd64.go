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
	"github.com/cockroachdb/errors"
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

	sentryClient, close, err := NewSentryClient(sentryDsn, isDebugLogging)
	if err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	defer close()

	response, err := lambdaMain(body)

	if err != nil {
		sentryClient.CaptureException(err)

		fmt.Printf("[ERROR] %s, error=%v\n", response, err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       response,
		}, errors.WithStack(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       response,
	}, nil
}

func lambdaMain(body string) (string, error) {
	slackOAuthAccessToken, err := GetParameterStoreOrEnv("SLACK_OAUTH_ACCESS_TOKEN", os.Getenv("SLACK_OAUTH_ACCESS_TOKEN_KEY"), true)
	if err != nil {
		return "Failed: slackOAuthAccessToken", errors.WithStack(err)
	}

	slackVerificationToken, err := GetParameterStoreOrEnv("SLACK_VERIFICATION_TOKEN", os.Getenv("SLACK_VERIFICATION_TOKEN_KEY"), false)
	if err != nil {
		return "Failed: slackVerificationToken", errors.WithStack(err)
	}

	gitlabAPIEndpoint, err := GetParameterStoreOrEnv("GITLAB_API_ENDPOINT", os.Getenv("GITLAB_API_ENDPOINT_KEY"), true)
	if err != nil {
		return "Failed: gitlabAPIEndpoint", errors.WithStack(err)
	}

	gitlabBaseURL, err := GetParameterStoreOrEnv("GITLAB_BASE_URL", os.Getenv("GITLAB_BASE_URL_KEY"), true)
	if err != nil {
		return "Failed: gitlabBaseURL", errors.WithStack(err)
	}

	gitlabPrivateToken, err := GetParameterStoreOrEnv("GITLAB_PRIVATE_TOKEN", os.Getenv("GITLAB_PRIVATE_TOKEN_KEY"), true)
	if err != nil {
		return "Failed: gitlabPrivateToken", errors.WithStack(err)
	}

	s := webhook.NewSlackWebhook(
		slackOAuthAccessToken,
		slackVerificationToken,
		&gitlab.URLParserParams{
			APIEndpoint:     gitlabAPIEndpoint,
			BaseURL:         gitlabBaseURL,
			PrivateToken:    gitlabPrivateToken,
			GitPandaVersion: Version,
			IsDebugLogging:  isDebugLogging,
		},
	)
	response, err := s.Request(body, truncateLines)

	if err != nil {
		return response, errors.WithStack(err)
	}

	return response, nil
}

// GetParameterStoreOrEnv returns Environment variable or Parameter Store variable
func GetParameterStoreOrEnv(envKey string, parameterStoreKey string, required bool) (string, error) {
	if parameterStoreKey != "" {
		// Get from Parameter Store
		d, err := NewSsmDecrypter()

		if err != nil {
			return "", errors.WithStack(err)
		}

		decryptedValue, err := d.Decrypt(parameterStoreKey)
		if err != nil {
			return fmt.Sprintf("Failed: Decrypt %s", parameterStoreKey), errors.WithStack(err)
		}

		return decryptedValue, nil
	}

	// Get from env
	if os.Getenv(envKey) != "" {
		return os.Getenv(envKey), nil
	}

	if required {
		return "", fmt.Errorf("Either %s or %s is required", envKey, parameterStoreKey)
	}

	return "", nil
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
		return nil, errors.WithStack(err)
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
		return "", errors.WithStack(err)
	}
	return *resp.Parameter.Value, nil
}
