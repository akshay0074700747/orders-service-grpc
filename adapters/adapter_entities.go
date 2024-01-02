package adapters

import "gorm.io/gorm"

type OrderAdapter struct {
	DB *gorm.DB
}

func NewOrderAdapter(db *gorm.DB) *OrderAdapter {
	return &OrderAdapter{
		DB: db,
	}
}