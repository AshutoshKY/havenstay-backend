package guest

import (
	"time"
)

type Model struct {
	GuestID   int64     `gorm:"primaryKey;column:guest_id;autoIncrement"`
	Name      string    `gorm:"column:name;type:varchar(255)"`
	Email     string    `gorm:"column:email;type:varchar(255);unique"`
	Phone     string    `gorm:"column:phone;type:varchar(255)"`
	Country   string    `gorm:"column:country;type:varchar(255)"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Model) TableName() string {
	return "guests"
}
