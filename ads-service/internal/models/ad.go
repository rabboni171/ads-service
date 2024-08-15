package models

import "time"

type Ad struct {
	Id          int       `json:"id,omitempty" swaggerignore:"true"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       int       `json:"price,omitempty"`
	Photos      []string  `json:"photos,omitempty"`
	UserId      int       `json:"user_id,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty" swaggerignore:"true"`
}
