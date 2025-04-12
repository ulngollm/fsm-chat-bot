package bot

import (
	"time"

	tele "gopkg.in/telebot.v4"
)

// todo убрать эту обертку. Это никто не сможет использовать
type Bot struct {
	bot *tele.Bot
}

func New(token string) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return bot, nil
}
