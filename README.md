# 💬 Real-Time Chat App Backend (Golang + WebSocket + Redis + MongoDB)

A robust real-time chat server built in **Golang**, supporting **public broadcasting** and **private 1-to-1 messaging** over WebSockets. It features JWT-based authentication, Redis-powered message broadcasting, and MongoDB-based persistent message storage.

---

## 🚀 Features

- ⚡ Real-time messaging using WebSockets
- 🔒 Private 1-to-1 chats
- 🌐 Public broadcast chat support
- 🔁 Scalable Redis Pub/Sub for multi-instance support
- 🗃️ MongoDB integration to persist chat history
- 🔐 JWT-based user authentication
- 🔄 Graceful client connection and disconnection handling
- 📜 REST API to fetch private and public chat history
- 📂 Clean modular structure for scalability and clarity

---

## 🛠️ Tech Stack

| Component      | Technology            |
| -------------- | --------------------- |
| Programming    | Go (Golang)           |
| Web Framework  | Gin                   |
| Real-time Comm | Gorilla WebSocket     |
| Messaging      | Redis (Pub/Sub)       |
| Database       | MongoDB               |
| Auth           | JWT (JSON Web Tokens) |

---

## 📁 Project Structure

```bash
chat-app/
│
├── auth/           # JWT generation and validation
├── models/         # Message/User schema definitions
├── mongo/          # MongoDB connection setup
├── redis/          # Redis Pub/Sub logic
├── websocket/      # WebSocket handling & client manager
├── api/            # REST APIs (e.g., login, chat history)
├── cmd/            # App entry point
└── main.go         # Server bootstrap and route registration
```

## 🔧 Prerequisites

Make sure the following services are running:

- Redis on `localhost:6379`
- MongoDB on `mongodb://localhost:27017`

👨‍💻 Author
Dishank Agrawal
