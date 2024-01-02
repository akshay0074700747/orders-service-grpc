package db

import (
	"github.com/akshay0074700747/orders-service/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(connectTo string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(connectTo), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entities.Orders{})
	db.AutoMigrate(&entities.OrderItems{})

	return db, nil
}
