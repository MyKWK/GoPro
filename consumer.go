package main

import (
	"awesomeProject/common"
	"awesomeProject/rabbitmq"
	"awesomeProject/repositories"
	"awesomeProject/services"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		panic(err)
	}
	productRepository := repositories.NewProductManager(db)
	productService := services.NewProductService(productRepository)
	orderRepository := repositories.NewOrderManagerRepository(db)
	orderService := services.NewOrderService(orderRepository)
	// rabbit
	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("iMooc")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}
