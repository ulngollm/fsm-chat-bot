package main

import (
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/ulngollm/teleflow"
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
	b, err := NewBot(opts.BotToken)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	pool := teleflow.NewMemoryPool()
	flowManager := teleflow.NewFlowManager(pool)
	router := teleflow.NewFlowRouter(flowManager)

	defaultFlowController := NewDefaultFlowController(flowManager)

	g := router.Group("default")
	g.AddHandler(stateFirst, defaultFlowController.handleFirst)
	g.AddHandler(stateSecond, defaultFlowController.handleSecond)
	g.AddHandler(stateThird, defaultFlowController.handlerThird)
	g.AddHandler(stateLast, defaultFlowController.handleLast)
	g.AddHandler(stateClosed, defaultFlowController.handleClose)

	b.Handle(tele.OnText, defaultFlowController.handleInit, router.Middleware())

	b.Start()
	return nil
}
