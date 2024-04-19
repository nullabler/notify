package provider

import (
	"notify/pkg/application"
	"notify/pkg/model"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

const TELEGRAM_UID = "telegram"

type TelegramProvider struct {
	app *application.App

	bot *telego.Bot
	bh  *th.BotHandler
}

func NewTelegramProvider(app *application.App) (*TelegramProvider, error) {
	t := &TelegramProvider{
		app: app,
	}

	err := t.initTelego()
	if err != nil {
		return t, err
	}

	return t, nil
}

func (t *TelegramProvider) initTelego() error {
	botOption := telego.WithDiscardLogger()
	if t.app.Config.Debug {
		botOption = telego.WithDefaultDebugLogger()
	}

	bot, err := telego.NewBot(t.app.Config.Telegram.Token, botOption)
	if err != nil {
		return err
	}
	t.bot = bot

	return nil
}

func (t *TelegramProvider) Send(notify model.Notify) {
	message, ok := t.getMessage(notify.Action, notify.Body)
	if !ok {
		return
	}

	for _, chatId := range t.app.Config.Telegram.TemplateToChats[notify.Action] {
		msg := tu.Message(
			tu.ID(chatId),
			message,
		).WithProtectContent()

		t.bot.SendMessage(msg)
	}
}

func (t *TelegramProvider) getMessage(action string, req model.Request) (string, bool) {
	req = t.applyAliases(req)
	msg := t.app.Config.Templates[action]
	for key, val := range req {
		msg = strings.ReplaceAll(msg, "{{"+key+"}}", val)
	}

	return msg, true
}

func (t *TelegramProvider) applyAliases(req model.Request) model.Request {
	result := make(model.Request)
	for key, val := range req {
		if alias, ok := t.app.Config.Telegram.Aliases[key][val]; ok {
			result[key] = alias
			continue
		}
		result[key] = val
	}

	return result
}

func (t *TelegramProvider) CmdHandler() {
	updates, _ := t.bot.UpdatesViaLongPolling(nil)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msg := ""

		switch update.Message.Text {
		case t.app.Config.Telegram.Trigger + " ChatID":
			msg = "ChatID: " + strconv.Itoa(int(chatID))
		case t.app.Config.Telegram.Trigger + " Ping":
			msg = "Pong"
		}

		if msg != "" {
			t.bot.SendMessage(
				tu.Message(
					tu.ID(chatID),
					msg,
				),
			)
		}
	}
}

func (t *TelegramProvider) Close() {
	t.bot.StopLongPolling()
}
