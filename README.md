# Code Colosseum

Code Colosseum is a competitive programming platform built around real-time 1v1 matches.

Players get matched based on rating, solve programming problems head-to-head, and gain or lose ELO based on match results.

## Features

- Real-time 1v1 matches
- Rating-based matchmaking
- ELO rating system
- Live match updates over WebSockets
- Code execution in Docker
- Multiple programming languages
- Problem and test case management
- Submissions and match history
- Leaderboard
- LeetCode problem importing

## Tech stack

**Backend**

- Go
- Gin
- PostgreSQL
- sqlc
- Goose
- WebSockets
- Docker

**Frontend**

- Next.js
- TypeScript
- Tailwind CSS
- shadcn/ui
- Lucide

**Seed tooling**

- TypeScript
- Bun
- LeetCode API
- AI-generated problems

## Getting started

### Requirements

You'll need:

- Go
- Node.js
- Bun
- Docker
- [Goose](https://github.com/pressly/goose)
- [sqlc](https://sqlc.dev/)

## Code execution

Submissions run inside Docker containers.

Supported languages:

- C++
- Go
- Java
- JavaScript
- TypeScript
- Python

Language Dockerfiles are located at:

```text
backend/internal/language/dockerfiles/
```

Judge implementation lives in:

```text
backend/internal/judge/
backend/internal/language/
```

Rebuild judge images after changing their Dockerfiles:

```bash
cd backend
go run scripts/build_images.go
```

## Project structure

```text
backend/
├── cmd/api/              # API entrypoint
├── internal/
│   ├── auth/             # Authentication
│   ├── db/               # Database, migrations, sqlc
│   ├── elo/              # ELO calculations
│   ├── judge/            # Code execution
│   ├── language/         # Language runtimes
│   ├── leaderboard/      # Leaderboard
│   ├── match/            # Matchmaking and matches
│   ├── problems/         # Problems
│   ├── submissions/      # Submissions
│   ├── tags/             # Problem tags
│   ├── testcases/        # Test cases
│   └── ws/               # WebSocket handling
│
frontend/
└── src/
    ├── app/              # Next.js routes
    ├── components/       # UI components
    ├── lib/api/          # API clients
    └── providers/        # React providers

seed/
├── src/                  # Problem importing/generation
└── data/                 # Seed data
```

## Development

Typical local setup:

```bash
# Start database
docker compose up -d

# Run migrations
goose up

# Generate database code
sqlc generate -f backend/internal/db/sqlc.yaml

# Build judge images and start backend
cd backend
go run scripts/build_images.go
go run ./cmd/api
```

Then, in another terminal:

```bash
cd frontend
bun dev
```

## Seed

The `seed` directory contains tooling for importing and generating programming problems.
