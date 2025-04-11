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
	"github.com/ulngollm/msg-constructor/internal/state"
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
	b, err := bot.New(opts.BotToken)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	pool := lflow.NewPool()
	flowManager := lflow.New(pool)
	stateManager := state.NewStateManager()

	defaultFlowHandler := middleware.NewFlowHandler("default", stateFirst, flowManager, stateManager)
	defaultFlowHandler.AddStateHandler(stateFirst, handleFirst)
	defaultFlowHandler.AddStateHandler(stateSecond, handleSecond)
	defaultFlowHandler.AddStateHandler(stateThird, handlerThird)
	defaultFlowHandler.AddStateHandler(stateLast, handleLast)
	defaultFlowHandler.AddStateHandler(stateClosed, func(c tele.Context) error {
		flow := getCurrentFlow(c)
		if err := flowManager.InvalidateFlow(flow); err != nil {
			return fmt.Errorf("invalidateFlow: %w", err)
		}
		return c.Send(flow.Data())
	})

	b.RegisterFlowHandler(tele.OnText, nil, defaultFlowHandler)

	b.Start()
	return nil
}

func handleFirst(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := checkoutState(flow, eventAskedFirst)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("first")
}

func handleSecond(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := checkoutState(flow, eventAskedThird)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("2")
}

func handlerThird(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := checkoutState(flow, eventAskedSecond)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("3")
}

func handleLast(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := checkoutState(flow, eventClose)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("finally")
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

func getCurrentFlow(c tele.Context) *lflow.Flow {
	return c.Get("flow").(*lflow.Flow)
}
