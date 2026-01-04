# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-01-05

### Added

- **International mobile number validation**: `SendMessage` now validates recipient numbers against known mobile prefixes for 180+ countries worldwide
- New `IsValidMobileNumber()` function for validating mobile numbers from any supported country
- New `mobile_prefixes.go` containing mobile prefix data sourced from [Wikipedia](https://en.wikipedia.org/wiki/List_of_mobile_telephone_prefixes_by_country)
- New `ErrInvalidMobileNumber` sentinel error for use with `errors.Is()`

### Changed

- `SendMessage` now rejects numbers that don't match known mobile prefixes (previously no validation)

### Deprecated

- `IsValidAustralianMobile()` - use `IsValidMobileNumber()` instead, which supports all countries

## [1.0.0] - 2025-01-04

### Added

- Initial release
- OAuth2 authentication with automatic token caching and refresh
- Messages API (send, get, update, delete, tags)
- Virtual Numbers API (assign, get, update, delete, optouts)
- Free Trial Numbers API
- Reports API
- Health Check API
- Sender Names API
- Logs API
- Custom error types with `errors.Is()` support

[1.1.0]: https://github.com/Its-donkey/gopholler/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/Its-donkey/gopholler/releases/tag/v1.0.0
