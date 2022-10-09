package natslibs

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

type natslibs struct {
	ctx        context.Context
	conn       *nats.Conn
	streamName string
}

//create new connection and stream
//return context, *nats.Conn, and streamname
func New(ctx context.Context, host, streamName string) *natslibs {
	nc, err := nats.Connect(host)
	if err != nil {
		log.Fatal("failed to connect nats server")
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	stream, err := js.StreamInfo(streamName)
	if stream == nil || err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamName + ".*"},
		})
		if err != nil {
			log.Fatal("failed to connect nats server")
		}
	}
	return &natslibs{
		ctx:        ctx,
		conn:       nc,
		streamName: streamName,
	}
}

//publish message as payload thru subject
//return error if exist
func (c *natslibs) Publish(subject string, payload map[string]interface{}) error {
	js, err := c.conn.JetStream()
	if err != nil {
		return err
	}
	requestJSON, _ := json.Marshal(payload)
	_, err = js.Publish(c.streamName+"."+subject, requestJSON)
	return err
}

//subcribe to subject and wait
//if there is a message put it into msg
func (c *natslibs) Subscribe(subject string, msg chan []byte) error {
	js, err := c.conn.JetStream()
	if err != nil {
		return err
	}
	sub, _ := js.PullSubscribe(c.streamName+"."+subject, "muunship")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
			}

			msgs, _ := sub.Fetch(128 * 1024)
			for _, m := range msgs {
				msg <- m.Data
				m.Ack()
			}
		}
	}()
	return nil
}
