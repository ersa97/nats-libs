package services

import (
	"encoding/json"
	"log"

	"github.com/ersa97/nats-libs/helpers"
	"github.com/ersa97/nats-libs/models"
	"github.com/nats-io/nats.go"
)

func PublishMessage(connection models.NatsConnection, streamName, streamSubjects, subjectName, Command string, Param, Data interface{}) error {

	url := connection.Ip + ":" + connection.Port

	log.Println("[NATS] " + url)

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
		SubjectName:    subjectName,
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

	err = Publish(js, ssubject, request)
	if err != nil {
		helpers.ErrorMessage(err.Error(), err)
		return err
	}

	return nil
}

func Publish(js nats.JetStreamContext, ssubject models.StreamSubject, request models.Request) error {

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
