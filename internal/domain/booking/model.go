package booking

import (
	"time"
)

type Model struct {
	BookingID   int64     `gorm:"primaryKey;column:booking_id;autoIncrement"`
	PropertyID  int64     `gorm:"column:property_id;index"`
	GuestID     int64     `gorm:"column:guest_id;index"`
	StartDate   string    `gorm:"column:start_date;type:varchar(10)"` // YYYY-MM-DD
	EndDate     string    `gorm:"column:end_date;type:varchar(10)"`   // YYYY-MM-DD
	Status      string    `gorm:"column:status;type:varchar(50)"`
	PaymentInfo string    `gorm:"column:payment_info;type:varchar(255)"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Model) TableName() string {
	return "bookings"
}
