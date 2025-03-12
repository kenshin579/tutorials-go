package model

import "time"

type User struct {
	ID      string  `bson:"_id,omitempty"`
	Name    string  `bson:"name,omitempty"`
	Address Address `bson:"address,omitempty"`
}

type Address struct {
	Street    string    `bson:"street,omitempty"`
	City      string    `bson:"city,omitempty"`
	State     string    `bson:"state,omitempty"`
	VisitedAt time.Time `bson:"visitedAt,omitempty"`
}
