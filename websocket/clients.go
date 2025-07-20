package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func AddClient(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	clients[conn] = true
}
func RemoveClient(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, conn)
}
func BroadcastToClients(message string) {
	mu.Lock()
	defer mu.Unlock()
	for conn := range clients {
		conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
var userConnections = make(map[string]*websocket.Conn) // username -> connection
//var mu sync.Mutex

func RegisterUser(username string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	userConnections[username] = conn
}

func SendPrivateMessage(to string, message string) {
	mu.Lock()
	defer mu.Unlock()

	if conn, ok := userConnections[to]; ok {
		conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
