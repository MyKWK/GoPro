package rabbitMQ

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"time"
)

// 构建实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error

	rabbitmq.Conn, err = amqp.Dial(rabbitmq.MQ_url)
	rabbitmq.FailOnErr(err, "Failed to connect to RabbitMQ")

	rabbitmq.Channel, err = rabbitmq.Conn.Channel()
	rabbitmq.FailOnErr(err, "Failed to open a channel")

	return rabbitmq
}

// 发送消息
func (r *RabbitMQ) TopicPublish(msg string, key string) {
	// 1.交换机
	err := r.Channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil)
	r.FailOnErr(err, "Failed to declare an exchange")

	// 2.发送：指定交换机和key
	err = r.Channel.Publish(
		r.Exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	r.FailOnErr(err, "Failed to publish a message")
}

// 接收 通过通配符去匹配
func (r *RabbitMQ) TopicConsume() {
	// 1. 交换机
	err := r.Channel.ExchangeDeclare(
		r.Exchange,
		"topic",
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
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	go func() {
		for d := range msgs {
			log.Printf("Queue: %s \t Key:%s  \t message: %s", queue.Name, r.Key, d.Body)
		}
	}()
}

func MainTopicPublish() {
	rabbitmq := NewRabbitMQTopic("Topic", "")
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			rabbitmq.TopicPublish(strconv.Itoa(i), "key.one")
		} else {
			rabbitmq.TopicPublish(strconv.Itoa(i), "key.two")
		}
		time.Sleep(1 * time.Second)
	}
}
func MainTopicConsume() {
	go func() {
		rabbitmq1 := NewRabbitMQTopic("Topic", "key.#")
		rabbitmq1.TopicConsume()
	}()
	go func() {
		rabbitmq2 := NewRabbitMQTopic("Topic", "key.one")
		rabbitmq2.TopicConsume()
	}()
}
