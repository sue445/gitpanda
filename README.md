# gitpanda :panda_face:
**Git**Lab URL ex**pande**r for Slack

[![Latest Version](https://img.shields.io/github/v/release/sue445/gitpanda)](https://github.com/sue445/gitpanda/releases)
[![CircleCI](https://circleci.com/gh/sue445/gitpanda.svg?style=svg)](https://circleci.com/gh/sue445/gitpanda)
[![docker](https://github.com/sue445/gitpanda/actions/workflows/docker.yml/badge.svg)](https://github.com/sue445/gitpanda/actions/workflows/docker.yml)
[![Maintainability](https://api.codeclimate.com/v1/badges/003d4dd72d10220e2564/maintainability)](https://codeclimate.com/github/sue445/gitpanda/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/003d4dd72d10220e2564/test_coverage)](https://codeclimate.com/github/sue445/gitpanda/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/sue445/gitpanda)](https://goreportcard.com/report/github.com/sue445/gitpanda)

## Example
![example1](img/example1.png)

![example2](img/example2.png)

## Requirements
* GitLab API v4
* Slack app and OAuth Access Token
  * see [CREATE_SLACK_APP.md](CREATE_SLACK_APP.md)

## Supported URL format
* User URL
  * e.g. `${GITLAB_BASE_URL}/:username`
* Group URL
  * e.g. `${GITLAB_BASE_URL}/:groupname`
* Project URL
  * e.g. `${GITLAB_BASE_URL}/:namespace/:reponame`
* Issue URL
  * e.g. `${GITLAB_BASE_URL}/:namespace/:reponame/issues/:iid`
* MergeRequest URL
  * e.g. `${GITLAB_BASE_URL}/:namespace/:reponame/merge_requests/:iid`
* Job URL
  * e.g. `${GITLAB_BASE_URL}/:namespace/:reponame/-/jobs/:id`
* Pipeline URL
  * e.g. `${GITLAB_BASE_URL}/:namespace/:reponame/pipelines/:id`
* Blob URL
  * e.g. `${GITLAB_BASE_URL}/:namespace/:reponame/blob/:sha1/:filename`
* Project snippet URL
  * e.g. `${GITLAB_BASE_URL}/:namespace/:reponame/snippets/:id`
* Snippet URL
  * e.g. `${GITLAB_BASE_URL}/snippets/:id`

## Running standalone
Download latest binary from https://github.com/sue445/gitpanda/releases

```bash
PORT=8000 \
GITLAB_API_ENDPOINT=https://your-gitlab.example.com/api/v4 \
GITLAB_BASE_URL=https://your-gitlab.example.com \
GITLAB_PRIVATE_TOKEN=xxxxxxxxxx \
SLACK_OAUTH_ACCESS_TOKEN=xoxp-0000000000-0000000000-000000000000-00000000000000000000000000000000 \
SLACK_VERIFICATION_TOKEN=xxxxxxxxx \
TRUNCATE_LINES=5 \
./gitpanda
```

### Environment variables
* `PORT`
  * default is `8000`
* `GITLAB_API_ENDPOINT`
  * e.g. `https://your-gitlab.example.com/api/v4`
* `GITLAB_BASE_URL`
  * e.g. `https://your-gitlab.example.com`
* `GITLAB_PRIVATE_TOKEN`
  * Generate a personal access token with `api` scope
* `SLACK_OAUTH_ACCESS_TOKEN`
  * see [CREATE_SLACK_APP.md](CREATE_SLACK_APP.md)
  * e.g. `xoxp-0000000000-0000000000-000000000000-00000000000000000000000000000000`
* `SLACK_VERIFICATION_TOKEN`
  * Token for verify slack requests. This is optional, but **recommended**
  * see. https://api.slack.com/docs/verifying-requests-from-slack#app_management_updates
* `TRUNCATE_LINES`
  * Line count to truncate the text (default. no truncate)
* `SENTRY_DSN`
  * [Sentry](https://sentry.io/) DSN
  * e.g. `https://xxxxxxxxxxxxx@sentry.example.com/0000000`

## Running with AWS (Lambda + API Gateway + Parameter Store)
Use latest `gitpanda_linux_amd64` on https://github.com/sue445/gitpanda/releases

### Environment variables
One of the following is required

| Environment                | Key of Parameter Store         |
| -------------------------- | ------------------------------ |
| `GITLAB_API_ENDPOINT`      | `GITLAB_API_ENDPOINT_KEY`      |
| `GITLAB_BASE_URL`          | `GITLAB_BASE_URL_KEY`          |
| `GITLAB_PRIVATE_TOKEN`     | `GITLAB_PRIVATE_TOKEN_KEY`     |
| `SLACK_OAUTH_ACCESS_TOKEN` | `SLACK_OAUTH_ACCESS_TOKEN_KEY` |
| `SLACK_VERIFICATION_TOKEN` | `SLACK_VERIFICATION_TOKEN_KEY` |
| `TRUNCATE_LINES`           |                                |

When you want to store to Parameter Store, please store as `SecureString`

![aws parameter_store](img/aws-parameter_store.png)

### Example
* [Example template for AWS SAM](examples/aws_sam_template.yaml)

## Arguments
```bash
$ ./gitpanda --help
Usage of ./gitpanda:
  -version
    	Whether showing version
```

## Running with docker
Run latest version

```bash
docker run --rm -it ghcr.io/sue445/gitpanda
```

Run with specified version

```bash
docker run --rm -it ghcr.io/sue445/gitpanda:vX.Y.Z
```

Available tags are followings

https://github.com/users/sue445/packages/container/package/gitpanda

## Development
Recommend to use https://github.com/direnv/direnv

```bash
cp .envrc.example .envrc
vi .envrc
direnv allow
```

## Heroku deploy
[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)
