package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/looplab/fsm"
	"github.com/ulngollm/msg-constructor/internal/bot"
	lflow "github.com/ulngollm/msg-constructor/internal/flow"
	"github.com/ulngollm/msg-constructor/internal/middleware"
	tele "gopkg.in/telebot.v4"
)

type options struct {
	BotToken string `long:"token" env:"BOT_TOKEN" required:"true" description:"telegram bot token"`
}

const flowDefault string = "default"

type App struct {
	DefaultFlowController *DefaultFlowController
}

func NewApp(defaultFlowController *DefaultFlowController) *App {
	return &App{DefaultFlowController: defaultFlowController}
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
	b, err := bot.New(opts.BotToken)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	// сложная инициализация. Нужно 5 струкрур инициализировать, чтобы все заработало
	// todo упростить или инкапсулировать
	// может просто вынести билдинг flowв другое место
	// todo разобрать, как в норм проектах инициализируют такие наборы сложных вложенных структур

	pool := lflow.NewPool()
	flowManager := lflow.New(pool)

	defaultFlowController := NewDefaultFlowController(flowManager)
	flowFinder := middleware.NewFlowFinder(flowManager)

	//flowFinder.Group("flowDefault")
	flowFinder.AddHandler(flowDefault, stateFirst, defaultFlowController.handleFirst)
	flowFinder.AddHandler(flowDefault, stateSecond, defaultFlowController.handleSecond)
	flowFinder.AddHandler(flowDefault, stateThird, defaultFlowController.handlerThird)
	flowFinder.AddHandler(flowDefault, stateLast, defaultFlowController.handleLast)
	flowFinder.AddHandler(flowDefault, stateClosed, defaultFlowController.handleClose)

	// see flow может быть инициализирован из разных мест. Например, начаться с команды или с сообщения
	//b.Group() todo может это использовать?
	b.Handle(tele.OnText, func(c tele.Context) error {
		id := c.Sender().ID
		flow, err := flowManager.InitFlow(id, stateFirst, "default")
		if err != nil {
			return err
		}
		c.Set("flow", flow)
		return nil
	}, flowFinder.Handle)

	b.Start()
	return nil
}

func checkoutState(flow *lflow.Flow, e string) error {
	//todo handle null flow
	f := fsm.NewFSM(
		stateFirst,
		fsm.Events{
			{Name: eventAskedFirst, Src: []string{stateFirst}, Dst: stateSecond},
			{Name: eventAskedThird, Src: []string{stateSecond}, Dst: stateThird},
			{Name: eventAskedSecond, Src: []string{stateThird}, Dst: stateLast},
			{Name: eventClose, Src: []string{stateLast}, Dst: stateClosed},
		},
		fsm.Callbacks{},
	)
	f.SetState(flow.GetCurrentState())
	if err := f.Event(context.Background(), e); err != nil {
		return fmt.Errorf("event: %w", err)
	}
	flow.SetState(f.Current())
	return nil
}

// как бы это прокидывать в tele.Context и доставать как c.Flow...
func getCurrentFlow(c tele.Context) *lflow.Flow {
	tele.NewContext(c.Bot(), c.Update())
	return c.Get("flow").(*lflow.Flow)
}
