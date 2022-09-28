package services

import (
	"log"
	"net/http"
	"time"

	"github.com/ersa97/nats-libs/helpers"
	"github.com/ersa97/nats-libs/models"
	"github.com/nats-io/nats.go"
)

func NatsPublish(connection *models.NatsConnection, data interface{}) (responseStatus int, responseMessage string, responseData interface{}) {
	url := connection.Ip + ":" + connection.Port

	log.Println("[NATS] " + url + "|" + connection.Subject)

	conn, err := nats.Connect("nats://"+connection.Username+":"+connection.Password+"@"+url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.ReconnectHandler(func(c *nats.Conn) {}))

	if err != nil {
		return http.StatusBadRequest, "failed to connect to nats", err
	}
	defer conn.Close()

	err = conn.Publish(connection.Subject, []byte(helpers.JSONEncode(data)))

	if err != nil {
		return http.StatusBadRequest, "failed to publish message to nats", err
	}

	return http.StatusOK, "Message => successfully to publish a message to nats.", nil

	// err = conn.PublishMsg(&nats.Msg{
	// 	Subject: connection.Subject,
	// 	Data:    []byte(helpers.JSONEncode(data)),
	// 	Sub: &nats.Subscription{
	// 		Subject: connection.Subject,
	// 	},
	// })

	// if err != nil {
	// 	return http.StatusBadRequest, "failed to publish message to nats", err
	// }

	// return http.StatusOK, "Message => successfully to publish a message to nats.", nil

}

func NatsSubscribe(connection models.NatsConnection) (data interface{}, err error) {
	url := connection.Ip + ":" + connection.Port

	log.Println("[NATS] " + url + "|" + connection.Subject)

	conn, err := nats.Connect("nats://" + connection.Username + ":" + connection.Password + "@" + url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sub, err := conn.SubscribeSync(connection.Subject)
	defer sub.Unsubscribe()

	msg, err := sub.NextMsg(10 * time.Second)

	if err != nil {
		return nil, err
	}

	// log.Println("data", msg.Data)

	return msg.Data, nil
}
