package plugin

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"voip-shovel-go/src/shovel/rabbitmq"
)

type RabbitMQConf struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	VirtualHost string `yaml:"virtualHost"`

	Exchange struct {
		AutoDelete   bool   `yaml:"autoDelete"`
		Durable      bool   `yaml:"durable"`
		ExchangeName string `yaml:"exchangeName"`
		ExchangeType string `yaml:"exchangeType"`
	}

	Queue struct {
		AutoDelete bool                   `yaml:"autoDelete"`
		Durable    bool                   `yaml:"durable"`
		Exclusive  bool                   `yaml:"exclusive"`
		QueueName  string                 `yaml:"queueName"`
		RoutingKey string                 `yaml:"routingKey"`
		Arguments  map[string]interface{} `yaml:"arguments,omitempty"`
	}
}

type Config struct {
	HandlerNumbers   int            `yaml:"handlerNumbers"`
	RabbitMQConfSrc  []RabbitMQConf `yaml:"rabbitMQSrc"`
	RabbitMQConfDest RabbitMQConf   `yaml:"rabbitMQDest"`
}

func LoadConfiguration(configFilePath string) (config *Config, err error) {
	configBuf, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		err = fmt.Errorf("error loading configuration file %s: %v", configFilePath, err)
	} else {
		config = new(Config)
		if err = yaml.Unmarshal(configBuf, &config); err != nil {
			err = fmt.Errorf("error parsing configuration file %s: %v", configFilePath, err)
		}
	}
	return config, err
}

type Shovel struct {
	cons *rabbitmq.Consumer
	pub  *rabbitmq.Publisher
}

func NewShovel(
	srcConf *RabbitMQConf,
	destConf *RabbitMQConf,
	handlerNumbers int,
) *Shovel {

	addrPub := fmt.Sprintf(
		"amqp://%v:%v@%v:%v%v",
		destConf.Username,
		destConf.Password,
		destConf.Host,
		destConf.Port,
		destConf.VirtualHost,
	)

	addrCons := fmt.Sprintf(
		"amqp://%v:%v@%v:%v%v",
		srcConf.Username,
		srcConf.Password,
		srcConf.Host,
		srcConf.Port,
		srcConf.VirtualHost,
	)

	c := &Shovel{
		pub: rabbitmq.NewPublisher(addrPub, destConf.Exchange.ExchangeName, destConf.Exchange.ExchangeType),
	}
	cons := rabbitmq.NewConsumer(
		addrCons,
		srcConf.Exchange.ExchangeName,
		srcConf.Exchange.ExchangeType,
		srcConf.Queue.QueueName,
		srcConf.Queue.RoutingKey,
		srcConf.Queue.Durable,
		srcConf.Queue.AutoDelete,
		srcConf.Queue.Arguments,
		c.shovel,
		handlerNumbers)

	c.cons = cons
	return c
}

func (s *Shovel) shovel(msg []byte) error {
	s.pub.Send(string(msg), "")
	return nil
}

func (s *Shovel) Start () {
	s.cons.Start()
}
