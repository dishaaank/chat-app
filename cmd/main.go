package main

import (
	"chat-app/api"
	"chat-app/mongo"
	"chat-app/redis"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	mongo.ConnectMongo()
	redis.InitRedis()
	r := gin.Default()

	api.RegisterRoutes(r)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
