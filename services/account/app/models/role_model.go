package models

import "time"

type Role struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
}
