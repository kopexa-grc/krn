# Changelog

## 1.0.0 (2026-01-14)


### ⚠ BREAKING CHANGES

* Remove service from KRN format

### Features

* add optional service subdomain support ([1771dd7](https://github.com/kopexa-grc/krn/commit/1771dd772a5446fdfd5b8063ba576d283ccee8cf))
* initial KRN implementation ([fd7d053](https://github.com/kopexa-grc/krn/commit/fd7d053a6cf5bcc5d490240eb33cecc09e2248d1))


### Bug Fixes

* add fallback token support for release-please ([86626a7](https://github.com/kopexa-grc/krn/commit/86626a7c081a9f9065111dbd2cfc9fc4bc595b2c))
* CI workflow compatibility issues ([1cb0b22](https://github.com/kopexa-grc/krn/commit/1cb0b2217037a95b775177f7995786b63ae46953))
* convert if-else chain to switch for gocritic linter ([8f9b47e](https://github.com/kopexa-grc/krn/commit/8f9b47ed003363862307d189a4defc41a31df890))
* exclude gocyclo from test files ([2020c2d](https://github.com/kopexa-grc/krn/commit/2020c2dba0e81efacb7e40d067239f626dd09ce9))
* relax coverage threshold to 90% and fix formatting ([fe02f57](https://github.com/kopexa-grc/krn/commit/fe02f5724a048d4e034514b47cc773d82d4ae373))
* update golangci-lint config for v2 format ([9b3f617](https://github.com/kopexa-grc/krn/commit/9b3f6173f276a01d1cc4547a31e7e216a97121ee))
* use KOPEXA_CLOUD_REPO_ACCESS_TOKEN for releases ([ad47c93](https://github.com/kopexa-grc/krn/commit/ad47c939a660b4cdf2134c416a57a11a351bc8f4))


### Code Refactoring

* simplify KRN format to //kopexa.com/{collection}/{id} ([203e98c](https://github.com/kopexa-grc/krn/commit/203e98cd60a34a0ef07c1745a1019f7297f4a859))
