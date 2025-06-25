package rabbitMQ

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"time"
)

// 订阅模式实例
func NewRabbitMQPubSub(ExchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", ExchangeName, "")
	var err error

	rabbitmq.Conn, err = amqp.Dial(rabbitmq.MQ_url)
	rabbitmq.FailOnErr(err, "Failed to connect to RabbitMQ")

	rabbitmq.Channel, err = rabbitmq.Conn.Channel()
	rabbitmq.FailOnErr(err, "Failed to open a channel")

	return rabbitmq
}

// 订阅模式的消息生产
func (r *RabbitMQ) PubSubPublish(msg string) {
	// 1.交换机申请
	err := r.Channel.ExchangeDeclare(
		r.Exchange,
		"fanout", // （广播型）
		true,
		false,
		false,
		false,
		nil)
	r.FailOnErr(err, "Failed to declare a exchange")

	// 2.发送消息
	err = r.Channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	r.FailOnErr(err, "Failed to publish a message")
}

// 订阅模式的消费
func (r *RabbitMQ) PubSubConsume() {
	// 1. 尝试创建交换机
	err := r.Channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	r.FailOnErr(err, "Failed to declare a exchange")

	// 2.创建队列，队列必须随机生成
	// 每个消费者可以有自己的队列
	queue, err := r.Channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	r.FailOnErr(err, "Failed to declare a queue")
	// 3.队列与Exchange绑定,fanout特性就是如此，它会广播to所有绑定的队列
	err = r.Channel.QueueBind(
		queue.Name,
		"",
		r.Exchange,
		false,
		nil)
	r.FailOnErr(err, "Failed to bind a queue")
	// 4.获取消息
	msgs, err := r.Channel.Consume(
		r.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
}

// demo演示
func MainPubSubPublish() {
	rabbitmq := NewRabbitMQPubSub("PubSub")
	rabbitmq.PubSubPublish("Hello")
	for i := 0; i < 1000; i++ {
		rabbitmq.PubSubPublish(strconv.Itoa(i))
		time.Sleep(1 * time.Second)
	}
}
func MainPubSubConsume() {
	go func() {
		rabbitmq1 := NewRabbitMQPubSub("PubSub")
		rabbitmq1.PubSubConsume()
	}()
	go func() {
		rabbitmq2 := NewRabbitMQPubSub("PubSub")
		rabbitmq2.PubSubConsume()
	}()
}
