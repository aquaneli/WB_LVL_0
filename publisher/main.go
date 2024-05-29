package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	f, err := os.Open("/Users/aquaneli/WB_LVL_0/model.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully Opened /Users/aquaneli/WB_LVL_0/model.json")
	defer f.Close()

	byteValue, _ := io.ReadAll(f)

	ns, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	err = ns.Publish("a", byteValue)
	if err != nil {
		panic(err)
	}
	log.Println("Message published successfully")
}
