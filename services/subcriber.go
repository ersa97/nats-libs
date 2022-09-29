package services

import (
	"fmt"
	"log"

	"github.com/ersa97/nats-libs/helpers"
	"github.com/ersa97/nats-libs/models"
	"github.com/nats-io/nats.go"
)

func SubscribeNats(connection models.NatsConnection) (data *nats.Msg, err error) {
	log.Println("connect to NATS")

	conn, err := StartConnection(connection)
	if err != nil {
		helpers.ErrorMessage("failed to connect NATS", err)
	}

	data, err = NatsSubscribe(connection, conn)
	if err != nil {
		helpers.ErrorMessage("failed to subscribe NATS", err)
		return
	}

	fmt.Println("data -> ", data.Data)

	return

}