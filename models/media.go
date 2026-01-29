package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Media representa una serie, manga, novela, etc.
type Media struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title    string             `bson:"title" json:"title"`
	Type     string             `bson:"type" json:"type"`
	Progress string             `bson:"progress" json:"progress"`
	Link     string             `bson:"link" json:"link"`
}
