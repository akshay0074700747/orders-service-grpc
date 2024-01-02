package initializer

import (
	"github.com/akshay0074700747/orders-service/adapters"
	"github.com/akshay0074700747/orders-service/service"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB) *service.OrderService {
	adapter := adapters.NewOrderAdapter(db)
	service := service.NewOrderService(adapter)

	return service
}
