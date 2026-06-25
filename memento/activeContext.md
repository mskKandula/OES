# OES — Active Context

## Current State
The project has been fully read and understood for the first time (initial analysis session). No code changes have been made. The memento documentation set is being created as the baseline for all future sessions.

## What Was Just Done
- Read all Dockerfiles, `docker-compose.yml`, and `nginx.conf` files for both Client and LiveStreamingServer.
- Read all Go source files: `main.go`, `api/api.go`, all handlers, services, repositories, data-source helpers, WebSocket pool/client, running-process goroutine, config, models.
- Read `MQServer/main.go`, `FileServer/files.go`, `IntelligenceSupport/questgen/server.py`, the `.proto` file.
- Read the MySQL schema (`MySQL/oes.sql`), Vue router, Vuex store, `ws.js`, `Exam.vue`, WASM Go source.
- Created all six core memento files.

## Architecture Understanding Summary

### Service Dependency Startup Order
```
messageq (RabbitMQ) → healthy
db (MySQL 8) → healthy
cache (Redis) → healthy
isupport (Python gRPC) → started
                              ↓
                        server (Go API) → starts
                              ↓
                        client (Nginx SPA) → starts
                        fileserver (Go) → starts (independent)
                        mqserver (Go+ffmpeg) → starts after messageq
                        liveserver (nginx-rtmp) → starts (independent)
```

### Key Inter-Service Communication Paths

| Path | Protocol | Parties | Purpose |
|---|---|---|---|
| Browser → Nginx (port 8080) | HTTP/WS | Browser ↔ `oes_client` | All browser traffic entry point |
| Nginx → Go API | HTTP/WS | `oes_client` → `oes_server:9000` | API calls + WebSocket upgrade |
| Nginx → FileServer | HTTP | `oes_client` → `oes_fileserver:8887` | Serve media files (`/cdn/*`) |
| Go API → MySQL | TCP/MySQL protocol | `oes_server` → `oes_db:3306` | All CRUD operations |
| Go API → Redis | TCP/RESP | `oes_server` → `cache:6379` | Route cache, video cache, WS pub/sub |
| Go API → RabbitMQ | AMQP 0-9-1 | `oes_server` → `messageq:5672` | Publish `encode` and `email` jobs |
| Go API → Python gRPC | HTTP/2 gRPC | `oes_server` → `isupport:50051` | AI question generation |
| MQServer → RabbitMQ | AMQP 0-9-1 | `oes_mqserver` → `messageq:5672` | Consume `encode` and `email` jobs |
| MQServer → shared volume | filesystem | `oes_mqserver` ↔ `./media` | Read MP4, write HLS, delete original |
| Browser → LiveServer | RTMP | Browser/OBS → `oes_liveserver:1935` | Examiner broadcasts live stream |
| Redis Pub/Sub | RESP | `oes_server` ↔ `cache` | Fan-out WebSocket messages across potential server instances |

## Active Decisions & Considerations

### In-Memory Question Cache
Questions uploaded via [`QuestionsUpload`](Server/api/handler/userHandler.go:99) are stored in a [`QuestionCache`](Server/api/handler/handler.go:10) (a `sync.RWMutex`-protected `map[int64][]string`). This is intentionally not persisted to Redis or MySQL — it means:
- Questions are lost if the server restarts.
- A future improvement would be to persist questions in MySQL or Redis.

### Video Encoding Pipeline
Video upload flow: HTTP upload → save MP4 → DB record → publish to `encode` queue → return 200 to browser. The actual HLS encoding happens async in `oes_mqserver`. The client receives success immediately but the HLS stream may not be ready for a few seconds/minutes.

### WebSocket & Redis Pub/Sub Scaling
The Redis pub/sub `general` channel is already in place to support **horizontal scaling** of the WebSocket server. Adding a second `oes_server` instance would work without code changes because all broadcast messages go through Redis. Currently only one instance runs.

### WASM Binary
The pre-built `main.wasm` in `Client/oes/public/` is served by Nginx with `application/wasm` MIME type (handled via `mime.types`). The `wasm_exec.js` Go runtime shim is included in `src/`. The WASM exports a single `countWords` function used in exam submission.

### Multi-Tenant Design
`clientId` (UUID) is the tenant isolation key — it is set on Examiner sign-up and propagated to all students registered by that examiner. All API calls that require data isolation use `clientId` extracted from the JWT claim.

## Next Steps (If Development Continues)
1. **Persist questions** — store question lists in MySQL/Redis instead of in-memory cache.
2. **JWT refresh token** — implement refresh token rotation so users aren't kicked every 15 minutes.
3. **Externalise JWT secret** — move `7yt65U745TR57lo9h%$fre#$TR43EW` to env var / Docker secret.
4. **Implement NLP in `isupport`** — replace the echo stub with an actual sentence-BERT or T5 model for question generation.
5. **Add HTTPS** — enable a proper TLS cert (Let's Encrypt) rather than self-signed.
6. **Horizontal scaling test** — spin up 2× `oes_server` instances behind Nginx load balancer to validate the Redis pub/sub WebSocket fan-out.
7. **Exam answer persistence** — currently `StudentExamMarks` table exists but no API writes to it.
8. **Live stream viewing** — the `VideoPlay` view exists in the router but the HLS player wiring to the live server's output needs verification.
