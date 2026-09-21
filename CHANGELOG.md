# Changelog

All notable changes to this package are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project follows
[Semantic Versioning](https://semver.org/).

## [1.2.0]

### Added
- `initializers/scheduler`: `RunOnce` and `Names`, so a service can run a single job from the
  command line instead of waiting for its next tick.

### Changed
- `Job.TransientFailures` is now `Job.FailuresBeforeAlarm` and defaults to 1: a failure is
  reported at error level straight away unless the job asks for more patience. Only frequent
  jobs should raise it.
- `EverySpec` takes the fallback interval from the caller instead of assuming five seconds.

## [1.1.1]

### Fixed
- `initializers/scheduler`: the seconds field is optional again, so a five-field spec such as
  `15 3 * * *` registers alongside `0 30 3 * * *`.

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
