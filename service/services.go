package service

import (
	"context"
	"fmt"

	"github.com/akshay0074700747/orders-service/adapters"
	"github.com/akshay0074700747/orders-service/entities"
	"github.com/akshay0074700747/proto-files-for-microservices/pb"
	"github.com/golang/protobuf/ptypes/empty"
)

type OrderService struct {
	Adapter *adapters.OrderAdapter
	pb.UnimplementedOrderServiceServer
}

var (
	userConn    pb.UserServiceClient
	productConn pb.ProductServiceClient
)

func Initializer(userconn pb.UserServiceClient, productconn pb.ProductServiceClient) {
	userConn = userconn
	productConn = productconn
}

func NewOrderService(adapter *adapters.OrderAdapter) *OrderService {
	return &OrderService{
		Adapter: adapter,
	}
}

func (order *OrderService) AddOrder(ctx context.Context, request *pb.AddOrderRequest) (*pb.AddOrderResponce, error) {

	req := entities.Orders{
		UserID: uint(request.UserId),
		Status: "Ordered",
	}

	var products []*pb.AddProductResponce
	var prodIds []uint

	for _, productid := range request.ProductIDs {

		prodIds = append(prodIds, uint(productid))
		product, err := productConn.GetProduct(ctx, &pb.GetProductByID{
			Id: uint32(productid),
		})
		if err != nil {
			fmt.Println(err.Error())
			return &pb.AddOrderResponce{}, err
		}
		products = append(products, product)
		req.Amount += int(product.Price)

	}

	userRes, err := userConn.GetUser(ctx, &pb.UserRequest{
		Id: request.UserId,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	orderRes, err := order.Adapter.AddOrdersAdapter(req, prodIds)
	if err != nil {
		return &pb.AddOrderResponce{}, err
	}

	responce := pb.AddOrderResponce{
		Id:       uint32(orderRes.Id),
		User:     userRes,
		Price:    int32(orderRes.Amount),
		Products: products,
	}

	return &responce, nil

}

func (order *OrderService) GetOrdersByUser(ctx context.Context, request *pb.GetOrdersByUserRequest) (*pb.GetOrdersWithoutUser, error) {

	userID := request.UserId

	orders, err := order.Adapter.GetOrderFromUser(uint(userID))
	if err != nil {
		return &pb.GetOrdersWithoutUser{}, err
	}

	if len(orders) == 0 {
		return &pb.GetOrdersWithoutUser{}, fmt.Errorf("You havent ordered anything")
	}

	var ordersRes []*pb.OrdersOnly
	for _, order := range orders {
		var products []*pb.AddProductResponce
		for _, prodId := range order.Products {
			product, err := productConn.GetProduct(ctx, &pb.GetProductByID{
				Id: uint32(prodId),
			})
			if err != nil {
				fmt.Println(err.Error())
				return &pb.GetOrdersWithoutUser{}, err
			}
			products = append(products, product)
		}
		responce := &pb.OrdersOnly{
			Id:       uint32(order.Id),
			Price:    int32(order.Amount),
			Status:   order.Status,
			Products: products,
		}
		ordersRes = append(ordersRes, responce)
	}

	return &pb.GetOrdersWithoutUser{
		Orders: ordersRes,
	}, nil

}

func (order *OrderService) GetAllOrdersResponce(ctx context.Context, request *empty.Empty) (*pb.GetOrdersByUserResponce, error) {

	orders, err := order.Adapter.GetOrders()
	if err != nil {
		return &pb.GetOrdersByUserResponce{}, err
	}

	var results []*pb.AddOrderResponce

	for _, order := range orders {
		var products []*pb.AddProductResponce
		for _, prodId := range order.Products {
			product, err := productConn.GetProduct(ctx, &pb.GetProductByID{
				Id: uint32(prodId),
			})
			if err != nil {
				fmt.Println(err.Error())
				return &pb.GetOrdersByUserResponce{}, err
			}
			products = append(products, product)
		}
		user, err := userConn.GetUser(ctx, &pb.UserRequest{
			Id: uint32(order.UserID),
		})
		if err != nil {
			return &pb.GetOrdersByUserResponce{}, err
		}
		res := &pb.AddOrderResponce{
			Id:       uint32(order.Id),
			User:     user,
			Price:    int32(order.Amount),
			Products: products,
		}

		results = append(results, res)
	}

	return &pb.GetOrdersByUserResponce{
		Orders: results,
	}, nil
}
