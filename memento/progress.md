# OES — Progress

## Current Status
**All 8 Docker images built and pushed to DockerHub (`mskkandula/oes`). Kubernetes deployment fully operational on a 4-node kind cluster. CodeExecutor API endpoint tested and confirmed working.**

## Docker Images — Published to DockerHub (`mskkandula/oes`)

| Tag | Service | Source Dockerfile |
|---|---|---|
| `mskkandula/oes:db` | MySQL 8 + SQL init scripts | `MySQL/Dockerfile` |
| `mskkandula/oes:server` | Go API Server | `Server/Dockerfile` |
| `mskkandula/oes:client` | Vue.js + Nginx SPA | `Client/Dockerfile` |
| `mskkandula/oes:fileserver` | Go File Server | `FileServer/Dockerfile` |
| `mskkandula/oes:mqserver` | Go MQ Worker + ffmpeg | `MQServer/Dockerfile` |
| `mskkandula/oes:isupport` | Python gRPC | `IntelligenceSupport/questgen/Dockerfile` |
| `mskkandula/oes:liveserver` | nginx-rtmp Live Streaming | `LiveStreamingServer/Dockerfile` |
| `mskkandula/oes:code-executor` | Go CodeExecutor service | `CodeExecutor/Dockerfile` |

All tags live at: https://hub.docker.com/r/mskkandula/oes/tags

## Kubernetes Deployment — Active on 4-node kind cluster

All `k8s/` deployment YAMLs applied via `kubectl apply -k k8s/`.

### Pods Running (verified 2026-07-07)
| Pod | Status |
|---|---|
| `oes-cache-0` (Redis) | ✅ Running |
| `oes-client` (2 replicas, Nginx SPA) | ✅ Running |
| `oes-code-executor` | ✅ Running |
| `oes-db-0` (MySQL) | ✅ Running |
| `oes-fileserver` | ✅ Running |
| `oes-liveserver` (nginx-rtmp) | ✅ Running |
| `oes-messageq-0` (RabbitMQ) | ✅ Running |
| `oes-mqserver` (Go+ffmpeg worker) | ✅ Running |
| `oes-runner-go` (2 replicas) | ✅ Running |
| `oes-runner-nodejs` (2 replicas) | ✅ Running |
| `oes-runner-python` (2 replicas) | ✅ Running |
| `oes-server` (Go API) | ✅ Running |
| `oes-vectordb` (ChromaDB) | ✅ Running |
| `oes-isupport` (Python gRPC RAG) | ⏳ Init (ollama pull in progress) |
| `oes-ollama` (LLM) | ⏳ Init (model pull in progress) |

### Service Access
- Server NodePort: `kubectl port-forward svc/oes-server-nodeport 30900:9000 -n oes`
- Server health: `GET http://localhost:30900/o/status`

## CodeExecutor API — Tested & Working

### Endpoint
`POST /r/executeCode` — **no auth** (intentionally removed for testing)

### Test Results

**Python (success)**:
```json
{"submissionId":"35143f41...","pending":false,"status":"completed","stdout":"Hello from OES CodeExecutor!\n","durationMs":554}
```

**Node.js (success)**:
```json
{"submissionId":"c70ea510...","pending":false,"status":"completed","stdout":"Hello from Node.js runner!\n","durationMs":686}
```

**Python error handling**:
```json
{"submissionId":"f4f2de17...","status":"failed","stderr":"ZeroDivisionError: division by zero\n","exitCode":1,"durationMs":1431}
```

**Go runner**: Fails with `go.mod file not found` — `go run /dev/stdin` requires module context. Fix: add `GONOSUMDB=*` + `GOFLAGS=-mod=mod` env vars to `runner-go-deployment.yaml`.

### Fast Path (200 OK)
Results returned within 5s window synchronously — all Python/Node tests returned `pending:false` confirming the Redis pub/sub fast path is working.

## k8s Manifest Changes Made

1. [`k8s/kustomization.yaml`](k8s/kustomization.yaml) — Added CodeExecutor resources:
   - `codeexecutor-rbac.yaml` (ServiceAccount + Role + RoleBinding)
   - `codeexecutor-deployment.yaml` (code-executor Deployment)
   - `runner-python-deployment.yaml` (Python warm runner pool, 2 replicas)
   - `runner-go-deployment.yaml` (Go warm runner pool, 2 replicas)
   - `runner-nodejs-deployment.yaml` (Node.js warm runner pool, 2 replicas)

2. [`k8s/client-deployment.yaml`](k8s/client-deployment.yaml) — Fixed for local kind cluster:
   - Reduced nginx `worker_processes 1` + `worker_rlimit_nofile 1024` (was OOMKilling at 256Mi)
   - Changed topology spread `whenUnsatisfiable: ScheduleAnyway` (was `DoNotSchedule`)

3. [`k8s/codeexecutor-deployment.yaml`](k8s/codeexecutor-deployment.yaml) — Changed `imagePullPolicy: IfNotPresent` (was `Always` which caused stuck pod on kind)

4. [`CodeExecutor/internal/model/`](CodeExecutor/internal/model/model.go) — Created `internal/model/` package directory, moved `model.go` from `internal/k8s/` (was causing Go build error: two packages in same directory)

## What Works (Implemented & Wired Up)

### Infrastructure
- [x] Full 9-container Docker Compose stack boots with health-checked dependency ordering
- [x] Multi-stage Docker builds for all Go, Node, and Python services
- [x] Shared `./media` bind-mount volume across `server`, `fileserver`, `mqserver`
- [x] Nginx reverse proxy routing (API, CDN, WebSocket, SPA fallback)
- [x] Self-signed TLS cert generated at build time in client container
- [x] All 8 images pushed to DockerHub `mskkandula/oes`
- [x] Kubernetes manifests deployed on 4-node kind cluster

### Authentication & Authorisation
- [x] Examiner sign-up (`POST /o/signUp`) — bcrypt hash, UUID `clientId`, DB insert + UserRole
- [x] Login (`POST /o/login`) — UNION query across Examiners + Students, bcrypt verify, JWT in HttpOnly cookie
- [x] Role-based middleware (`Auth("Common")`, `Auth("Examiner")`, `Auth("Student")`)
- [x] Logout — cookie cleared (`MaxAge: -1`)
- [x] Dynamic role-based menu (`GET /r/getRoutes`) — Redis-cached for 24 h

### Code Execution (NEW)
- [x] `POST /r/executeCode` endpoint — no auth for testing
- [x] Python code execution via warm runner pods (kubectl exec)
- [x] Node.js code execution via warm runner pods (kubectl exec)
- [x] Result delivered via Redis pub/sub fast path (200 OK within 5s)
- [x] Error/exception handling — stderr + exitCode returned correctly
- [x] code-executor RBAC (ServiceAccount + Role + RoleBinding for pod watch/exec/delete)
- [x] Warm runner pool (2 pods per language, 6 total)
- [ ] Go code execution — fails due to missing go.mod (needs GONOSUMDB fix)

### Examiner Features
- [x] Bulk student registration from `.xlsx` via `POST /r/multipleStudentsRegistration`
- [x] Welcome email published to RabbitMQ `email` queue → consumed by `MQServer`
- [x] Question paper upload (`POST /r/uploadQuestionFile`) — in-memory cache
- [x] Video upload (`POST /r/uploadVideoContent`) — MP4 saved, encode job published
- [x] Video HLS encoding — `MQServer` consumes `encode` queue, runs `ffmpeg`
- [x] Get all students, download students as Excel
- [x] AI question generation stub (`POST /r/questionGeneration`) — gRPC call to `isupport`

### Student Features
- [x] Get exam questions from in-memory cache
- [x] Upload exam proof zip — async unzip worker pool
- [x] In-browser answer evaluation — Go WASM `countWords`

### Real-Time (WebSocket)
- [x] WebSocket endpoint with JWT cookie auth
- [x] 32-shard pool with netpoll edge-triggered reads
- [x] Redis pub/sub `general` channel for broadcast fan-out
- [x] Message type routing (types 1-5 + new type 6 for code execution results)

### Video
- [x] Redis-cached video list
- [x] FileServer serving media directory
- [x] Live streaming server: RTMP ingest on port `1935`

## What Is NOT Yet Done / Known Gaps

### Critical Gaps
- [ ] **Go runner needs fix** — `GONOSUMDB=*` + `GOFLAGS=-mod=mod` needed in runner-go-deployment.yaml
- [ ] **Exam answers not persisted** — `StudentExamMarks` table exists but no API endpoint writes marks
- [ ] **In-memory question cache** — questions lost on server restart
- [ ] **JWT refresh tokens** — session expires in 15 minutes with no recovery
- [ ] **JWT secret hardcoded** — `7yt65U745TR57lo9h%$fre#$TR43EW`

### AI / NLP Features (Stub Only)
- [ ] `isupport` Python gRPC service — RAG pipeline (langchain + ollama + chromadb) — initializing (slow)
- [ ] Document Similarity, Text Summariser, Student Mood Analyser

### Security / Production Hardening
- [ ] HTTPS / valid TLS certificate
- [ ] Cookie `Secure: false`
- [ ] No rate limiting on auth endpoints

### WASM
- [ ] `main.wasm` pre-built binary committed — no automated rebuild
- [ ] Word-counter uses hardcoded test dictionary

## Known Issues
- `imagePullPolicy: Always` on most k8s deployments causes slow starts on kind (images already present but registry is contacted). Fixed for code-executor; others may benefit from `IfNotPresent`.
- Go runner pods fail with `go.mod file not found` — needs `GONOSUMDB=*` and `GOFLAGS=-mod=mod` env vars.
- `MQServer` uses `auto-ack` for `encode` queue — encode jobs lost if ffmpeg crashes.
- `publishMu` mutex serialises RabbitMQ publish calls during bulk student registration.
- Redis `ReadVideos` cache uses `clientId` directly as key (non-descriptive format).
