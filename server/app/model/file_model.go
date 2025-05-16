package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	UserID      primitive.ObjectID  `bson:"user_id"`
	DirectoryID *primitive.ObjectID `bson:"directory_id,omitempty"`
	Name        string              `bson:"name"`
	Hash        string              `bson:"hash"`
	Size        int64               `bson:"size"`
	ContentType string              `bson:"content_type"`
	Status      string              `bson:"status"`
	IsDeleted   bool                `bson:"is_deleted"`
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
}

func (File) GetIndexes() []primitive.D {
	return nil
}
