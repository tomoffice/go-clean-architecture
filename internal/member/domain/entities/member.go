package entities

import "time"

type Member struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}
