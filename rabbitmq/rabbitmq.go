package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"

	"awesomeProject/datamodels"
	"awesomeProject/services"
	"encoding/json"
	"os"
	"strconv"
	"sync"
)

// 连接信息
const MQURL = "amqp://root:1368@127.0.0.1:5672/vhost"

// rabbitMQ结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
	sync.Mutex
	// 用来做确认
	ackCh chan amqp.Confirmation // 表示Broker 对一条发布消息的确认信息
	retCh chan amqp.Return       // 表示一条无法被路由到任何队列的消息被 Broker 退回时携带的完整信息
}

// 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

// 断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
	}
}

// 创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabb"+
		"itmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	// 开启 confirm 模式（一次即可）
	err = rabbitmq.channel.Confirm(false)
	rabbitmq.failOnErr(err, "channel confirm mode failed")

	rabbitmq.ackCh = rabbitmq.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	// 指示 broker 对每条发布按顺序回一个 Ack/Nack；
	//这个通道里异步接收 Broker 发回的 发布确认事件（Ack/Nack）（接受/拒绝）
	rabbitmq.retCh = rabbitmq.channel.NotifyReturn(make(chan amqp.Return, 1)) //
	// 用来接收 Broker 退回的消息
	return rabbitmq
}

// 直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		true,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		return err
	}
	//调用channel 发送消息到队列中
	err = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		true,  // mandatory: 无人接受则返回消息给发布者
		false, // immediate 已废弃，保持 false
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent, // 你已设置：持久化消息
			Body:         []byte(message),
		},
	)
	if err != nil {
		return err
	}

	// 等待 Return/Ack/Nack/超时（谁先来用谁）
	select {
	case ret := <-r.retCh:
		// 不可路由（通常是 exchange/routingKey/绑定错）
		return fmt.Errorf("publish returned (unroutable): code=%d text=%s key=%s",
			ret.ReplyCode, ret.ReplyText, ret.RoutingKey)

	case c := <-r.ackCh:
		if !c.Ack {
			return fmt.Errorf("publish NACKed by broker")
		}
		// Ack=确认成功（对持久消息+持久队列，表示已安全接收并写盘）

	case <-time.After(2 * time.Second):
		// 超时：视为失败，让上层重试（建议做指数退避）
		return fmt.Errorf("publish confirm timeout")
	}

	return nil
}

// ConsumeSimple simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
	}

	// 可配置的 prefetch（默认 32，可通过环境变量 MQ_PREFETCH 覆盖）
	prefetch := 32
	if v := os.Getenv("MQ_PREFETCH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			prefetch = n
		}
	}
	// 在ACK之前，mq还能推送多少信息给当前消费者
	r.channel.Qos(
		prefetch, // 未ack消息条数
		0,        // 未ack消息最大字节
		false,    //如果设置为true ，则当前连接的所有channel共享这套限流
	)

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		//这里要改掉，我们用手动应答
		false, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	// 并发消费：启动多个 worker 读取同一个 msgs（默认 16 个，可用 MQ_WORKERS 调整）
	workers := 16
	if v := os.Getenv("MQ_WORKERS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			workers = n
		}
	}
	for i := 0; i < workers; i++ {
		go func() {
			for d := range msgs {
				//消息逻辑处理，可以自行设计逻辑
				log.Printf("Received a message: %s", d.Body)
				message := &datamodels.Message{}
				err := json.Unmarshal(d.Body, message) // 没有tag的话，则按照字段名首字母小写匹配
				if err != nil {
					fmt.Println(err)
				}
				//插入订单
				_, err = orderService.InsertOrderByMessage(message)
				if err != nil {
					fmt.Println(err)
				}
				//扣除商品数量
				err = productService.SubNumberOne(message.ProductID)
				if err != nil {
					fmt.Println(err)
				}
				//如果为true表示确认所有未确认的消息，
				//为false表示确认当前消息
				d.Ack(false)
			}
		}()
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}

}
