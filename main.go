package main

import (
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/ulngollm/msg-constructor/internal/bot"
	lflow "github.com/ulngollm/msg-constructor/internal/flow"
	"github.com/ulngollm/msg-constructor/internal/middleware"
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
	b, err := bot.NewBot(opts.BotToken)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	// todo посмотреть, как в норм проектах инициализируют такие наборы сложных вложенных структур
	pool := lflow.NewPool()
	flowManager := lflow.New(pool)
	router := middleware.NewFlowFinder(flowManager)

	defaultFlowController := NewDefaultFlowController(flowManager)

	g := router.Group("default")
	g.AddHandler(stateFirst, defaultFlowController.handleFirst)
	g.AddHandler(stateSecond, defaultFlowController.handleSecond)
	g.AddHandler(stateThird, defaultFlowController.handlerThird)
	g.AddHandler(stateLast, defaultFlowController.handleLast)
	g.AddHandler(stateClosed, defaultFlowController.handleClose)

	// see flow может быть инициализирован из разных мест. Например, начаться с команды или с сообщения
	b.Handle(tele.OnText, defaultFlowController.handleInit, router.Handle)

	b.Start()
	return nil
}
