package property

import (
	"time"
)

// Entity is the core domain representation of a Property
type Entity struct {
	ID          int64
	HostID      int64
	Name        string
	Location    string
	Description string
	Price       int64
	CreatedAt   time.Time
}
