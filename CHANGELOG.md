## Unreleased
[full changelog](http://github.com/sue445/gitpanda/compare/v0.9.1...master)

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
