package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//RabbitMQ ...
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列的名称
	QueueName string
	//交换器
	Exchange string
	//key
	key string
	//连接信息
	Mqurl string
	//callback
	CallBackFunc func(delivery amqp.Delivery)
}

//NewRabbitMQ ..创建Rabbitmq实例
func NewRabbitMQ(hostUrl, queuename, exchange, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queuename,
		Exchange:  exchange,
		key:       key,
		Mqurl:     hostUrl,
	}
	var err error
	//创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.Err(err, "NewRabbitMQSimple-->创建连接错误")
	//创建隧道
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.Err(err, "NewRabbitMQSimple-->获取隧道错误")
	return rabbitmq
}

//Close ..断开连接
func (r *RabbitMQ) Close() {
	r.conn.Close()
	r.channel.Close()
}

//Err 错误日志
func (r *RabbitMQ) Err(err error, message string) {
	if err != nil {
		log.Fatal(message, err)
	}
}

//NewRabbitMQPubSub ..订阅模式创建Rabbitmq实例
func NewRabbitMQPubSub(hostUrl,queuename string) *RabbitMQ {
	rbmq := NewRabbitMQ(hostUrl,queuename, "", "")
	return rbmq
}

//ConsumeSub ..订阅模式下消费者
func (r *RabbitMQ) ConsumeSub() {
	//接收消息
	mesgs, err := r.channel.Consume(
		r.QueueName,
		"",    //区分多个消费者
		true,  //是否自动应答，
		false, //排他性
		false, //为true，不能将同一个connection中发送的消息传递给这个connection中的消费者
		false, //是否设置为阻塞
		nil)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for d := range mesgs {
			//实现我们要处理的逻辑函数
			r.CallBackFunc(d)
		}
	}()
	log.Printf("[*] waiting for messages, to exit press ctrl+c!!!")
	select {}
}

//PublishSub .. 订阅模式下生产者
func (r *RabbitMQ) PublishSub(message []byte) error {
	//发送消息
	err := r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	return err
}
