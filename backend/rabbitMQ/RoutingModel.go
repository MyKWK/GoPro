package rabbitMQ

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"time"
)

// 构建实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error

	rabbitmq.Conn, err = amqp.Dial(rabbitmq.MQ_url)
	rabbitmq.FailOnErr(err, "Failed to connect to RabbitMQ")

	rabbitmq.Channel, err = rabbitmq.Conn.Channel()
	rabbitmq.FailOnErr(err, "Failed to open a channel")

	return rabbitmq
}

// 发送消息
func (r *RabbitMQ) RoutingPublish(msg string) {
	// 1.交换机
	err := r.Channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil)
	r.FailOnErr(err, "Failed to declare an exchange")

	// 2.发送：指定交换机和key
	err = r.Channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	r.FailOnErr(err, "Failed to publish a message")
}

// 接收
func (r *RabbitMQ) RoutingConsume() {
	// 1. 交换机
	err := r.Channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil)
	r.FailOnErr(err, "Failed to declare a exchange")

	// 2.创建队列，队名必须随机生成
	queue, err := r.Channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	r.FailOnErr(err, "Failed to declare a queue")
	// 3.队列绑定，direct特性--指定key
	err = r.Channel.QueueBind(
		queue.Name,
		r.Key,
		r.Exchange,
		false,
		nil)
	r.FailOnErr(err, "Failed to bind a queue")
	// 4.获取消息
	msgs, err := r.Channel.Consume(
		//r.QueueName,
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	go func() {
		for d := range msgs {
			log.Printf("%s Have received a message: %s", queue.Name, d.Body)
		}
	}()
}

func MainRoutingPublish() {
	rabbitmq := NewRabbitMQRouting("Routing", "RoutingKey")
	rabbitmqMini := NewRabbitMQRouting("Routing", "RoutingKeyMini")
	rabbitmq.RoutingPublish("Hello")
	rabbitmqMini.RoutingPublish("Hello")
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			rabbitmq.RoutingPublish(strconv.Itoa(i))
		} else {
			rabbitmqMini.RoutingPublish(strconv.Itoa(i))
		}
		time.Sleep(1 * time.Second)
	}
}
func MainRoutingConsume() {
	go func() {
		rabbitmq1 := NewRabbitMQRouting("Routing", "RoutingKey")
		rabbitmq1.RoutingConsume()
	}()
	go func() {
		rabbitmq2 := NewRabbitMQRouting("Routing", "RoutingKeyMini")
		rabbitmq2.RoutingConsume()
	}()
}
