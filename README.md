# gitpanda :panda_face:
**Git**Lab URL ex**pande**r for Slack

[![Latest Version](https://img.shields.io/github/v/release/sue445/gitpanda)](https://github.com/sue445/gitpanda/releases)
[![test](https://github.com/sue445/gitpanda/actions/workflows/test.yml/badge.svg)](https://github.com/sue445/gitpanda/actions/workflows/test.yml)
[![build](https://github.com/sue445/gitpanda/actions/workflows/build.yml/badge.svg)](https://github.com/sue445/gitpanda/actions/workflows/build.yml)
[![docker-ghcr](https://github.com/sue445/gitpanda/actions/workflows/docker-ghcr.yml/badge.svg)](https://github.com/sue445/gitpanda/actions/workflows/docker-ghcr.yml)
[![docker-gcp](https://github.com/sue445/gitpanda/actions/workflows/docker-gcp.yml/badge.svg)](https://github.com/sue445/gitpanda/actions/workflows/docker-gcp.yml)
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
* Commit URL
  * e.g. `${GITLAB_BASE_URL}/:namespace/:reponame/commit/:sha1`
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

#### Images
* [GitHub Container Registry](https://github.com/sue445/gitpanda/pkgs/container/gitpanda) **(Recommended)**
  * `ghcr.io/sue445/gitpanda:latest`: Use latest version
  * `ghcr.io/sue445/gitpanda:vX.Y.Z`: Use specified version
* [Google Artifact Registry](https://console.cloud.google.com/artifacts/docker/gitpanda/us/gitpanda/app): If you want to run this app on [Cloud Run](https://cloud.google.com/run), use this image
  * `us-docker.pkg.dev/gitpanda/gitpanda/app:latest`: Use latest version
  * `us-docker.pkg.dev/gitpanda/gitpanda/app:vX.Y.Z`: Use specified version
  * `us-docker.pkg.dev/gitpanda/gitpanda/app:edge`: The contents of the main branch are pushed to this tag

## Running with Kubernetes
### 1. [GKE] Create static IP and managed certificate if necessary
c.f https://cloud.google.com/kubernetes-engine/docs/how-to/managed-certs

### 2. Download and edit files in [examples/kubernetes_manifests/](examples/kubernetes_manifests/)

### 3. Create secrets
e.g.

```bash
kubectl create secret generic gitpanda-secret \
    --from-literal=gitlab_private_token=${GITLAB_PRIVATE_TOKEN} \
    --from-literal=slack_oauth_access_token=${SLACK_OAUTH_ACCESS_TOKEN} \
    --from-literal=slack_verification_token=${SLACK_VERIFICATION_TOKEN}
```

### 4. apply manifests
e.g.

```bash
kubectl apply -f path/to/kubernetes_manifests/*.yaml
```

## Running with Cloud Run
See followings

* [.github/workflows/docker-gcp.yml](.github/workflows/docker-gcp.yml) `deploy_to_cloud_run` job
* [.github/cloud_run-service.yml](.github/cloud_run-service.yml)

## Development
Recommend to use https://github.com/direnv/direnv

```bash
cp .envrc.example .envrc
vi .envrc
direnv allow
```
