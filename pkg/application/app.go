package application

import (
	"log"
	"notify/pkg/model"

	"github.com/IBM/sarama"
)

type App struct {
	Config   *model.Config
	Producer sarama.SyncProducer
}

func NewApp(pathToConf string) *App {
	app := &App{}
	app.initConfig(pathToConf)

	return app
}

func (app *App) initConfig(pathToConf string) {
	conf, err := model.NewConfig(pathToConf)
	if err != nil {
		log.Panicf("Failed to init config: %s", err)
	}
	app.Config = conf
}

func (app *App) Dump(msg interface{}) {
	if app.Config.Debug {
		log.Printf("Dump: %v", msg)
	}
}

func (app *App) Logf(format string, v ...any) {
	log.Printf(format, v...)
}