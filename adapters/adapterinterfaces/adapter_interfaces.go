package adapterinterfaces

import (
	"github.com/akshay0074700747/orders-service/entities"
	"github.com/akshay0074700747/orders-service/responce"
)

type AdapterInterface interface {
	AddOrdersAdapter(orderReq entities.Orders, prodIDs []uint) (responce.ProductRes, error)
	GetOrders() ([]responce.ProductRes, error)
	GetOrderFromUser(userID uint) ([]responce.ProductRes, error)
}
