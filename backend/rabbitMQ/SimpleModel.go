package rabbitMQ

import (
	"fmt"
	"strconv"
	"time"
)

func MainSimplePublish() {
	rabbitSimple := NewRabbitMQSimple("simpleMode")
	rabbitSimple.PublishSimple("Hello World!")
	for i := 0; i < 1000; i++ {
		fmt.Println("Publishing message: ", i)
		rabbitSimple.PublishSimple(strconv.Itoa(i))
		time.Sleep(1 * time.Second)
	}
}

func MainSimpleReceive() {
	rabbitMQ := NewRabbitMQSimple("simpleMode")
	rabbitMQ.ConsumeSimple()
}
