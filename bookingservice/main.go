package main

import (
	"flag"

	"github.com/GO_NATIVE/bookingservice/listener"
	"github.com/GO_NATIVE/lib/configuration"
	msgqueue_amqp "github.com/GO_NATIVE/lib/msgqueue/amqp"
	"github.com/GO_NATIVE/lib/persistence/dblayer"
	"github.com/streadway/amqp"
)

func main() {
	confPath := flag.String("config", "./configuration/config.json", "path to config file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)

	dblayer, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		panic(err)
	}

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}

	eventListener, err := msgqueue_amqp.NewAMQPEventListener(conn, "events")
	if err != nil {
		panic(err)
	}

	processor := &listener.EventProcessor{eventListener, dblayer}
	go processor.ProcessEvents()

}
