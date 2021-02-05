# go-rabbitMQ-shovel

![](https://w0w-public-pics.oss-cn-shenzhen.aliyuncs.com/github-MD/rabbitmq-shovel-go.png)
### 功能说明
Shovel原本是RabbitMQ的一个插件，插件的功能就是将源节点的消息发布到目标节点。本程序实现了同样的功能，负责从一个节点读取数据然后发送到目标节点。

### 使用方法
shovel -path /root/yourConfigFile.yml


### 配置文件说明

| 字段                               | 说明                                         |
| ---------------------------------- | -------------------------------------------- |
| rabbitMQSrc                           | 源MQ连接参数                                |
| rabbitMQDest                          | 目标MQ连接参数                              |
| handlerNumbers                   | 源MQ消费者数量                              |
