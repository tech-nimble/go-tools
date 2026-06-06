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

## Initializers

Each `initializers/*` package reads configuration from environment variables and
returns a ready-to-use component, e.g.:

```go
pool, err := pg.Initialize()           // initializers/pg  — *pgxpool.Pool
logger := logs.InitializeLogs(writer)  // initializers/logs — zerolog.Logger
```

## License

[MIT](LICENSE) © Nimble Tech
