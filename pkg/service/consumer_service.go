package service

import (
	"context"
	"encoding/json"
	"log"
	"notify/pkg/application"
	"notify/pkg/model"
	"notify/pkg/provider"

	"github.com/IBM/sarama"
)

type ConsumerSvc struct {
	app              *application.App
	telegramProvider *provider.TelegramProvider

	consumerGroup sarama.ConsumerGroup
}

func NewConsumerSvc(app *application.App) (*ConsumerSvc, error) {
	c := &ConsumerSvc{
		app: app,
	}

	if err := c.initCounsumerGroup(); err != nil {
		return c, err
	}

	if err := c.initTelegramProvider(); err != nil {
		return c, err
	}

	return c, nil
}

func (c *ConsumerSvc) initCounsumerGroup() error {
	config := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{c.app.Config.Kafka.Address},
		"notify-group",
		config,
	)
	if err != nil {
		return err
	}
	c.consumerGroup = consumerGroup

	return nil
}

func (c *ConsumerSvc) initTelegramProvider() error {
	telegramProvider, err := provider.NewTelegramProvider(c.app)
	if err != nil {
		return err
	}
	c.telegramProvider = telegramProvider
	go c.telegramProvider.CmdHandler()

	return nil
}

func (c *ConsumerSvc) Invoice() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		if err := c.consumerGroup.Consume(ctx, []string{c.app.Config.Kafka.Topic}, c); err != nil {
			log.Printf("error from consumer: %v", err)
		}

		if ctx.Err() != nil {
			return
		}
	}
}

func (c *ConsumerSvc) Close() {
	c.consumerGroup.Close()
	c.telegramProvider.Close()
}

func (*ConsumerSvc) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (*ConsumerSvc) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerSvc) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.app.Logf("Start consume claim")
	for msg := range claim.Messages() {
		req := model.Request{}
		err := json.Unmarshal(msg.Value, &req)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}
		c.app.Logf("Got a message: %v", req)
		c.telegramProvider.Send("pipeline-stage", req)
		sess.MarkMessage(msg, "")
	}
	return nil
}
