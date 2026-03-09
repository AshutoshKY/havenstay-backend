package wishlist

import (
	"time"
)

type Entity struct {
	ID         int64
	GuestID    int64
	PropertyID int64
	CreatedAt  time.Time
}
