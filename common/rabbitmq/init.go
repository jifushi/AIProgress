package rabbitmq

var RMQMessage *RabbitMQ

func Init() {
	// 创建MQ
	// aichat作为生产者将消息写入消息队列中
	RMQMessage = NewWorkRabbitMQ("ai_progress")
	// 异步启动消费者，将消息队列中的数据写入数据库中
	go RMQMessage.Consume(MQMessage)
}

// DestroyRabbitMQ 销毁RabbitMQ
func DestroyRabbitMQ() {
	RMQMessage.Destroy()
}
