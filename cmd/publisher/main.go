package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	var path string
	fmt.Println("Enter the full path of the json file")
	fmt.Scan(&path)
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Successfully Opened %s\n", path)
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
	log.Println("The message has been sent")
}
