# Changelog

All notable changes to this package are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project follows
[Semantic Versioning](https://semver.org/).

## [1.1.0]

### Added
- `initializers/scheduler`: cron runtime shared by the services — a scheduler that understands
  seconds and skips overlapping runs, job registration bound to the application context, and a
  per-job failure tracker that keeps a single transient timeout at warning level and escalates
  to error once the failures pile up.

## [1.0.0]

First release under Nimble Tech.

### Added
- golangci-lint config, GitHub Actions CI and dependabot.

### Changed
- Naming aligned with Go conventions: `Options.URL`, `JSONAPIErrors.GetHTTPStatus`,
  `AddContext(ctx, err)`.
- Error helper methods use pointer receivers consistently.

### Fixed
- `RabbitMQ.tryReconnect` now honours `ConnectAttempts` instead of looping
  forever, and unreachable code was removed.
- Connection/channel close errors are explicitly handled.

### Removed
- `initializers/amqp_v2` moved to the `go-events` package (it wired go-events
  types and caused a module dependency cycle).
