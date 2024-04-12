package provider

import (
	"notify/pkg/model"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

const (
	TELEGRAM = "telegram"
)

type Telegram struct {
	conf      model.TelegramConf
	templates map[string]string

	bot *telego.Bot
	bh  *th.BotHandler
}

func NewTelegram(debug bool, conf model.TelegramConf, templates map[string]string) (*Telegram, error) {
	t := &Telegram{
		conf:      conf,
		templates: templates,
	}

	botOption := telego.WithDiscardLogger()
	if debug {
		botOption = telego.WithDefaultDebugLogger()
	}

	bot, err := telego.NewBot(t.conf.Token, botOption)
	if err != nil {
		return t, err
	}
	t.bot = bot
	go t.cmd()

	return t, nil
}

func (t *Telegram) Send(action string, req model.Request) {
	message, ok := t.getMessage(action, req)
	if !ok {
		return
	}

	for _, chatId := range t.conf.TemplateToChats[action] {
		msg := tu.Message(
			tu.ID(chatId),
			message,
		).WithProtectContent()

		t.bot.SendMessage(msg)
	}
}

func (t *Telegram) getMessage(action string, req model.Request) (string, bool) {
	req = t.applyAliases(req)
	msg := t.templates[action]
	for key, val := range req {
		msg = strings.ReplaceAll(msg, "{{"+key+"}}", val)
	}

	return msg, true
}

func (t *Telegram) applyAliases(req model.Request) model.Request {
	result := make(model.Request)
	for key, val := range req {
		if alias, ok := t.conf.Aliases[key][val]; ok {
			result[key] = alias
			continue
		}
		result[key] = val
	}

	return result
}

func (t *Telegram) cmd() {
	updates, _ := t.bot.UpdatesViaLongPolling(nil)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msg := ""

		switch update.Message.Text {
		case t.conf.Trigger + " ChatID":
			msg = "ChatID: " + strconv.Itoa(int(chatID))
		case t.conf.Trigger + " Ping":
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

func (t *Telegram) Clear() {
	t.bot.StopLongPolling()
}
