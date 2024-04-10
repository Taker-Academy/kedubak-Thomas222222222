package structures

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comments struct {
	ID        int    `bson:"id" json:"id"`
	FirstName string `bson:"firstName" json:"firstName"`
	Content   string `bson:"content" json:"content"`
}

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CreateAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UserID    string             `bson:"userId" json:"userId"`
	FirstName string             `bson:"firstName" json:"firstName"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Comments  []Comments         `bson:"comments" json:"comments"`
	UpVotes   []string           `bson:"upVotes" json:"upVotes"`
}
