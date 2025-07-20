# Real-Time Chat App Backend (Golang + WebSocket + Redis + MongoDB)

This project is a real-time chat server built using Golang. It supports public chat broadcasting and private 1-on-1 messaging using WebSocket. Redis is used for Pub/Sub messaging, and MongoDB is used to persist message history. Authentication is handled using JWT tokens.

---

## 🚀 Features

- ✅ **Real-time messaging** with WebSockets
- ✅ **Private 1-to-1 chat** and **public chat broadcasting**
- ✅ **Redis Pub/Sub** for scalable message distribution
- ✅ **MongoDB integration** for storing chat history
- ✅ **JWT-based authentication**
- ✅ Graceful handling of client connect/disconnect
- ✅ REST API to retrieve chat history between two users

---

## 🛠️ Tech Stack

- **Golang**
- **Gin** – Web framework
- **WebSocket** – Real-time communication
- **Redis** – Pub/Sub system
- **MongoDB** – NoSQL database to store messages
- **JWT (JSON Web Tokens)** – Authentication

---

## 📁 Project Structure

chat-app/
├── auth/ # JWT token creation and validation
├── models/ # Message and user model definitions
├── mongo/ # MongoDB connection setup
├── redis/ # Redis Pub/Sub functions
├── websocket/ # WebSocket logic, client management
├── main.go # App entry point and routes


## 🔧 Prerequisites

Make sure the following services are running:

- Redis on `localhost:6379`
- MongoDB on `mongodb://localhost:27017`



🙋‍♂️ Author
Dishank Agrawal
