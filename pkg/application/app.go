package application

import (
	"encoding/json"
	"log"
	"notify/pkg/model"
	"notify/pkg/provider"

	"github.com/IBM/sarama"
)

type App struct {
	Config           *model.Config
	producer         sarama.SyncProducer
	TelegramProvider *provider.Telegram
	route            Route
}

type Route map[string][]string

func NewApp(pathToConf string) *App {
	app := &App{}
	app.initConfig(pathToConf)
	app.initProducer()
	// app.initTelegram()
	// app.initRouter()

	return app
}

func (app *App) initConfig(pathToConf string) {
	conf, err := model.NewConfig(pathToConf)
	if err != nil {
		log.Panicf("Failed to init config: %s", err)
	}
	app.Config = conf
}

func (app *App) initTelegram() {
	t, err := provider.NewTelegram(app.Config.Debug, app.Config.Telegram, app.Config.Templates)
	if err != nil {
		log.Panicf("Failed to init telegram: %s", err)
	}
	app.TelegramProvider = t
}

func (app *App) initRouter() {
	app.route = make(Route)
	for template := range app.Config.Telegram.TemplateToChats {
		app.route[template] = append(app.route[template], provider.TELEGRAM)
	}
}

func (app *App) Send(action string, req model.Request) {
	for _, router := range app.route[action] {
		switch router {
		case provider.TELEGRAM:
			app.TelegramProvider.Send(
				action,
				req,
			)
		}
	}
}

func (app *App) initProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{app.Config.Kafka.Address}, config)
	if err != nil {
		log.Panicf("Failed to init producer: %s", err)
	}
	app.producer = producer
}

func (app *App) SendMessage(action string, req model.Request) error {
	notificationJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: app.Config.Kafka.Topic,
		Key:   sarama.StringEncoder(provider.TELEGRAM),
		Value: sarama.StringEncoder(notificationJSON),
	}
	_, _, err = app.producer.SendMessage(msg)

	return err
}

func (app *App) Close() {
	app.producer.Close()
	app.TelegramProvider.Close()
}
