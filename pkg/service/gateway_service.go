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

func (g *GatewaySvc) SendMessage(action string, req model.Request) error {
	notifyJSON, err := json.Marshal(model.NewNotify(action, req))
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: g.app.Config.Kafka.Topic,
		Key:   sarama.StringEncoder(provider.TELEGRAM_UID),
		Value: sarama.StringEncoder(notifyJSON),
	}

	_, _, err = g.producer.SendMessage(msg)

	return err
}

func (g *GatewaySvc) Close() error {
	return g.producer.Close()
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
