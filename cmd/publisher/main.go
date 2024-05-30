package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	flag.Args()
	f, err := os.Open("/Users/aquaneli/WB_LVL_0/model2.json")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully Opened /Users/aquaneli/WB_LVL_0/model2.json")
	defer f.Close()

	byteValue, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	ns, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}

	err = ns.Publish("orders", byteValue)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Message published successfully")
}
