package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	From      string             `bson:"from" json:"from"`
	To        string             `bson:"to,omitempty" json:"to,omitempty"` // optional for public messages
	Content   string             `bson:"content" json:"content"`
	IsPrivate bool               `bson:"is_private" json:"is_private"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

func SaveMessage(collection *mongo.Collection, msg Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, msg)
	//if err != nil {
		return err
	//}
}
