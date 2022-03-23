package mq

import (
	"github.com/streadway/amqp"
	"log"
	"testing"
)



func TestConsumeSub(t *testing.T) {
	sub := NewRabbitMQPubSub("amqp://zktube:zktube@127.0.0.1:5672/machine","create")
	defer sub.Close()
	sub.CallBackFunc = func(delivery amqp.Delivery) {
		//模拟1秒接受一个任务
		//time.Sleep(1 * time.Second)
		log.Printf("Received a message :%s", delivery.Body)
	}
	sub.ConsumeSub()

}
