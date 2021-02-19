package main

import (
	"flag"
	"log"
	"voip-shovel-go/src/shovel/plugin"
	"voip-shovel-go/src/shovel/version"
)

func main() {
	configPath := flag.String("path", "config/shovel.yaml", "The config path of shovel App")
	flag.Parse()

	log.Println("version: ", version.Full())
	config, err := plugin.LoadConfiguration(*configPath)
	if err != nil {
		panic(err)
	}

	for _, rabbitSrcConf := range config.RabbitMQConfSrc {
		go plugin.NewShovel(&rabbitSrcConf, &config.RabbitMQConfDest, config.HandlerNumbers).Start()
	}
	select{}
}
