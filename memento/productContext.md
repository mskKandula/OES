# OES — Product Context

## Why This Project Exists
OES was created as a technology-focused proof-of-concept to demonstrate how a modern, distributed online examination and learning platform can be built using production-grade open-source tools. The secondary goal is to evolve it into a fully functional LMS with AI-assisted features (NLP, CV).

## Problems It Solves

| Problem | Solution |
|---|---|
| Conducting exams remotely without physical invigilation | BlackBox invigilation (webcam image capture + zip upload) + WebRTC video proctoring |
| Manually registering hundreds of students | Bulk Excel upload with automated welcome email via RabbitMQ async job |
| Static question papers being leaked | Questions uploaded as `.txt`, cached in-memory server-side, fetched per session |
| Teaching without a physical classroom | WebRTC video call + WebSockets whiteboard |
| Students with typing difficulties | On-screen virtual keyboard + WebSpeechRecognition voice input |
| Evaluating answers without a human grader | WebAssembly (Go → WASM) word-counter evaluation runs entirely in browser |
| Serving video to students at variable bandwidths | HLS adaptive bitrate streaming (360p / 480p / 720p) via ffmpeg transcoding |
| AI-powered question generation from text | gRPC stub to a Python service (placeholder, NLP not yet implemented) |

## How It Should Work — User Journeys

### Examiner Journey
1. **Sign up** → POST `/o/signUp` → account created in `Examiners` table with a UUID `clientId`.
2. **Login** → POST `/o/login` → JWT set in `HttpOnly` cookie; `userType=Examiner`, `clientId`, `userId` returned.
3. **Get dynamic menu** → GET `/r/getRoutes` → Redis-cached role-based menu items fetched from MySQL.
4. **Register students** → POST `/r/multipleStudentsRegistration` (Excel `.xlsx`) → concurrent bcrypt hash per student → DB insert → `email` message published to RabbitMQ → MQServer sends welcome email via Gmail SMTP.
5. **Upload question paper** → POST `/r/uploadQuestionFile` → `.txt` parsed line by line → exam record created in `Exams` table → questions cached in-memory `QuestionCache` (map keyed by `examId`).
6. **Upload video** → POST `/r/uploadVideoContent` → MP4 saved to `media/video/` → video metadata stored in MySQL → `encode` message published to RabbitMQ → MQServer runs `ffmpeg` to generate 360p/480p/720p HLS streams.
7. **Monitor students** → WebSocket connection at `ws://.../r/ws` → real-time join/leave events, notifications, whiteboard data, 1-to-1 WebRTC signals.
8. **Generate AI questions** → POST `/r/questionGeneration` → gRPC call to `isupport:50051` (Python service).

### Student Journey
1. **Login** (same endpoint, same JWT mechanism, `userType=Student`).
2. **Get dynamic menu** → different menu items (StudentDashboard, Exam, VideoContent, WhiteBoard, BroadcastVideo).
3. **Take exam** → GET `/r/getQuestions?examId=X` → answer per question in textarea / voice / on-screen keyboard → on Submit, **Go WASM** (`countWords` function) counts matching keyword answers in browser → result displayed in `/wordCounter`.
4. **Invigilation** → webcam stream captured every 5 seconds (10 frames total) → zipped with JSZip → POST `/r/uploadExamProof` → server unzips async via worker goroutine pool → image paths stored in `StudentExamProofs`.
5. **Watch videos** → GET `/r/getVideos` → Redis-cached list → HLS `.m3u8` playlist served from FileServer at `/cdn/...`.
6. **Video captions** → WebkitSpeechRecognition overlaid on Canvas during video playback.
7. **Collaborate** → WebSocket for whiteboard (type 3 messages), text chat (type 2), broadcast (type 4).

## User Experience Goals
- **Zero friction login** — one form, cookie-based session, role-detected menu rendered automatically.
- **Responsive real-time feedback** — WebSocket notifications, online-user list, whiteboard sync all happen without page reload.
- **Browser-native features** — voice input, on-screen keyboard, image capture, WebAssembly evaluation all run without any plugin.
- **Adaptive video** — HLS ensures smooth playback on slow connections.
- **Async heavy work** — email sending and video encoding never block the HTTP response.
