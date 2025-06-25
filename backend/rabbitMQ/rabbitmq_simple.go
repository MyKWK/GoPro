package rabbitMQ

import (
	"github.com/streadway/amqp"
	"log"
)

const MQURL = "amqp://guest:guest@localhost:5672/myvhost"

type RabbitMQ struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	MQ_url    string
}

func NewRabbitMQ(QueueName, Exchange, Key string) *RabbitMQ {
	return &RabbitMQ{
		QueueName: QueueName,
		Exchange:  Exchange,
		Key:       Key,
		MQ_url:    MQURL,
	}
}

func (r *RabbitMQ) Disconnect() {
	err := r.Channel.Close()
	if err != nil {
		return
	}
	err = r.Conn.Close()
	if err != nil {
		return
	}
}

func (r *RabbitMQ) FailOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func NewRabbitMQSimple(QueueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(QueueName, "", "")
	var err error
	rabbitmq.Conn, err = amqp.Dial(rabbitmq.MQ_url)
	rabbitmq.FailOnErr(err, "Failed to connect to RabbitMQ")
	rabbitmq.Channel, err = rabbitmq.Conn.Channel()
	rabbitmq.FailOnErr(err, "Failed to open a Channel")
	return rabbitmq
}
func (r *RabbitMQ) PublishSimple(Message string) {
	_, err := r.Channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil)
	r.FailOnErr(err, "Failed to declare a queue")
	if r.Key == "" {
		r.Key = r.QueueName
	}
	r.Channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(Message),
		})
}
func (r *RabbitMQ) ConsumeSimple() {
	_, err := r.Channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil)
	r.FailOnErr(err, "Failed to declare a queue")
	msgs, err := r.Channel.Consume(
		r.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil)
	// 只读通道，元素类型是 Delivery
	r.FailOnErr(err, "Failed to register a consumer")
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Println("Waiting for messages. To exit press CTRL+C")
}
