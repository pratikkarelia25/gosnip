# GoSnip

A simple URL shortener with a Go backend and React frontend.

> This is a practice project I'm using to get better at Golang — expect it to evolve as I learn.

## Features

- Shorten a long URL to a short link
- Optional custom short codes, auto-generated when not provided
- Duplicate detection — shortening the same URL returns the existing code
- Optional expiration — links stop redirecting after a set time
- Clean React + shadcn/ui frontend

## Tech Stack

**Backend**
- [Go](https://go.dev/) + [chi](https://github.com/go-chi/chi) — HTTP router
- [SQLite](https://www.sqlite.org/) via [go-sqlite3](https://github.com/mattn/go-sqlite3)

**Frontend**
- [React](https://react.dev/) + [Vite](https://vite.dev/) + TypeScript
- [shadcn/ui](https://ui.shadcn.com/) + Tailwind CSS

## Project Structure

```
backend/
  cmd/gosnip/        # entry point
  internal/api/      # HTTP handlers
  internal/store/    # database layer
  internal/shortcode/ # code generation
  internal/validate/ # URL validation
  database/          # SQLite db (gitignored)
frontend/
  src/
    App.tsx          # main UI
    api.ts           # backend API calls
    components/ui/   # shadcn components
```

## Getting Started

### Backend

```bash
cd backend
go run ./cmd/gosnip
```

Runs on `http://localhost:8080`.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Runs on `http://localhost:5173`.

## API

### POST /shorten

```json
{
  "long_url": "https://example.com",
  "short_code": "mycode",
  "expires_in_seconds": 3600
}
```

`short_code` and `expires_in_seconds` are optional.

### GET /{code}

Redirects to the original URL (`307`), or `404` if not found / expired.

## License

See [LICENSE](LICENSE).
