package services

import (
	"awesomeProject/datamodels"
	"awesomeProject/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
}

type OrderService struct {
	orderRepository repositories.IOrderRepository
}

func NewOrderService(orderRepository repositories.IOrderRepository) IOrderService {
	return &OrderService{
		orderRepository: orderRepository,
	}
}

func (s *OrderService) GetOrderByID(id int64) (*datamodels.Order, error) {
	return s.orderRepository.SelectByKey(id)
}

func (s *OrderService) DeleteOrderByID(id int64) bool {
	order := &datamodels.Order{ID: id}
	return s.orderRepository.Delete(order)
}

func (s *OrderService) UpdateOrder(order *datamodels.Order) error {
	return s.orderRepository.Update(order)
}

func (s *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return s.orderRepository.Insert(order)
}

func (s *OrderService) GetAllOrder() ([]*datamodels.Order, error) {
	return s.orderRepository.SelectAll()
}

func (s *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {

	return s.orderRepository.SelectAllWithInfo()
}
