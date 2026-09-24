# TaskForge

TaskForge is a Go backend job processing system that accepts jobs through a REST API, stores them in PostgreSQL, and processes them asynchronously using RabbitMQ workers.

The project demonstrates backend development concepts including REST APIs, asynchronous job processing, database persistence, message queues, retries, error handling, graceful shutdown, Docker, and automated testing.

## Features

- Create and retrieve jobs through a REST API
- Store job data and results in PostgreSQL
- Process jobs asynchronously with RabbitMQ
- Support job execution with worker processes
- Retry failed jobs up to 3 attempts
- Track job status: `pending`, `running`, `completed`, and `failed`
- Store execution results and error messages
- Graceful shutdown for the API and worker
- Docker Compose setup for the full application stack
- Automated tests with Go's testing framework

## Architecture

TaskForge uses a REST API, PostgreSQL database, RabbitMQ message broker, and a background worker.

```text
Client
  │
  ▼
Go API
  │
  ├──────────────► PostgreSQL
  │
  ▼
RabbitMQ
  │
  ▼
Go Worker
  │
  └──────────────► PostgreSQL
```

### Request flow

1. A client sends a job to the REST API.
2. The API stores the job in PostgreSQL with a `pending` status.
3. The API publishes the job ID to RabbitMQ.
4. The worker consumes the job from RabbitMQ.
5. The worker retrieves the job from PostgreSQL.
6. The worker executes the job.
7. The worker stores the result and updates the job status.
8. Failed jobs are retried up to 3 attempts before being marked as failed.

## Tech Stack

- **Go** — REST API and background worker
- **PostgreSQL** — persistent job storage
- **RabbitMQ** — asynchronous job queue
- **Docker** — containerized application services
- **Docker Compose** — local development and service orchestration
- **Go testing package** — automated tests

## Project Structure

```text
TaskForge/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── jobs.go
│   │   └── jobs_test.go
│   ├── database/
│   │   └── database.go
│   ├── jobs/
│   │   ├── executor.go
│   │   ├── executor_test.go
│   │   ├── jobs.go
│   │   ├── repository.go
│   │   └── repository_test.go
│   └── queue/
│       └── rabbitmq.go
├── migrations/
│   └── 001_create_jobs_table.sql
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

### Key directories

- `cmd/api` — starts the HTTP API server
- `cmd/worker` — starts the background job worker
- `internal/api` — HTTP handlers and API tests
- `internal/database` — PostgreSQL connection setup
- `internal/jobs` — job models, execution logic, and repository
- `internal/queue` — RabbitMQ integration
- `migrations` — database schema migrations

## Getting Started

### Prerequisites

Make sure you have the following installed:

- Go 1.26+
- Docker
- Docker Compose

### Run with Docker Compose

Clone the repository and enter the project directory:

```bash
git clone https://github.com/RivellionCS/TaskForge.git
cd TaskForge
```

Copy the example environment file:

```bash
cp .env.example .env
```

Start the application:

```bash
docker compose up -d --build
```

Check that all services are running:

```bash
docker compose ps
```

The application uses the following services:

- API — `localhost:8080`
- PostgreSQL — `localhost:5432`
- RabbitMQ — `localhost:5672`
- RabbitMQ Management UI — `localhost:15672`

## API Usage

### Create a Job

Create a `sleep` job using the REST API:

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"type":"sleep","payload":{"seconds":5}}'
```

The API returns a job ID and its initial status:

```json
{
  "id": "JOB_ID",
  "status": "pending"
}
```

The job is then published to RabbitMQ and processed asynchronously by the worker.

### Get a Job

Use the returned job ID to retrieve the job:

```bash
curl http://localhost:8080/jobs/JOB_ID
```

A completed job returns information such as:

```json
{
  "id": "JOB_ID",
  "type": "sleep",
  "status": "completed",
  "payload": {
    "seconds": 5
  },
  "result": {
    "message": "slept for 5 seconds"
  },
  "attempts": 0,
  "created_at": "2026-09-24T07:27:47.438709Z",
  "started_at": "2026-09-24T07:27:47.444674Z",
  "completed_at": "2026-09-24T07:27:52.476208Z"
}
```

## Job Processing

TaskForge currently supports the `sleep` job type.

A sleep job accepts a payload containing the number of seconds to wait:

```json
{
  "type": "sleep",
  "payload": {
    "seconds": 5
  }
}
```

The worker executes the job asynchronously and stores the result in PostgreSQL.

### Retry Behavior

If a job fails during execution, the worker increments its attempt count and requeues the job.

Jobs are retried up to 3 attempts. After the third failed attempt, the job is marked as `failed` and the error is stored as the job result.

## Testing

Run the complete test suite with:

```bash
go test ./...
```

The project includes tests for:

- API handlers
- Job execution
- Job repository operations
- Successful job processing
- Job failure and retry behavior

## Graceful Shutdown

Both the API and worker support graceful shutdown.

When the processes receive `SIGINT` or `SIGTERM`, they shut down cleanly and close their database and RabbitMQ connections.

## Environment Configuration

TaskForge uses environment variables for its database and RabbitMQ connections.

The `.env.example` file documents the required configuration:

```env
# Local development
DATABASE_URL=postgres://taskforge:taskforge@localhost:5432/taskforge
RABBITMQ_URL=amqp://taskforge:taskforge@localhost:5672/

# Docker Compose
DATABASE_URL_DOCKER=postgres://taskforge:taskforge@postgres:5432/taskforge
RABBITMQ_URL_DOCKER=amqp://taskforge:taskforge@rabbitmq:5672/
```

The `.env` file is excluded from version control through `.gitignore`.

## License

This project is intended as a backend development portfolio project.
