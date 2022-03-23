package mq

import (
	"log"
	"testing"
	"time"
)

type Msg struct {
	CTime time.Time
}

var msg = `{"machine_uuid":"xxxxxxxxxx","pledge_address":"0x24d77de0e45132xxxxxxx83bded23c0b72bb5e"}`
func TestPublishSub(t *testing.T) {
	sub := NewRabbitMQPubSub("amqp://zktube:zktube@127.0.0.1:5672/machine","create")
	defer sub.Close()
	for {
		//date, _ := json.Marshal(Msg{
		//	CTime: time.Now(),
		//})
		err := sub.PublishSub([]byte(msg))
		if err !=  nil {
			log.Fatal(err)
		}
		//模拟1秒发送一个任务
		time.Sleep(1 * time.Second)
	}
}
