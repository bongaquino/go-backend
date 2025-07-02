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
	Access      string              `bson:"access"`          // "private", "public", "password", "email"
	Salt        string              `bson:"salt,omitempty"`  // Used for encrypted files
	Nonce       string              `bson:"nonce,omitempty"` // Used for encrypted files
	IsEncrypted bool                `bson:"is_encrypted"`
	IsDeleted   bool                `bson:"is_deleted"`
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
}

func (File) GetIndexes() []primitive.D {
	return nil
}
