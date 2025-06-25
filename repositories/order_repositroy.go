package repositories

import (
	"awesomeProject/common"
	"awesomeProject/datamodels"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	Insert(*datamodels.Order) (int64, error)
	Delete(*datamodels.Order) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManagerRepository struct {
	db *gorm.DB
}

func NewOrderManagerRepository(db *gorm.DB) *OrderManagerRepository {
	return &OrderManagerRepository{db: db}
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (int64, error) {
	result := o.db.Create(order)
	return order.ID, result.Error
}

func (o *OrderManagerRepository) Delete(order *datamodels.Order) bool {
	result := o.db.Delete(order)
	return result.Error == nil && result.RowsAffected > 0
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) error {
	result := o.db.Save(order)
	return result.Error
}

func (o *OrderManagerRepository) SelectByKey(productID int64) (*datamodels.Order, error) {
	order := &datamodels.Order{}
	result := o.db.Where("productID = ?", productID).First(order)
	return order, result.Error
}

func (o *OrderManagerRepository) SelectAll() ([]*datamodels.Order, error) {
	var orders []*datamodels.Order
	result := o.db.Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (o *OrderManagerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	sql := "Select o.ID,p.productName,o.orderStatus From imooc.order as o left join product as p on o.productID=p.ID"
	rows, errRows := o.db.Raw(sql).Rows()
	if errRows != nil {
		return nil, errRows
	}
	defer rows.Close()
	return common.GetResultRows(rows), err
}
