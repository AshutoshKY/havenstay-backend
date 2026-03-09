package host

import (
	"time"
)

// Entity is the core domain representation of a Host
type Entity struct {
	ID         int64
	UserName   string
	Email      string
	Phone      string
	Location   string
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
