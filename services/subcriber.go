package services

import (
	"context"
	"log"

	"github.com/ersa97/nats-libs/helpers"
	"github.com/ersa97/nats-libs/models"
	"github.com/nats-io/nats.go"
)

func SubscribeNats(connection models.NatsConnection, subjectname string) (sub *nats.Subscription, ctx context.Context, err error) {
	url := connection.Ip + ":" + connection.Port

	log.Println("[NATS] " + url)

	nc, err := nats.Connect(url)
	if err != nil {
		helpers.ErrorMessage(err.Error(), err)
		return nil, nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		helpers.ErrorMessage(err.Error(), err)
		return nil, nil, err
	}

	subSubjectName := subjectname

	sub, _ = js.PullSubscribe(subSubjectName, "muunship")

	ctx, _ = context.WithCancel(context.Background())

	return

}
