package wishlist

import (
	"time"
)

type Model struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	GuestID    int64     `gorm:"column:guest_id;index"`
	PropertyID int64     `gorm:"column:property_id;index"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (Model) TableName() string {
	return "wishlists"
}
