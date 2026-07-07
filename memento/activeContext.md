# OES — Active Context

## Current State
All 8 Docker images built, pushed to DockerHub (`mskkandula/oes`), and deployed on a 4-node kind (Kubernetes in Docker) cluster. The CodeExecutor API endpoint has been tested and confirmed working for Python and Node.js.

## What Was Just Done (2026-07-07)

### CodeExecutor Service — New since last session
A new `CodeExecutor/` microservice was added to the project. It:
- Consumes code execution jobs from RabbitMQ queues (`code.execute.python`, `code.execute.go`, `code.execute.nodejs`)
- Picks a warm runner pod from a per-language pool
- Executes code via `kubectl exec` (Kubernetes pods/exec subresource)
- Deletes the pod after execution (Deployment controller replaces it automatically)
- Publishes results to Redis: `PUBLISH result:<submissionId>` + `SET result:<submissionId>` + `PUBLISH general` (WebSocket type-6)

### Build & Push
All 8 images pushed to `mskkandula/oes` on DockerHub:
- `db`, `server`, `client`, `fileserver`, `mqserver`, `isupport`, `liveserver` (existing 7)
- `code-executor` (new — fixed build error: `model.go` was in wrong package directory `internal/k8s/`, moved to `internal/model/`)

### k8s Manifests Fixed
1. [`k8s/kustomization.yaml`](k8s/kustomization.yaml) — Added 5 CodeExecutor resources (copied from `CodeExecutor/k8s/` to `k8s/` to satisfy kustomize security boundary)
2. [`k8s/client-deployment.yaml`](k8s/client-deployment.yaml) — Fixed OOMKill: reduced nginx worker limits; changed topology spread to `ScheduleAnyway`
3. [`k8s/codeexecutor-deployment.yaml`](k8s/codeexecutor-deployment.yaml) — Changed `imagePullPolicy: IfNotPresent` (was `Always`, causing stuck pod on kind)

### Deployment
- `kubectl apply -k k8s/` — all resources created in `oes` namespace
- `kubectl port-forward svc/oes-server-nodeport 30900:9000 -n oes` — used for API access

### API Testing Results
```
POST http://localhost:30900/r/executeCode
```

| Language | Code | Result |
|---|---|---|
| python | `print('Hello from OES CodeExecutor!')` | ✅ `completed`, stdout correct, durationMs=554 |
| nodejs | `console.log('Hello from Node.js runner!')` | ✅ `completed`, stdout correct, durationMs=686 |
| python | `x = 1/0` | ✅ `failed`, ZeroDivisionError in stderr, exitCode=1 |
| go | `package main...` | ❌ `failed` — `go.mod file not found` (needs fix) |

All Python/Node tests returned `pending: false` — confirming **Redis pub/sub fast path** (200 OK within 5s) is working end-to-end.

## Architecture Understanding Summary

### Service Dependency Startup Order (Kubernetes)
```
oes-db (MySQL) → healthy
oes-cache (Redis) → healthy
oes-messageq (RabbitMQ) → healthy
                              ↓
                        oes-server (Go API) → starts
                        oes-mqserver (Go+ffmpeg) → starts after messageq
                        oes-code-executor → starts after messageq + redis
oes-vectordb (ChromaDB) → starts (independent)
oes-ollama (LLM) → starts (independent, slow — model download)
                              ↓
                        oes-isupport (Python gRPC RAG) → starts after ollama + vectordb
oes-client (Nginx SPA) → starts (independent, fast)
oes-fileserver (Go) → starts (independent)
oes-liveserver (nginx-rtmp) → starts (independent)
Runner pods (python/go/nodejs) → start (independent, warm pool)
```

### Key Inter-Service Communication Paths (Kubernetes DNS names)

| Path | Protocol | Notes |
|---|---|---|
| Browser → port-forward → `oes-server:9000` | HTTP | API calls |
| `oes-server` → `oes-messageq:5672` | AMQP | Publish code jobs |
| `oes-code-executor` → `oes-messageq:5672` | AMQP | Consume code jobs |
| `oes-code-executor` → Kubernetes API | HTTPS | Watch/exec/delete runner pods |
| `oes-code-executor` → `oes-cache:6379` | Redis | Publish results |
| `oes-server` → `oes-cache:6379` | Redis | Subscribe to results (5s timeout) |
| `oes-server` → `oes-db:3306` | MySQL | All CRUD |

## Active Decisions & Considerations

### CodeExecutor Model Package Fix
The `CodeExecutor/internal/k8s/model.go` declared `package model` but lived in the `k8s` directory alongside files with `package k8s`. This caused the Go compiler to reject the build. Fix: created `CodeExecutor/internal/model/model.go` (proper location) and deleted the misplaced file.

### imagePullPolicy on kind
`imagePullPolicy: Always` forces Kubernetes to contact the registry even when the image is already in the node's containerd cache. On kind this can cause pods to be stuck `PodInitializing` for extended periods if DockerHub is slow. Fixed for `codeexecutor-deployment.yaml` by using `IfNotPresent`. Other deployments should also be updated.

### Go Runner Module Issue
`go run /dev/stdin` without a module context fails with `go.mod file not found`. Fix options:
1. Add `GONOSUMDB=*` + `GOFLAGS=-mod=mod` env vars to `runner-go-deployment.yaml`  
2. Or provide a minimal `go.mod` in the runner container's working directory via ConfigMap

### Two kustomization files
- [`k8s/kustomization.yaml`](k8s/kustomization.yaml) — the **canonical** deploy root (includes all platform + CodeExecutor resources)
- [`CodeExecutor/k8s/kustomization.yaml`](CodeExecutor/k8s/kustomization.yaml) — standalone reference (references files in `k8s/` parent, only usable from that directory)

Use `kubectl apply -k k8s/` always.

## Next Steps (If Development Continues)
1. **Fix Go runner** — add `GONOSUMDB=*` + `GOFLAGS=-mod=mod` to [`k8s/runner-go-deployment.yaml`](k8s/runner-go-deployment.yaml)
2. **Change remaining `imagePullPolicy: Always`** to `IfNotPresent` in other k8s deployments for faster kind startup
3. **Persist questions** — store in MySQL/Redis instead of in-memory cache
4. **JWT refresh token** — implement refresh token rotation
5. **Externalise JWT secret** — move `7yt65U745TR57lo9h%$fre#$TR43EW` to env var
6. **Exam answer persistence** — write marks to `StudentExamMarks` table
7. **Implement NLP in `isupport`** — RAG pipeline with langchain+ollama+chromadb (infrastructure is deployed, pipeline code needs completion)
