package service

import (
	"encoding/json"
	"notify/pkg/application"
	"notify/pkg/model"
	"notify/pkg/provider"

	"github.com/IBM/sarama"
)

type GatewaySvc struct {
	app *application.App

	producer sarama.SyncProducer
}

func NewGatewaySvc(app *application.App) (*GatewaySvc, error) {
	svc := &GatewaySvc{
		app: app,
	}
	err := svc.initProducer()

	return svc, err
}

func (g *GatewaySvc) initProducer() error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{g.app.Config.Kafka.Address}, config)
	if err != nil {
		return err
	}
	g.producer = producer

	return nil
}

func (g *GatewaySvc) SendMessage(action string, req model.Request) error {
	notificationJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: g.app.Config.Kafka.Topic,
		Key:   sarama.StringEncoder(provider.UID_TELEGRAM),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = g.producer.SendMessage(msg)

	return err
}

func (g *GatewaySvc) Close() error {
	return g.producer.Close()
}
