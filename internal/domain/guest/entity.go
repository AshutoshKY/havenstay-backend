package guest

import (
	"time"
)

type Entity struct {
	GuestID   int64
	Name      string
	Email     string
	Phone     string
	Country   string
	CreatedAt time.Time
}
