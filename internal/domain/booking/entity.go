package booking

import (
	"time"
)

type Entity struct {
	BookingID   int64
	PropertyID  int64
	GuestID     int64
	StartDate   string
	EndDate     string
	Status      string
	PaymentInfo string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
