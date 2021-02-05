package main

import (
	"flag"
)

func main() {
	configPath := flag.String("path", "config/shovel.yaml", "The config path of shovel App")
	flag.Parse()

	config, err := LoadConfiguration(*configPath)
	if err != nil {
		panic(err)
	}

	for _, rabbitSrcConf := range config.RabbitMQConfSrc {
		go NewShovel(&rabbitSrcConf, &config.RabbitMQConfDest, config.HandlerNumbers).Start()
	}
	select{}
}
