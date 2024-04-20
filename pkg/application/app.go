package application

import (
	"log"
	"notify/pkg/model"
)

type App struct {
	Config *model.Config
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

func (app *App) Log(v ...any) {
	log.Println(v...)
}

func (app *App) Logf(format string, v ...any) {
	log.Printf(format, v...)
}
