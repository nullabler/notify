package service

import (
	"context"
	"encoding/json"
	"notify/pkg/application"
	"notify/pkg/model"
	"notify/pkg/provider"

	"github.com/IBM/sarama"
)

const CONSUMER_GROUP_ID = "notify-group"

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

func (c *ConsumerSvc) Invoke() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		if err := c.consumerGroup.Consume(ctx, []string{c.app.Config.Kafka.Topic}, c); err != nil {
			c.app.Logf("Error from consumer: %v", err)
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
		notify := model.Notify{}
		err := json.Unmarshal(msg.Value, &notify)
		if err != nil {
			c.app.Logf("Failed to unmarshal notification: %v", err)
			continue
		}

		c.app.Logf("Got a message: %v", notify)
		switch string(msg.Key) {
		case provider.TELEGRAM_UID:
			c.telegramProvider.Send(notify)
			break
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (c *ConsumerSvc) initCounsumerGroup() error {
	config := sarama.NewConfig()
	if c.app.Config.Debug {
		config.Consumer.Return.Errors = true
	}

	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{c.app.Config.Kafka.Address},
		CONSUMER_GROUP_ID,
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
	if !c.telegramProvider.IsEnabled() {
		c.app.Log("Telegram disabled")
		return nil
	}

	go c.telegramProvider.CmdHandler()

	return nil
}
