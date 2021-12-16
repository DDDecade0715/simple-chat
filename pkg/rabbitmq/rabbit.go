package rabbitmq

import (
	"fmt"
	"gin-derived/config"
	"github.com/streadway/amqp"
)

//RabbitMQ 结构体
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
	//
	Error chan error
}

var RabbitmqPool = make(map[string]*RabbitMQ)

// GetRabbitMQ 获取连接体
func GetRabbitMQ(name string, config config.Params, url string) *RabbitMQ {
	if c, ok := RabbitmqPool[name]; ok {
		return c
	}
	c := CreateRabbitMQ(config.Exchange, config.Queues, url, config.Key)

	RabbitmqPool[name] = c

	return c
}

//CreateRabbitMQ 创建Rabbitmq
func CreateRabbitMQ(exchange string, queues string, url string, key string) *RabbitMQ {
	//实例
	c := &RabbitMQ{
		Exchange:  exchange,
		QueueName: queues,
		Mqurl:     url,
		Key:       key,
	}

	if err := c.ConnectRabbitMQ(); err != nil {
		c.failOnErr(err, "failed to connect rabbitmq!")
	}

	return c
}

func (r *RabbitMQ) ConnectRabbitMQ() error {
	var err error
	r.conn, err = amqp.Dial(r.Mqurl)
	r.failOnErr(err, "failed to connect rabbitmq!")
	r.channel, err = r.conn.Channel()
	r.failOnErr(err, "failed to open a channel")

	return nil
}

//Destory 断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		//global.GVA_LOG.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//BindQueue 绑定队列
func (r *RabbitMQ) BindQueue() error {
	if _, err := r.channel.QueueDeclare(r.QueueName, true, false, false, false, nil); err != nil {
		return fmt.Errorf("error in declaring the queue %s", err)
	}
	if err := r.channel.QueueBind(r.QueueName, "my_routing_key", r.Exchange, false, nil); err != nil {
		return fmt.Errorf("Queue  Bind error: %s", err)
	}

	return nil
}

//Reconnect 重连
func (r *RabbitMQ) Reconnect() error {
	if err := r.ConnectRabbitMQ(); err != nil {
		return err
	}
	if err := r.BindQueue(); err != nil {
		return err
	}
	return nil
}

//PublishSimple 直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string) {
	select {
	case err := <-r.Error:
		if err != nil {
			r.Reconnect()
		}
	default:
	}
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
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
		fmt.Println(err)
	}
	//调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//ConsumeSimple simple模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
	select {
	case err := <-r.Error:
		if err != nil {
			r.Reconnect()
		}
	default:
	}
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
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
		fmt.Println(err)
	}

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
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

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			fmt.Printf("\n Received a message: %s \n", d.Body)
		}
	}()
	fmt.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
