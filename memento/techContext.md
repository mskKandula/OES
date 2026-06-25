# OES — Tech Context

## Technology Stack

### Frontend
| Technology | Version | Purpose |
|---|---|---|
| Vue.js 2 | `^2.6.11` | SPA framework |
| Vue Router | `^3.2.0` | Client-side routing (history mode) |
| Vuex | `^3.4.0` | Global state (WebSocket conn, notifications, chat, whiteboard, online users) |
| Axios (`vue-axios`) | via `vue-axios` | HTTP client for API calls |
| Bootstrap | via SCSS | UI layout and components |
| Material Design Icons | custom font plugin | Icon set |
| JSZip | — | Zip webcam images in browser for exam proof upload |
| Chart.js (`vue-chartjs`) | — | Real-time dashboard charts (Bar, Line) |
| WebAssembly (Go WASM) | — | Answer word-count evaluation running natively in browser |
| WebRTC | browser native | Peer-to-peer video/audio (exam proctoring, classroom, broadcast) |
| WebSockets | browser native | Real-time notifications, chat, whiteboard |
| WebkitSpeechRecognition | browser native | Voice-to-text input during exam |
| SpeechSynthesisUtterance | browser native | Text-to-speech |
| Canvas API | browser native | Whiteboard drawing, video caption overlay, image capture |
| ImageCapture API | browser native | Grab frames from webcam stream |
| HLS.js or native | browser native | Adaptive bitrate video playback |

### Backend (Main API Server)
| Technology | Version | Purpose |
|---|---|---|
| Go | `1.26` (alpine) | Core language |
| Gin | `github.com/gin-gonic/gin` | HTTP router/framework |
| gobwas/ws | `github.com/gobwas/ws` | Low-level WebSocket library (custom framing) |
| mailru/easygo netpoll | `github.com/mailru/easygo/netpoll` | Edge-triggered epoll for high-concurrency WS reads |
| golang-jwt/jwt | `v5` | JWT generation and validation |
| go-redis/redis | `v8` | Redis client (caching + pub/sub) |
| go-sql-driver/mysql | — | MySQL driver |
| rabbitmq/amqp091-go | — | RabbitMQ AMQP 0-9-1 client |
| google.golang.org/grpc | — | gRPC client to `isupport` |
| bcrypt | `golang.org/x/crypto` | Password hashing |
| errgroup | `golang.org/x/sync` | Concurrent operations with error collection |
| google/uuid | — | UUID generation for `clientId` |
| gorilla/schema | — | URL query-param decoding for WebSocket handshake |
| tealeg/xlsx | `v3` | Excel file read/write |
| tidwall/gjson + sjson | — | Fast JSON get/set for Excel→JSON conversion |
| gin-contrib/pprof | — | Go profiling endpoint |

### MQ Worker Server
| Technology | Version | Purpose |
|---|---|---|
| Go | `1.26` (alpine) | Core language |
| amqp091-go | — | RabbitMQ consumer |
| ffmpeg | Alpine package | Video transcoding (HLS encode), thumbnail extraction |
| gomail.v2 | — | Gmail SMTP email dispatch |
| html/template | stdlib | Welcome email HTML rendering |

### Intelligence Support (gRPC)
| Technology | Version | Purpose |
|---|---|---|
| Python | `3.14` (alpine) | Runtime |
| grpcio | `requirements.txt` | gRPC server |
| grpcio-tools | `requirements.txt` | Protobuf codegen (used at dev time) |

### Infrastructure / Persistence
| Technology | Version | Purpose |
|---|---|---|
| MySQL | `8.0` | Primary relational database |
| Redis | `8.0-alpine` | Caching + WebSocket pub/sub bus |
| RabbitMQ | `4-alpine` | Message broker (2 durable queues: `encode`, `email`) |
| Nginx | `stable-alpine` | Frontend static serving + reverse proxy |
| nginx-rtmp (`tiangolo/nginx-rtmp`) | — | RTMP ingest + HLS/DASH live streaming output |
| Docker / Docker Compose | — | Containerisation and orchestration |

## Development Setup

### Prerequisites
- Docker Desktop (or Docker Engine + Compose plugin)
- No local Go / Node / Python installation required — all builds happen inside containers

### Running the Full Stack
```bash
docker-compose up --build
```
This brings up all 9 services. The client is available at `http://localhost:8080`.

### Port Map (Host → Container)
| Host Port | Service | Protocol |
|---|---|---|
| `8080` | `oes_client` (Nginx) | HTTP/WebSocket |
| `9000` | `oes_server` (Go API) | HTTP/WebSocket (also directly accessible) |
| `8887` | `oes_fileserver` | HTTP |
| `5672` | `messageq` (RabbitMQ) | AMQP |
| `50051` | `oes_isupport` (gRPC) | gRPC/HTTP2 |
| `1935` | `oes_liveserver` (nginx-rtmp) | RTMP |

### Configuration
Server configuration is in [`Server/config.json`](Server/config.json):
```json
{
  "database": {
    "MySqlDSN":    "root:root@tcp(db:3306)/OES",
    "RedisDSN":    "redis://cache:6379/0",
    "RabbitMQDSN": "amqp://rabbitmq:rabbitmq@messageq/",
    "GRPCDSN":     "isupport:50051"
  }
}
```
All DSNs use Docker service names as hostnames.

MQServer reads `RABBITMQ_DSN` env var (defaults to the same DSN if unset).
SMTP credentials for email: `SMTP_EMAIL` and `SMTP_PASSWORD` env vars on `oes_mqserver`.

### Build Stages
All services use **multi-stage Docker builds** to keep final images small:

| Service | Stage 1 | Stage 2 |
|---|---|---|
| `server` | `golang:1.26-alpine` — compile binary | `alpine:3.23` — binary + `config.json` only |
| `client` | `node:20-alpine` — `npm run build` | `nginx:stable-alpine` — dist + self-signed cert |
| `fileserver` | `golang:1.26-alpine` — compile binary | `alpine:3.23` — binary only |
| `mqserver` | `golang:1.26-alpine` — compile binary | `alpine:3.21` + ffmpeg |
| `isupport` | `python:3.14-alpine` + build tools — pip install | `python:3.14-alpine` + runtime libs only |
| `liveserver` | N/A (single stage) | `tiangolo/nginx-rtmp` + custom nginx.conf |
| `db` | N/A (single stage) | `mysql:8.0` + SQL init scripts |

### Shared Volume
`./media` host directory is bind-mounted into `server`, `fileserver`, and `mqserver` at `/app/OES/media`. This is the shared file store for:
- `media/video/<clientId>/<videoName>/` — MP4 originals (pre-encode) and HLS segments
- `media/examProofs/<clientId>/<examId>/<userId>/` — extracted webcam images
- `media/video/liveData/` — HLS live stream fragments (also mounted into `liveserver` at `/mnt/hls`)

## Technical Constraints
- **JWT secret** is hardcoded (`7yt65U745TR57lo9h%$fre#$TR43EW`) — must be externalised via env var before production use.
- **No TLS on gRPC** — `grpc.WithInsecure()` is used; acceptable only within the private Docker network.
- **Self-signed TLS cert** on Nginx client — browsers will show a warning; cookie `Secure: false` is set because of this.
- **JWT 15-minute expiry** with no refresh token mechanism; users must re-login after expiry.
- **WASM binary** (`main.wasm`) is pre-compiled and committed to `Client/oes/public/main.wasm`; must be rebuilt with `GOOS=js GOARCH=wasm go build` after changes to [`Client/wasmGo/main.go`](Client/wasmGo/main.go).
- **RabbitMQ channel** is not goroutine-safe — a `sync.Mutex` (`publishMu`) guards concurrent `Publish` calls in student service.
- **No horizontal scaling** for the WebSocket server — the Redis pub/sub bus is in place to support it, but the current single-server setup doesn't exercise multi-instance scaling.
- **MySQL healthcheck** uses `mysqladmin ping` with up to 10 retries and 30 s start period to handle slow MySQL initialisation.
