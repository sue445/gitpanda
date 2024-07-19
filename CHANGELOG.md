## Unreleased
[full changelog](http://github.com/sue445/gitpanda/compare/v0.12.6...main)

## [v0.12.6](https://github.com/sue445/gitpanda/releases/tag/v0.12.6)
[full changelog](http://github.com/sue445/gitpanda/compare/v0.12.5...v0.12.6)

* Migrate to github.com/aws/aws-sdk-go-v2
  * https://github.com/sue445/gitpanda/pull/1586
* Upgrade dependencies

## [v0.12.5](https://github.com/sue445/gitpanda/releases/tag/v0.12.5)
[full changelog](http://github.com/sue445/gitpanda/compare/v0.12.4...v0.12.5)

* Upgrade to Go 1.22 :rocket:
  * https://github.com/sue445/gitpanda/pull/1467
* Add golangci-lint
  * https://github.com/sue445/gitpanda/pull/1580
* Upgrade dependencies

## v0.12.4
[full changelog](http://github.com/sue445/gitpanda/compare/v0.12.3...v0.12.4)

* ignore URL parameters in blobFetcher
  * https://github.com/sue445/gitpanda/pull/1307
* Upgrade dependencies

## v0.12.3
[full changelog](http://github.com/sue445/gitpanda/compare/v0.12.2...v0.12.3)

* Upgrade to Go 1.21 :rocket:
  * https://github.com/sue445/gitpanda/pull/1304
* Upgrade dependencies

## v0.12.2
[full changelog](http://github.com/sue445/gitpanda/compare/v0.12.1...v0.12.2)

* Upgrade to Go 1.20 :rocket:
  * https://github.com/sue445/gitpanda/pull/1135
* Migrate to github.com/cockroachdb/errors
  * https://github.com/sue445/gitpanda/pull/1275
* Upgrade dependencies

## v0.12.1
[full changelog](http://github.com/sue445/gitpanda/compare/v0.12.0...v0.12.1)

* Re-released due to unintentional disappearance of `us-docker.pkg.dev/gitpanda/gitpanda/app:v0.12.0` and `us-docker.pkg.dev/gitpanda/gitpanda/app:latest`
* Upgrade dependencies

## v0.12.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.11.0...v0.12.0)

* Support commit URL
  * https://github.com/sue445/gitpanda/pull/1062
* Fixed. 404 error for URLs containing anchors
  * https://github.com/sue445/gitpanda/pull/1061

## v0.11.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.10.2...v0.11.0)

* Upgrade to Go 1.19
  * https://github.com/sue445/gitpanda/pull/985
* Push docker image to GCP Artifact Registry
  * https://github.com/sue445/gitpanda/pull/1010
* Upgrade dependencies

## v0.10.2
[full changelog](http://github.com/sue445/gitpanda/compare/v0.10.1...v0.10.2)

* update nlopes/slack@v0.6.0 to slack-go/slack@v0.10.3
  * https://github.com/sue445/gitpanda/pull/923
* Upgrade dependencies

## v0.10.1
[full changelog](http://github.com/sue445/gitpanda/compare/v0.10.0...v0.10.1)

* Fixed failure to send error to Sentry when not Lambda
  * https://github.com/sue445/gitpanda/pull/896
* Upgrade to Go 1.18
  * https://github.com/sue445/gitpanda/pull/874
* Add Kubernetes examples
  * https://github.com/sue445/gitpanda/pull/895
* Upgrade dependencies

## v0.10.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.9.6...v0.10.0)

* Supports to job URL with the line number
  * https://github.com/sue445/gitpanda/pull/668
* Upgrade dependencies

## v0.9.6
[full changelog](http://github.com/sue445/gitpanda/compare/v0.9.5...v0.9.6)

* Fixed.`x509: certificate signed by unknown authority`
  * https://github.com/sue445/gitpanda/pull/662
* Upgrade dependencies

## v0.9.5
[full changelog](http://github.com/sue445/gitpanda/compare/v0.9.4...v0.9.5)

* Support GitHub Container Registry (ghcr.io)
  * https://github.com/sue445/gitpanda/pull/594
* Upgrade to Go 1.16
  * https://github.com/sue445/gitpanda/pull/554
* Upgrade dependencies

## v0.9.4
[full changelog](http://github.com/sue445/gitpanda/compare/v0.9.3...v0.9.4)

* Wrap all error with `errors.WithStack`
  * https://github.com/sue445/gitpanda/pull/490
* Upgrade dependencies

## v0.9.3
[full changelog](http://github.com/sue445/gitpanda/compare/v0.9.2...v0.9.3)

* Support snippet routes including `/-/` for GitLab 13.3.0+
  * https://github.com/sue445/gitpanda/pull/446
* Upgrade dependencies

## v0.9.2
[full changelog](http://github.com/sue445/gitpanda/compare/v0.9.1...v0.9.2)

* Expand only UTF-8 file
  * https://github.com/sue445/gitpanda/pull/430
* Upgrade to Go 1.15
  * https://github.com/sue445/gitpanda/pull/429
* Upgrade dependencies

## v0.9.1
[full changelog](http://github.com/sue445/gitpanda/compare/v0.9.0...v0.9.1)

* Support blob link without line hash
  * https://github.com/sue445/gitpanda/pull/320
* Upgrade to Go 1.14
  * https://github.com/sue445/gitpanda/pull/321
* Upgrade dependencies

## v0.9.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.8.0...v0.9.0)

* Support project routes including `/-/`
  * https://github.com/sue445/gitpanda/pull/251
* Upgrade dependencies

## v0.8.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.7.3...v0.8.0)

* Support [Sentry](https://sentry.io/)
  * https://github.com/sue445/gitpanda/pull/213
* Upgrade dependencies

## v0.7.3
[full changelog](http://github.com/sue445/gitpanda/compare/v0.7.2...v0.7.3)

* Build with Go 1.13
  * https://github.com/sue445/gitpanda/pull/97
* Upgrade dependencies

## v0.7.2
[full changelog](http://github.com/sue445/gitpanda/compare/v0.7.1...v0.7.2)

* Expand urls when contains both valid and invalid urls
  * https://github.com/sue445/gitpanda/pull/79

## v0.7.1
[full changelog](http://github.com/sue445/gitpanda/compare/v0.7.0...v0.7.1)

* Verify slack request with signature token
  * https://github.com/sue445/gitpanda/pull/74
* Add newline to log
  * https://github.com/sue445/gitpanda/pull/75

## v0.7.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.6.0...v0.7.0)

* Support sub-group URL
  * https://github.com/sue445/gitpanda/pull/72

## v0.6.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.5.0...v0.6.0)

* Support snippet URL
  * https://github.com/sue445/gitpanda/pull/71

## v0.5.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.4.1...v0.5.0)

* Support group URL
  * https://github.com/sue445/gitpanda/pull/68

## v0.4.1
[full changelog](http://github.com/sue445/gitpanda/compare/v0.4.0...v0.4.1)

* Bugfixed. description is broken when text or url in markdown is blank
  * https://github.com/sue445/gitpanda/pull/66

## v0.4.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.10...v0.4.0)

* Support Job URL and Pipeline URL
  * https://github.com/sue445/gitpanda/pull/64

## v0.3.10
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.9...v0.3.10)

* Add duration in debug log
  * https://github.com/sue445/gitpanda/pull/61

## v0.3.9
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.8...v0.3.9)

* Bugfix. debug logging
  * https://github.com/sue445/gitpanda/pull/59

## v0.3.8
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.7...v0.3.8)

* Convert markdown link to Slack link
  * https://github.com/sue445/gitpanda/pull/58

## v0.3.7
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.6...v0.3.7)

* Performance tuning when call multiple APIs
  * https://github.com/sue445/gitpanda/pull/57

## v0.3.6
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.5...v0.3.6)

* Sanitize embed image in description
  * https://github.com/sue445/gitpanda/pull/56

## v0.3.5
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.4...v0.3.5)

* Bugfixed. couldn't be expanded url which ends with `/`
  * https://github.com/sue445/gitpanda/pull/54

## v0.3.4
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.3...v0.3.4)

* Fixed SEGV when non-user's project
  * https://github.com/sue445/gitpanda/pull/46

## v0.3.3
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.2...v0.3.3)

* Add debug log in gitlab/url_parser
  * https://github.com/sue445/gitpanda/pull/45

## v0.3.2
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.1...v0.3.2)

* Add footer to slack attachment
  * https://github.com/sue445/gitpanda/pull/44

## v0.3.1
[full changelog](http://github.com/sue445/gitpanda/compare/v0.3.0...v0.3.1)

* Set User-Agent for GitLab API
  * https://github.com/sue445/gitpanda/pull/41
* Performance tuning for multiple links
  * https://github.com/sue445/gitpanda/pull/43
* Refactor: split to sub packages
  * https://github.com/sue445/gitpanda/pull/42
* and other refactorings

## v0.3.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.2.0...v0.3.0)

* Support blob URL
  * https://github.com/sue445/gitpanda/pull/35

## v0.2.0
[full changelog](http://github.com/sue445/gitpanda/compare/v0.1.3...v0.2.0)

* Tweak page title
  * https://github.com/sue445/gitpanda/pull/33
* Support comment url
  * https://github.com/sue445/gitpanda/pull/34

## v0.1.3
[full changelog](http://github.com/sue445/gitpanda/compare/v0.1.2...v0.1.3)

* Add `TRUNCATE_LINES`
  * https://github.com/sue445/gitpanda/pull/30

## v0.1.2
[full changelog](http://github.com/sue445/gitpanda/compare/v0.1.1...v0.1.2)

* Add debug log in response
  * https://github.com/sue445/gitpanda/pull/26
* Update dependencies
  * https://github.com/sue445/gitpanda/pull/23``

## v0.1.1
[full changelog](http://github.com/sue445/gitpanda/compare/v0.1.0...v0.1.1)

* Remove check for parameter store
  * https://github.com/sue445/gitpanda/pull/22

## v0.1.0
* Initial release
