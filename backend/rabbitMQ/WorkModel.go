package rabbitMQ

func MainWorkReceive() {
	rabbitMQ1 := NewRabbitMQSimple("simpleMode")
	rabbitMQ2 := NewRabbitMQSimple("simpleMode")
	go func() {
		rabbitMQ1.ConsumeSimple()
	}()
	go func() {
		rabbitMQ2.ConsumeSimple()
	}()
}
