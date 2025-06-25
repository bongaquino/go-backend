package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileAccess struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	FileID      primitive.ObjectID  `bson:"file_id"`                // The shared file
	OwnerID     primitive.ObjectID  `bson:"owner_id"`               // Who created the share
	RecipientID *primitive.ObjectID `bson:"recipient_id,omitempty"` // Nil if public/password-based share
	Password    *string             `bson:"password,omitempty"`     // Hashed password if set
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
}

func (FileAccess) GetIndexes() []primitive.D {
	return nil
}
