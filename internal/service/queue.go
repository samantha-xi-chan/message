package service

import (
	"github.com/message/internal"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func QueueConnInit(exhange string) (channel <-chan amqp.Delivery, amqp_channel *amqp.Channel) {

	// 连接rabbitmq
	conn, err := amqp.Dial(internal.AMQP_URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	// 创建信道，通常一个消费者一个
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	//defer ch.Close()

	// 声明交换机
	err = ch.ExchangeDeclare(
		exhange,  // 交换机名，需要跟消息发送方保持一致
		"fanout", // 交换机类型
		true,     // 是否持久化
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 声明需要操作的队列
	q, err := ch.QueueDeclare(
		"",    //"",   // 队列名字，不填则随机生成一个
		false, // 是否持久化队列
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 队列绑定指定的交换机
	err = ch.QueueBind(
		q.Name,  // 队列名
		"",      // 路由参数，fanout类型交换机，自动忽略路由参数
		exhange, // 交换机名字，需要跟消息发送端定义的交换器保持一致
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	// 创建消费者
	msg_channel, err := ch.Consume(
		q.Name, // 引用前面的队列名
		"",     // 消费者名字，不填自动生成一个
		true,   // 自动向队列确认消息已经处理
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	return msg_channel, ch
}
