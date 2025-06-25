package datamodels

type Order struct {
	ID          int64 `gorm:"column:ID;primaryKey;autoIncrement" json:"id" sql:"ID"`
	UserId      int64 `gorm:"column:userID"       json:"userId" sql:"userID"`
	ProductId   int64 `gorm:"column:productID"    json:"productId" sql:"productID"`
	OrderStatus int64 `gorm:"column:orderStatus"  json:"orderStatus" sql:"orderStatus"`
}

const (
	OrderWait    = iota
	OrderSuccess //1
	OrderFailed  //2
)
