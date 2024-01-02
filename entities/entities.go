package entities

type Orders struct {
	Id       uint `gorm:"primaryKey"`
	UserID   uint /*`gorm:"foreignKey:UserID;references:users(ID)"`*/
	Amount   int
	Status   string
}

type OrderItems struct {
	Id       uint `gorm:"primaryKey"`
	OrderID   uint `gorm:"foreignKey:OrderID;references:orders(Id)"`
	ProductID uint
}
