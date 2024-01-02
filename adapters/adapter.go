package adapters

import (
	"fmt"

	"github.com/akshay0074700747/orders-service/entities"
	"github.com/akshay0074700747/orders-service/responce"
)

func (order *OrderAdapter) AddOrdersAdapter(orderReq entities.Orders, prodIDs []uint) (responce.ProductRes, error) {

	var res entities.Orders
	query := "INSERT INTO orders (user_id,amount,status) VALUES($1,$2,$3) RETURNING id,user_id,amount,status"
	query1 := "INSERT INTO order_items (order_id,product_id) VALUES($1,$2)"

	tx := order.DB.Begin()
	err := order.DB.Raw(query, orderReq.UserID, orderReq.Amount, orderReq.Status).Scan(&res).Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return responce.ProductRes{}, err
	}
	for _, prod := range prodIDs {
		if order.DB.Exec(query1, res.Id, prod).Error != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return responce.ProductRes{}, err
		}
	}

	tx.Commit()

	return responce.ProductRes{
		Id:       res.Id,
		UserID:   res.UserID,
		Amount:   res.Amount,
		Status:   res.Status,
		Products: prodIDs,
	}, nil
}

func (order *OrderAdapter) GetOrders() ([]responce.ProductRes, error) {

	var responcee []responce.ProductRes
	var res []entities.Orders
	query := "SELECT * FROM orders"
	query1 := "SELECT product_id FROM order_items WHERE order_id = $1"

	err := order.DB.Raw(query).Scan(&res).Error
	if err != nil {
		fmt.Println(err.Error())
		return []responce.ProductRes{}, err
	}

	for i := range res {
		responcee = append(responcee, responce.ProductRes{
			Id:     res[i].Id,
			UserID: res[i].UserID,
			Amount: res[i].Amount,
			Status: res[i].Status,
		})
		err := order.DB.Raw(query1, responcee[i].Id).Scan(&responcee[i].Products).Error
		if err != nil {
			fmt.Println(err.Error)
			return []responce.ProductRes{}, err
		}
	}

	return responcee, err

}

func (order *OrderAdapter) GetOrderFromUser(userID uint) ([]responce.ProductRes, error) {

	var responcee []responce.ProductRes
	var ress []entities.Orders
	query := "SELECT * FROM orders WHERE user_id = $1"
	query1 := "SELECT product_id FROM order_items WHERE order_id = $1"

	res := order.DB.Raw(query, userID).Scan(&ress)

	if res.Error != nil {
		return []responce.ProductRes{}, res.Error
	}

	if res.RowsAffected == 0 {
		return []responce.ProductRes{}, fmt.Errorf("No user foud with the given userName")
	}

	for i := range ress {
		responcee = append(responcee, responce.ProductRes{
			Id:     ress[i].Id,
			UserID: ress[i].UserID,
			Amount: ress[i].Amount,
			Status: ress[i].Status,
		})
		if err := order.DB.Raw(query1, responcee[i].Id).Scan(&responcee[i].Products).Error; err != nil {
			fmt.Println(err.Error)
			return []responce.ProductRes{}, err
		}
	}

	return responcee, nil
}

