package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/GO_NATIVE/eventsservice/rest"
	"github.com/GO_NATIVE/lib/configuration"
	msgqueue_amqp "github.com/GO_NATIVE/lib/msgqueue/amqp"
	"github.com/GO_NATIVE/lib/persistence/dblayer"
	"github.com/streadway/amqp"
)

func main() {

	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}

	emitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to database")
	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		fmt.Println(err)
	}
	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndpoint, dbhandler, emitter)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error:", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error:", err)
	}
}
