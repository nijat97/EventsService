package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/GO_NATIVE/eventsservice/rest"
	"github.com/GO_NATIVE/lib/configuration"
	"github.com/GO_NATIVE/lib/persistence/dblayer"
)

func main() {
	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()

	config, _ := configuration.ExtractConfiguration(*confPath)
	fmt.Println("Connecting to database")
	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		fmt.Println(err)
	}
	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndpoint, dbhandler)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error:", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error:", err)
	}
}
