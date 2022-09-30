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

}

func NatsSubscribe(connection models.NatsConnection, conn *nats.Conn) (*nats.Msg, error) {

	ch := make(chan *nats.Msg, 64)
	sub, err := conn.ChanSubscribe(connection.Subject, ch)

	var msgCh *nats.Msg
	for msg := range ch {
		msgCh = msg
	}
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()

	// msg, err := sub.NextMsg(10 * time.Second)

	if err != nil {
		return nil, err
	}

	// log.Println("data", msg.Data)

	return msgCh, nil
}

func StartConnection(connection models.NatsConnection) (*nats.Conn, error) {
	url := connection.Ip + ":" + connection.Port

	log.Println("[NATS] " + url + "|" + connection.Subject)

	conn, err := nats.Connect("nats://"+connection.Username+":"+connection.Password+"@"+url, nats.Timeout(10*time.Second))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
