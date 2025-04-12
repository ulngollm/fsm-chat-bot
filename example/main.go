package main

import (
	"log"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/ulngollm/teleflow"
	_default "github.com/ulngollm/teleflow/example/flow/default"
	tele "gopkg.in/telebot.v4"
)

type options struct {
	BotToken string `long:"token" env:"BOT_TOKEN" required:"true" description:"telegram bot token"`
}

func main() {
	var opts options
	p := flags.NewParser(&opts, flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		log.Printf("parse: %s", err)
		return
	}

	if err := run(opts); err != nil {
		log.Printf("run: %s", err)
	}
}

func run(opts options) error {
	pref := tele.Settings{
		Token:  opts.BotToken,
		Poller: &tele.LongPoller{Timeout: time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	pool := teleflow.NewMemoryPool()
	flowManager := teleflow.NewFlowManager(pool)
	router := teleflow.NewFlowRouter(flowManager)

	defaultFlowController := _default.NewDefaultFlowController(flowManager)

	g := router.Group("default")
	g.AddHandler(_default.StateFirst, defaultFlowController.HandleFirst)
	g.AddHandler(_default.StateSecond, defaultFlowController.HandleSecond)
	g.AddHandler(_default.StateThird, defaultFlowController.HandlerThird)
	g.AddHandler(_default.StateLast, defaultFlowController.HandleLast)
	g.AddHandler(_default.StateClosed, defaultFlowController.HandleClose)

	bot.Handle(tele.OnText, defaultFlowController.HandleInit, router.Middleware())

	bot.Start()
	return nil
}
