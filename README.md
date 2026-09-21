# go-tools

[![CI](https://github.com/tech-nimble/go-tools/actions/workflows/ci.yml/badge.svg)](https://github.com/tech-nimble/go-tools/actions/workflows/ci.yml)

Shared building blocks for Nimble Tech Go services: error handling, Gin helpers,
tracing, and ready-made initializers wired from environment variables.

## Install

```sh
go get github.com/tech-nimble/go-tools
```

## Packages

| Package | What it provides |
| --- | --- |
| `helpers/errors` | Extended error type with type, code, context and i18n-aware handler. |
| `helpers/gin` | Gin middleware (auth, content-type, headers) and JSON:API renderers. |
| `helpers/jaeger` | OpenTracing span helpers, including AMQP context propagation. |
| `helpers/sentry` | Sentry hub helpers. |
| `helpers/http` | HTTP header helpers. |
| `events`, `events/queue` | RabbitMQ client, event bus and repository primitives. |
| `repositories` | pgx-based repository base with transaction helpers. |
| `initializers/*` | Bootstrap helpers for pg, redis, logs, sentry, http server, env, amqp, errors. |
| `initializers/scheduler` | Cron runtime: scheduler with seconds, job registry and failure tracking. |

## Initializers

Each `initializers/*` package reads configuration from environment variables and
returns a ready-to-use component, e.g.:

```go
pool, err := pg.Initialize()           // initializers/pg  — *pgxpool.Pool
logger := logs.InitializeLogs(writer)  // initializers/logs — zerolog.Logger
```

## Scheduler

`initializers/scheduler` holds the cron runtime a service would otherwise copy into
`internal/ports/cronjob`. A service only describes its jobs:

```go
jobs := []scheduler.Job{{
    Name:                "dispatch.expire_offers",
    Spec:                scheduler.EverySpec(cfg.TickInterval, 5*time.Second),
    Timeout:             10 * time.Second,
    Quiet:               true,
    FailuresBeforeAlarm: 3,
    Run:                 func(ctx context.Context) error { return useCase.Execute(ctx) },
}}

cronScheduler := scheduler.Initialize()
scheduler.Register(jobsCtx, cronScheduler, jobs...)
cronScheduler.Start()
```

The context passed to `Register` is the application's, not a request's: canceling it stops the
jobs in flight. The seconds field is optional, so `15 3 * * *` and `0 30 3 * * *` both register,
and a job never runs twice at once.

A failed job is reported at error level from the first failure. A job that ticks every few seconds
can raise `FailuresBeforeAlarm`: the first failures are then logged at warning level and only a run
of them escalates, so a short database outage that the next tick already cured does not raise an
alert. A rare job leaves the field alone — three silent nights is worse than one noisy alert.

`scheduler.RunOnce(ctx, name, jobs...)` runs one job outside its schedule, and `scheduler.Names`
lists what a service declared: together they let a service expose its jobs on the command line.

## License

[MIT](LICENSE) © Nimble Tech
