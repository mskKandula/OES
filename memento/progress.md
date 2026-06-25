# OES — Progress

## Current Status
**All 7 Docker images built and pushed to DockerHub (`mskkandula/oes`). Kubernetes deployment YAMLs updated to reference DockerHub images.**

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

All tags live at: https://hub.docker.com/r/mskkandula/oes/tags

All `k8s/` deployment YAMLs updated — `image:` fields reference `mskkandula/oes:<tag>` with `imagePullPolicy: Always`.

## What Works (Implemented & Wired Up)

### Infrastructure
- [x] Full 9-container Docker Compose stack boots with health-checked dependency ordering
- [x] Multi-stage Docker builds for all Go, Node, and Python services
- [x] Shared `./media` bind-mount volume across `server`, `fileserver`, `mqserver`
- [x] Nginx reverse proxy routing (API, CDN, WebSocket, SPA fallback)
- [x] Self-signed TLS cert generated at build time in client container

### Authentication & Authorisation
- [x] Examiner sign-up (`POST /o/signUp`) — bcrypt hash, UUID `clientId`, DB insert + UserRole
- [x] Login (`POST /o/login`) — UNION query across Examiners + Students, bcrypt verify, JWT in HttpOnly cookie
- [x] Role-based middleware (`Auth("Common")`, `Auth("Examiner")`, `Auth("Student")`)
- [x] Logout — cookie cleared (`MaxAge: -1`)
- [x] Dynamic role-based menu (`GET /r/getRoutes`) — Redis-cached for 24 h

### Examiner Features
- [x] Bulk student registration from `.xlsx` via `POST /r/multipleStudentsRegistration` — concurrent bcrypt + DB insert per student
- [x] Welcome email published to RabbitMQ `email` queue → consumed by `MQServer` → sent via Gmail SMTP
- [x] Question paper upload (`POST /r/uploadQuestionFile`) — `.txt` parsed line-by-line, exam created in DB, questions cached in-memory
- [x] Video upload (`POST /r/uploadVideoContent`) — MP4 saved, DB record inserted, encode job published to `encode` queue
- [x] Video HLS encoding — `MQServer` consumes `encode` queue, runs `ffmpeg` for 360p/480p/720p variants + thumbnail, deletes original
- [x] Get all students (`GET /r/getStudents`)
- [x] Download students as Excel (`GET /r/downloadStudents`)
- [x] AI question generation stub (`POST /r/questionGeneration`) — gRPC call to Python `isupport` service

### Student Features
- [x] Get exam questions from in-memory cache (`GET /r/getQuestions?examId=X`)
- [x] Upload exam proof zip (`POST /r/uploadExamProof`) — saved to `media/examProofs/`, queued for async unzip
- [x] Async unzip worker pool — extracts images, bulk-inserts paths into `StudentExamProofs`, deletes zip
- [x] In-browser answer evaluation — Go WASM `countWords` function counts keyword matches

### Real-Time (WebSocket)
- [x] WebSocket endpoint (`GET /r/ws`) with JWT cookie auth
- [x] 32-shard pool with netpoll edge-triggered reads
- [x] Redis pub/sub `general` channel for broadcast fan-out
- [x] Per-shard dedicated write goroutines
- [x] Heartbeat ping/pong with dead client cleanup
- [x] Message type routing: notifications (1), chat (2), whiteboard (3), targeted WebRTC signal (4), presence (5)
- [x] Vuex store mutations for all WS message types on frontend

### Video
- [x] Redis-cached video list (`GET /r/getVideos`, 15 min TTL, invalidated on new upload)
- [x] FileServer serving media directory at port `8887` → `/cdn/*` via Nginx
- [x] Live streaming server: RTMP ingest on port `1935`, HLS output to `/mnt/hls/`, 5 adaptive bitrate variants

### Frontend (Vue.js)
- [x] Vue Router with lazy-loaded views (history mode)
- [x] Examiner views: Dashboard, MultipleStudentsRegistration, StudentsList, UploadQuestions, UploadVideo, BroadcastVideo, QuestionGen
- [x] Student views: StudentDashboard, Exam (with voice input + on-screen keyboard + webcam capture), VideoContent, VideoPlay, VideoCaptioning
- [x] Common views: WhiteBoard (WS + Canvas), WordCounter (WASM result display)
- [x] Charts: Bar, Line (real-time dashboard)
- [x] WebRTC video call component (`BroadcastVideo.vue`, `VideoCall.vue`, `VideoMenu.vue`)

## What Is NOT Yet Done / Known Gaps

### Critical Gaps
- [ ] **Exam answers not persisted** — `StudentExamMarks` table exists in schema but no API endpoint writes marks to it
- [ ] **In-memory question cache** — questions lost on server restart; not stored in MySQL or Redis
- [ ] **JWT refresh tokens** — no refresh mechanism; session expires in 15 minutes with no recovery
- [ ] **JWT secret hardcoded** — `7yt65U745TR57lo9h%$fre#$TR43EW` must be moved to env var

### AI / NLP Features (Stub Only)
- [ ] `isupport` Python gRPC service is an **echo stub** — returns `"Hello I am up and running received "..." message from you"` — no real NLP
- [ ] Document Similarity (NLP)
- [ ] Text Summariser (NLP)
- [ ] Student Mood Analyser (CNN)

### Security / Production Hardening
- [ ] HTTPS / valid TLS certificate (currently self-signed)
- [ ] gRPC uses `WithInsecure()` — no mTLS
- [ ] SMTP credentials via env vars (wired) but no secret management
- [ ] Cookie `Secure: false` (must be `true` with real HTTPS)
- [ ] No rate limiting on auth endpoints
- [ ] No input sanitisation beyond file type/size checks

### Operational
- [ ] No centralised logging / log aggregation
- [ ] No metrics / observability (pprof endpoint exists but no Prometheus)
- [ ] No graceful shutdown handling in `server` or `mqserver`
- [ ] Live stream HLS output not confirmed wired to frontend `VideoPlay` view

### WASM
- [ ] `main.wasm` pre-built binary is committed — no automated rebuild step in Docker build pipeline
- [ ] Word-counter logic uses a hardcoded test dictionary (`["Mobile", "computer", "flower"]`) — not configurable per exam

## Known Issues
- `MQServer` uses `auto-ack` for `encode` queue but `manual-ack` for `email` queue — encode jobs are lost if ffmpeg crashes mid-job.
- `publishMu` mutex in `studentService` serialises all RabbitMQ publish calls during bulk student registration — could be a bottleneck for very large Excel uploads.
- The `ResultPaths` channel in `studentHandler` is a package-level `var` (buffered 200) — it is closed in `main.go`'s defer, coupling handler and main tightly.
- Redis `ReadVideos` cache uses `clientId` directly as the cache key — key collisions are impossible but it's a non-descriptive key format compared to `routes:` prefixed keys.
