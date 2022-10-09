package main

import (
	"context"
	"fmt"
	"log"

	natslibs "github.com/ersa97/nats-libs"
	"github.com/nats-io/nats.go"
)

func main() {
	ctx, done := context.WithCancel(context.Background())
	defer done()
	n := natslibs.New(ctx, nats.DefaultURL, "test")

	msg := make(chan []byte)

	err := n.Subscribe("simple", msg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			m := <-msg
			fmt.Printf("Get Data -> %s\n", string(m))
		}
	}()

	fmt.Println("Running Subscriber")
	fmt.Scanf("input")
	fmt.Println("Shutdown")

}
