package services

import (
	"encoding/json"
	"log"

	"github.com/ersa97/nats-libs/helpers"
	"github.com/ersa97/nats-libs/models"
	"github.com/nats-io/nats.go"
)

// func NatsPublish(connection *models.NatsConnection, data interface{}) (responseStatus int, responseMessage string, responseData interface{}) {
// 	url := connection.Ip + ":" + connection.Port

// 	log.Println("[NATS] " + url + "|" + connection.Subject)

// 	conn, err := nats.Connect("nats://"+connection.Username+":"+connection.Password+"@"+url,
// 		nats.RetryOnFailedConnect(true),
// 		nats.MaxReconnects(10),
// 		nats.ReconnectWait(time.Second),
// 		nats.ReconnectHandler(func(c *nats.Conn) {}))

// 	if err != nil {
// 		return http.StatusBadRequest, "failed to connect to nats", err
// 	}
// 	defer conn.Close()

// 	err = conn.Publish(connection.Subject, []byte(helpers.JSONEncode(data)))

// 	if err != nil {
// 		return http.StatusBadRequest, "failed to publish message to nats", err
// 	}

// 	return http.StatusOK, "Message => successfully to publish a message to nats.", nil

// }

// func NatsSubscribe(connection models.NatsConnection, conn *nats.Conn) (*nats.Msg, error) {

// 	sub, err := conn.SubscribeSync(connection.Subject)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// defer sub.Unsubscribe()

// 	msg, err := sub.NextMsg(10 * time.Second)

// 	if err != nil {
// 		return nil, err
// 	}

// 	// log.Println("data", msg.Data)

// 	return msg, nil
// }

// func StartConnection(connection models.NatsConnection) (*nats.Conn, error) {
// 	url := connection.Ip + ":" + connection.Port

// 	log.Println("[NATS] " + url + "|" + connection.Subject)

// 	conn, err := nats.Connect("nats://"+connection.Username+":"+connection.Password+"@"+url, nats.Timeout(10*time.Second))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return conn, nil
// }

func PublishMessage(connection models.NatsConnection, streamName, streamSubjects, subjectName, Command string, Param, Data interface{}) error {

	url := connection.Ip + ":" + connection.Port

	log.Println("[NATS] " + url + "|" + connection.Subject)

	nc, err := nats.Connect(url)
	if err != nil {
		helpers.ErrorMessage(err.Error(), err)
		return err
	}

	js, err := nc.JetStream()
	if err != nil {
		helpers.ErrorMessage(err.Error(), err)
		return err
	}

	ssubject := models.StreamSubject{
		StreamName:     streamName,
		StreamSubjects: streamSubjects,
		SubjectName:    streamName,
	}

	err = createStream(js, ssubject)
	if err != nil {
		helpers.ErrorMessage(err.Error(), err)
		return err
	}

	request := models.Request{
		Command: Command,
		Param:   Param,
		Data:    Data,
	}

	err = createOrder(js, ssubject, request)
	if err != nil {
		helpers.ErrorMessage(err.Error(), err)
		return err
	}

	return nil
}

func createOrder(js nats.JetStreamContext, ssubject models.StreamSubject, request models.Request) error {

	requestJSON, _ := json.Marshal(request)
	_, err := js.Publish(ssubject.SubjectName, requestJSON)
	log.Println("Subject => ", ssubject.SubjectName)
	log.Println("Message => ", request)
	if err != nil {
		return err
	}
	return nil
}

func createStream(js nats.JetStreamContext, ssubject models.StreamSubject) error {
	stream, err := js.StreamInfo(ssubject.StreamName)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", ssubject.StreamName, ssubject.StreamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     ssubject.StreamName,
			Subjects: []string{ssubject.StreamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
