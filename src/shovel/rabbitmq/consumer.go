package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct {
	conn       *Connection
	queue      string
	routingKey string
	autoDelete bool
	durable    bool
	args       amqp.Table
	handler    func([]byte) error
	handlerNum int
}

func NewConsumer(
	addr string,
	exchange string,
	exchangeType string,
	queue string,
	routingKey string,
	durable bool,
	autoDelete bool,
	args map[string]interface{},
	handler func([]byte) error,
	handlerNum int) *Consumer {

	c := &Consumer{
		conn:       NewConnection(addr, exchange, exchangeType, true),
		queue:      queue,
		routingKey: routingKey,
		autoDelete: autoDelete,
		durable:    durable,
		handler:    handler,
		args:       args,
		handlerNum: handlerNum,
	}
	return c
}

func (c *Consumer) Start() error {
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Consumer) Stop() {
	close(c.conn.quit)

	if !c.conn.conn.IsClosed() {
		// 关闭 SubMsg message delivery
		if err := c.conn.channel.Cancel("", true); err != nil {
			log.Println("rabbitmq consumer - channel cancel failed: ", err)
		}

		if err := c.conn.conn.Close(); err != nil {
			log.Println("rabbitmq consumer - connection close failed: ", err)
		}
	}
}

func (c *Consumer) Run() error {
	var err error

	if err = c.conn.Connect(); err != nil {
		return err
	}

	if _, err = c.conn.channel.QueueDeclare(
		c.queue,      // name
		true,         // durable
		c.autoDelete, // delete when unused
		false,        // exclusive
		false,        // no-wait
		c.args,       // arguments
	); err != nil {
		log.Printf("queue declare error: %s", err)
		_ = c.conn.channel.Close()
		_ = c.conn.conn.Close()
		return err
	}

	if err = c.conn.channel.QueueBind(
		c.queue,
		c.routingKey,
		c.conn.exchange,
		false,
		nil,
	); err != nil {
		log.Printf("queue bind error: %s", err)
		_ = c.conn.channel.Close()
		_ = c.conn.conn.Close()
		return err
	}

	if err = c.conn.channel.Qos(
		c.handlerNum,
		0,
		false,
	); err != nil {
		log.Printf("channel Qos error: %s", err)
		_ = c.conn.channel.Close()
		_ = c.conn.conn.Close()
		return err
	}

	go c.consumeOnConnect()
	return err
}

func (c *Consumer) Handle(delivery <-chan amqp.Delivery) {
	for d := range delivery {
		go func(delivery amqp.Delivery) {
			if err := c.handler(delivery.Body); err == nil {
				_ = delivery.Ack(false)
			} else {
				// 重新入队，否则未确认的消息会持续占用内存
				_ = delivery.Reject(true)
			}
		}(d)
	}
}

func (c *Consumer) consumeOnConnect() error {
	var err error
	for {
		<-c.conn.connected
		for i := 0; i < c.handlerNum; i++ {
			var delivery <-chan amqp.Delivery
			if delivery, err = c.conn.channel.Consume(
				c.queue, // queue
				"",      // consumer
				false,   // auto-ack
				false,   // exclusive
				false,   // no-local
				false,   // no-wait
				nil,     // args
			); err != nil {
				if !c.conn.conn.IsClosed(){
					log.Printf("channel consume error: %s", err)
					_ = c.conn.channel.Close()
					_ = c.conn.conn.Close()
				}
			} else {
				go c.Handle(delivery)
			}
		}
	}
}
