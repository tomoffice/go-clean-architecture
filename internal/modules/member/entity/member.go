package entity

import "time"

type Member struct {
	ID        int       ` json:"id"`
	Name      string    ` json:"name"`
	Email     string    ` json:"email"`
	Password  string    ` json:"-"`
	CreatedAt time.Time ` json:"created_at"`
}
