package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	f, err := os.Open("/Users/aquaneli/go_lerning/model.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully Opened /Users/aquaneli/go_lerning/model.json")
	defer f.Close()

	byteValue, _ := io.ReadAll(f)

	ns, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	_ = ns.Publish("a", byteValue)

	time.Sleep(1 * time.Hour)
}
