# GoSnip

A simple URL shortener written in Go.

> This is a practice project I'm using to get better at Golang — expect it to evolve as I learn.

## Features

- Shorten a long URL to a short code
- Optional custom short codes, with auto-generated ones when not provided
- Duplicate detection — shortening the same long URL again returns the existing short code
- Optional expiration (`expires_in_seconds`) after which a short code stops redirecting
- Redirect from a short code to its original URL

## Tech Stack

- [Go](https://go.dev/)
- [chi](https://github.com/go-chi/chi) — HTTP router
- [SQLite](https://www.sqlite.org/) (via [go-sqlite3](https://github.com/mattn/go-sqlite3)) — storage

## Project Structure

```
cmd/gosnip/        # application entry point
internal/api/       # HTTP handlers and routing
internal/store/      # database access
internal/shortcode/  # short code generation
internal/validate/    # input validation
database/           # sqlite database file (gitignored)
```

## Getting Started

### Prerequisites

- Go 1.25+

### Run

```bash
go run ./cmd/gosnip
```

The server starts on `http://localhost:8080`.

## API

### Create a short URL

```
POST /shorten
Content-Type: application/json

{
  "long_url": "https://example.com",
  "short_code": "mycode",        // optional, auto-generated if omitted
  "expires_in_seconds": 3600     // optional, never expires if omitted
}
```

Response:

```json
{
  "message": "Short URL created successfully",
  "short_code": "mycode",
  "expires_at": "2026-06-17T16:00:00+05:30"
}
```

### Redirect to the original URL

```
GET /{code}
```

Redirects with `307 Temporary Redirect`, or `404 Not Found` if the code doesn't exist or has expired.

## License

See [LICENSE](LICENSE).
