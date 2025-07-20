package websocket

import (
	"chat-app/auth"
	"chat-app/models"
	"chat-app/mongo"
	"chat-app/redis"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(c *gin.Context) {
	// Extract JWT token from header
	tokenStr := c.Request.Header.Get("Sec-WebSocket-Protocol")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Validate token
	token, claims, err := auth.VerifyJWT(tokenStr)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract username from JWT claims
	username := claims.Username

	// Upgrade to WebSocket
	responseHeader := http.Header{}
	responseHeader.Set("Sec-WebSocket-Protocol", "Bearer "+tokenStr)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, responseHeader)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}
	defer conn.Close()

	AddClient(conn)
	RegisterUser(username, conn)
	defer RemoveClient(conn)

	go listenRedis()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Parse incoming JSON message
		var incoming struct {
			To      string `json:"to"`      // Optional: empty means broadcast
			Content string `json:"content"` // Required
		}
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Println("Invalid JSON format:", err)
			continue
		}

		// Save message to MongoDB
		message := models.Message{
			From:      username,
			To:        incoming.To,
			Content:   incoming.Content,
			Timestamp: time.Now(),
		}
		_, err = mongo.MessageCollection.InsertOne(context.Background(), message)
		if err != nil {
			log.Println("MongoDB Insert error:", err)
		}

		// Redis publish (for broadcast OR private)
		messageJSON, _ := json.Marshal(message)
		if incoming.To == "" {
			redis.PublishMessage("chat", string(messageJSON))
		} else {
			SendPrivateMessage(incoming.To, string(messageJSON))
		}
	}
}

// func HandleWebSocket(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("Upgrade error:", err)
// 		return
// 	}

// 	log.Println("WebSocket connected!")
// 	for {
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Read error:", err)
// 			return
// 		}

// 		log.Println("Received:", string(msg))
// 	}
// }

// func HandleWebSocket(c *gin.Context) {
// 	// Extract JWT token from header
// 	tokenStr := strings.TrimSpace(c.GetHeader("Sec-WebSocket-Protocol"))
// 	if tokenStr == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
// 		return
// 	}

// 	// Validate token
// 	token, claims, err := auth.VerifyJWT(tokenStr)
// 	if err != nil || !token.Valid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		return
// 	}

// 	// Extract username from JWT claims
// 	username := claims.Username

// 	// Upgrade to WebSocket
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, http.Header{
// 		"Sec-WebSocket-Protocol": {tokenStr},
// 	})
// 	if err != nil {
// 		log.Println("WebSocket Upgrade error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	AddClient(conn)
// 	RegisterUser(username, conn)
// 	defer RemoveClient(conn)

// 	go listenRedis()

// 	for {
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Read error:", err)
// 			break
// 		}

// 		// Parse incoming JSON message
// 		var incoming struct {
// 			To      string `json:"to"`      // Optional: empty means broadcast
// 			Content string `json:"content"` // Required
// 		}
// 		if err := json.Unmarshal(msg, &incoming); err != nil {
// 			log.Println("Invalid JSON format:", err)
// 			continue
// 		}

// 		// Save message to MongoDB
// 		message := models.Message{
// 			From:      username,
// 			To:        incoming.To,
// 			Content:   incoming.Content,
// 			Timestamp: time.Now(),
// 		}
// 		_, err = mongo.MessageCollection.InsertOne(context.Background(), message)
// 		if err != nil {
// 			log.Println("MongoDB Insert error:", err)
// 		}

// 		// Redis publish (for broadcast OR private)
// 		messageJSON, _ := json.Marshal(message)
// 		if incoming.To == "" {
// 			redis.PublishMessage("chat", string(messageJSON))
// 		} else {
// 			SendPrivateMessage(incoming.To, string(messageJSON))
// 		}
// 	}
// }

func listenRedis() {
	sub := redis.Subscribe("chat")
	ch := sub.Channel()

	for msg := range ch {
		BroadcastToClients(msg.Payload)
	}
}

func GetChatHistory(c *gin.Context) {
	tokenStr := strings.TrimSpace(c.GetHeader("Authorization"))
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Remove "Bearer " prefix if exists
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, claims, err := auth.VerifyJWT(tokenStr)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	username := claims.Username

	toUser := c.Query("to")

	filter := bson.M{}

	if toUser != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"from": username, "to": toUser},
				{"from": toUser, "to": username},
			},
		}
	} else {
		filter = bson.M{"to": ""}
	}

	cursor, err := mongo.MessageCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer cursor.Close(context.Background())

	var messages []models.Message
	if err := cursor.All(context.Background(), &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
