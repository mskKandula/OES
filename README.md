# OES

Online Examination System

An attempt to create OES, Which is a platform to conduct Exams seamlessly & effortlessly. Started as OES & Scaling as LMS (Learning Management System).

        Note: Focused mainly on technologies & used as POC implementations, but in a systematic way.

## Actors :

     - Examiner
     - Student

## Features :

     Examiner :

    	  - Multiple Students Registration through Excel & Email Notification
    	  - Classroom teaching through VideoCall (WebRTC)
    	  - WhiteBoard (WebSockets & Canvas)
    	  - Upload Videos
    	  - Upload a Question Paper
    	  - Students Status & Video Proctoring Online through WebRTC
    	  - Realtime Notifications to Students (WebSockets)
    	  - Realtime Chart Dashboards


      Student :

    	  - View uploaded Videos On Demand (HLS(Http Live Streaming))
    	  - Dynamic video captions (WebkitSpeechRecognition, Canvas)
    	  - One to One & One to Many Video/Text Chat (WebRTC & WebSockets)
    	  - Online Examination (Voice Typing & On Screen Keyboard)
    	  - BlackBox Invigilation (Captures the images through webcam & Uploads a zip file to server)
    	  - Online Answer Evaluation (WASM (WebAssembly) )


      Common :

    	  - Login Functionality (JWT Authentication)
    	  - Dynamic MenuList (User Based)
    	  - UserBased Authorization to API's

## Scope & Ideas yet to implement:

Fully functional LMS & Incorporation of Deep Learning to develop a novel quality.

        - Document Similarity (NLP)
        - Text Summarizer (NLP)
        - Automatic Question Generation (NLP)
        - Student Mood Analyzer while Exam (CNN)

## Technologies:

      Front-End : Vue.js, WebAssembly
      Back-End : Golang
      DB	: MySQL
      Web Server : Nginx
      Caching & Scaling : Redis
      Message Queue : RabbitMQ
      Deep Learning : Python
      Other : Docker, gRPC, HLS, WebSockets, WebRTC

## Conclusion:

      Clean Architecture, Concurrent & Highly Scalable Distributed System.
