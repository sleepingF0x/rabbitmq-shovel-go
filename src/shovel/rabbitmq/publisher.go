package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Publisher struct {
	conn     *Connection
	exchange string
}

func NewPublisher(
	addr string,
	exchange string,
	exchangeType string) *Publisher {

	c := &Publisher{
		conn: NewConnection(addr, exchange, exchangeType, true),
	}

	_ = c.conn.Connect()
	return c
}

func (p *Publisher) Send(msg string, routingKey string) {
	for {
	clear:
		for {
			select {
			case <-p.conn.connected:
			default:
				break clear
			}
		}
		if err := p.conn.channel.Publish(
			p.conn.exchange, // exchange
			routingKey,      // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType:  "text/plain",
				DeliveryMode: amqp.Persistent,
				Body:         []byte(msg),
			}); err != nil {
			log.Println("rabbitmq publish - failCheck: ", err)
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}
}
