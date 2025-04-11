package bot

import (
	"time"

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

func (b *Bot) RegisterHandlers() {
	b.bot.Handle("/start", b.handleStart)
}

func (b *Bot) handleStart(c tele.Context) error {
	return c.Send(c.Message())
}
