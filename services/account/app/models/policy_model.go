package models

import "time"

type Policy struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
}
