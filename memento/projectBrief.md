# OES — Project Brief

## Project Name
**OES** — Online Examination System (scaling towards LMS — Learning Management System)

## One-Line Summary
A full-stack, containerised distributed platform to conduct online examinations and deliver learning content seamlessly, built as a proof-of-concept for a range of modern technologies.

## Core Goals
1. Allow **Examiners** to register students in bulk, upload question papers, upload video content, proctor live exams, and communicate via whiteboard/video.
2. Allow **Students** to take online exams (with voice typing & on-screen keyboard), watch on-demand HLS video, and participate in real-time communication.
3. Demonstrate a **clean, concurrent, highly-scalable distributed architecture** using Go, Vue.js, WebSockets, WebRTC, gRPC, RabbitMQ, Redis, HLS, and WebAssembly.

## Actors
| Actor | Primary Role |
|---|---|
| Examiner | Admin/Teacher — manages students, content, exams |
| Student | End-user — takes exams, consumes content |

## Scope

### Implemented Features

**Examiner:**
- Bulk student registration via Excel upload + email notification
- Classroom teaching via Video Call (WebRTC)
- Interactive WhiteBoard (WebSockets + Canvas)
- Upload Videos (MP4 → HLS transcoding via ffmpeg)
- Upload Question Paper (`.txt`)
- Student status & video proctoring (WebRTC)
- Real-time notifications to students (WebSockets)
- Real-time chart dashboards

**Student:**
- View On-Demand Videos (HLS adaptive bitrate streaming)
- Dynamic video captions (WebkitSpeechRecognition + Canvas)
- One-to-one & one-to-many Video/Text Chat (WebRTC + WebSockets)
- Online Examination (voice typing + on-screen keyboard)
- BlackBox invigilation (captures webcam images, zips and uploads to server)
- Online Answer Evaluation (WebAssembly — Go compiled to WASM)

**Common:**
- Login / Logout (JWT in HttpOnly cookie)
- Dynamic menu list (user role-based)
- Role-based API authorisation

### Future Scope (NLP / DL ideas)
- Document Similarity
- Text Summariser
- Automatic Question Generation (placeholder gRPC service exists)
- Student Mood Analyser (CNN)

## Constraints & Notes
- Focused on technology demonstration / POC — not hardened for production security.
- JWT secret is currently hardcoded in source; SMTP credentials via env vars.
- Self-signed TLS certificate generated at build time inside the Client container.
- The gRPC `isupport` service is a stub that echoes back the input; real NLP not yet implemented.
