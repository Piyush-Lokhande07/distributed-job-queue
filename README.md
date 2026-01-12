# Distributed Job Queue (Go + Redis + Docker)

A robust, production-oriented distributed job queue system built with Go and Redis.  
This project focuses on **reliability, backpressure, graceful shutdown, and observability**, mirroring how real backend systems process asynchronous workloads.

---

## 🚀 Features

- **Asynchronous Processing**  
  Concurrent workers pull jobs from a Redis-backed distributed queue.

- **Backpressure & Rate Control**  
  Rejects incoming jobs when the queue depth exceeds safe limits to protect the system under load.

- **Retry with Exponential Backoff**  
  Failed jobs are retried automatically with increasing delay.

- **Failure Injection**  
  Supports deterministic and random failure injection to validate retry and resilience behavior.

- **Graceful Shutdown**  
  Ensures in-flight jobs complete before workers exit on `SIGINT` / `SIGTERM`.

- **Observability**  
  Structured JSON logs and a `/metrics` endpoint for real-time system health.

- **Persistence (Redis Hashes)**  
  Final job outcomes are persisted for auditing and inspection.

- **Dockerized**  
  One-command setup using Docker Compose.

---

## 🛠 Tech Stack

- **Language:** Go  
- **Queue & Persistence:** Redis (Lists + Hashes)  
- **Infrastructure:** Docker, Docker Compose  
- **Logging:** Structured JSON logging  
- **Concurrency:** Goroutines, Channels, WaitGroups  

---

## 🚦 Getting Started

### Prerequisites

- Docker & Docker Compose
- (Optional) Go installed for local runs

---

### Quick Start (Docker)

1. Clone the repository:
```
git clone https://github.com/Piyush-Lokhande07/distributed-job-queue.git
cd distributed-job-queue
```


2. Start the entire stack:

```bash
docker compose up --build
```

The API will be available at:

```
http://localhost:8080
```

---

## 📡 API Documentation

### Enqueue a Job

Submit a job to the queue.

* **Endpoint:** `/jobs`
* **Method:** `POST`

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"id":101}'
```

---

### Get System Metrics

View real-time system statistics.

* **Endpoint:** `/metrics`
* **Method:** `GET`

```bash
curl http://localhost:8080/metrics
```

**Example response:**

```json
{
  "processed": 120,
  "failed": 7,
  "retried": 15,
  "in_progress": 3,
  "queue_depth": 42
}
```

---

### Check Job Status

Query the final status of a job.

* **Endpoint:** `/status?id=<job_id>`
* **Method:** `GET`

---

## 🛡 Graceful Shutdown

On receiving `SIGINT` or `SIGTERM`, the system:

1. Stops accepting new jobs
2. Signals workers to stop fetching tasks
3. Waits for in-flight jobs to finish
4. Persists final job results
5. Closes Redis connections cleanly

This ensures no partial or lost jobs.


