package bot

import (
	"time"

	"github.com/ulngollm/msg-constructor/pkg/middleware"

	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	bot *tele.Bot
}

func New(token string) (*Bot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &Bot{bot: bot}, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) RegisterHandler(endpoint string, baseHandler tele.HandlerFunc, m middleware.Middleware) {
	b.bot.Handle(endpoint, baseHandler, m.Handle)
}
