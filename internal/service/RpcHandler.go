package service

import (
	"encoding/json"
	"github.com/message/config"
	"github.com/message/internal/domain"
	"github.com/streadway/amqp"
	"log"
)

var ch_high *amqp.Channel
var ch_normal *amqp.Channel

func InitQueue() {

	_, ch_high = QueueConnInit(config.EXCHANGE_HIGH)
	_, ch_normal = QueueConnInit(config.EXCHANGE_NORMAL)

}

func OnNewStatus(sessionID string, timestamp int64, status int32, payload string) {
	data := domain.UpdateSessionStatus{
		SessionID: sessionID,
		Timestamp: timestamp,
		EvtType:   status,
		Payload:   payload,
	}
	bytes, err := json.Marshal(&data)
	if err == nil {
		//fmt.Println("json.Marshal 编码结果: ", string(bytes))
		enQueue(ch_high, config.EXCHANGE_HIGH, bytes)
	}

}

func OnNewFeed(sessionID string, timestamp int64, feed string) {
	data := domain.FeedSessionStream{
		SessionID: sessionID,
		Timestamp: timestamp,
		EvtType:   domain.SESSION_ING,
		Payload:   feed,
	}
	bytes, err := json.Marshal(&data)
	if err == nil {
		//fmt.Println("json.Marshal 编码结果: ", string(bytes))
		enQueue(ch_normal, config.EXCHANGE_NORMAL, bytes)
	}

}

func enQueue(amqp_channel *amqp.Channel, queue string, body []byte) {

	err = amqp_channel.Publish(
		queue, // exchange（交换机名字，跟前面声明对应）
		"",    // 路由参数，fanout类型交换机，自动忽略路由参数，填了也没用。
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json", // 消息内容类型，这里是普通文本
			Body:        body,               // 消息内容
		})

	if err == nil {
		log.Printf("sent:  %s", body)
	} else {
		log.Print(err)
	}
}

var err error
var ch *amqp.Channel
