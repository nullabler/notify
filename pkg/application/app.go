package application

import (
	"log"
	"notify/pkg/model"
	"notify/pkg/provider"
)

type App struct {
	Config           *model.Config
	TelegramProvider *provider.Telegram
	route            Route
}

type Route map[string][]string

func NewApp(pathToConf string) *App {
	app := &App{}
	app.initConfig(pathToConf)
	app.initTelegram()
	app.initRouter()

	return app
}

func (app *App) initConfig(pathToConf string) {
	conf, err := model.NewConfig(pathToConf)
	if err != nil {
		log.Panic(err)
	}
	app.Config = conf
}

func (app *App) initTelegram() {
	t, err := provider.NewTelegram(app.Config.Debug, app.Config.Telegram, app.Config.Templates)
	if err != nil {
		log.Panic(err)
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
