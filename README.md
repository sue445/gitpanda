# gitpanda :panda_face:
GitLab URL expander for Slack

![example](img/example.png)

[![CircleCI](https://circleci.com/gh/sue445/gitpanda/tree/master.svg?style=svg&circle-token=f42c3df848d11f83347750c71494c0e14e7732dc)](https://circleci.com/gh/sue445/gitpanda/tree/master)

## Requirements
* GitLab API v4
* Slack app and OAuth Access Token
  * see [CREATE_SLACK_APP.md](CREATE_SLACK_APP.md)
  
## Running
Download latest binary from https://github.com/sue445/gitpanda/releases

```bash
PORT=8000 \
GITLAB_API_ENDPOINT=https://your-gitlab.example.com/api/v4 \
GITLAB_BASE_URL=https://your-gitlab.example.com \
GITLAB_PRIVATE_TOKEN=xxxxxxxxxx \
SLACK_OAUTH_ACCESS_TOKEN=xoxp-0000000000-0000000000-000000000000-00000000000000000000000000000000 \
./gitpanda
```

## Environment variables
* `PORT`
  * default is `8000`
* `GITLAB_API_ENDPOINT`
  * e.g. `https://your-gitlab.example.com/api/v4`
* `GITLAB_BASE_URL`
  * e.g. `https://your-gitlab.example.com`
* `GITLAB_PRIVATE_TOKEN`
* `SLACK_OAUTH_ACCESS_TOKEN`
  * see [CREATE_SLACK_APP.md](CREATE_SLACK_APP.md)
  * e.g. `xoxp-0000000000-0000000000-000000000000-00000000000000000000000000000000`

## Arguments
```bash
$ ./gitpanda --help
Usage of ./gitpanda:
  -version
    	Whether showing version
```

## Development
Recommend to use https://github.com/direnv/direnv

```bash
cp .envrc.example .envrc
vi .envrc
direnv allow
```
