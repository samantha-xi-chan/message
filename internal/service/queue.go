package service

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func QueueConnInit(amqpUrl string, exhange string) (channel <-chan amqp.Delivery, amqp_channel *amqp.Channel) {
	conn, err := amqp.Dial(amqpUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	//defer ch.Close()

	err = ch.ExchangeDeclare(
		exhange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"", // 路由参数，fanout类型交换机，自动忽略路由参数
		exhange,
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	// 创建消费者
	msg_channel, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	return msg_channel, ch
}

func ProdQueueConnInit(amqpUrl string, exhange string) (channel <-chan amqp.Delivery, amqp_channel *amqp.Channel) {
	conn, err := amqp.Dial(amqpUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	//defer ch.Close()

	err = ch.ExchangeDeclare(
		exhange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	return nil, ch
}
