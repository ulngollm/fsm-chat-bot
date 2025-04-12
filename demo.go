package main

import (
	"context"
	"fmt"
	"github.com/looplab/fsm"
	lflow "github.com/ulngollm/msg-constructor/internal/flow"
	tele "gopkg.in/telebot.v4"
)

//todo как должен выглядеть идеальный интерфейс:

type FlowController interface {
	CheckoutState(flow *lflow.Flow, e string) error
}

type DefaultFlowController struct {
	//	states можно хранить извне. Контроллер же ничего ен значет про роуты?
	manager *lflow.Manager
	//	flow manager
}

func NewDefaultFlowController(manager *lflow.Manager) *DefaultFlowController {
	return &DefaultFlowController{manager: manager}
}

func (r *DefaultFlowController) handleFirst(c tele.Context) error {
	flow := getCurrentFlow(c)
	err := checkoutState(flow, eventAskedFirst)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("first")
}

func (r *DefaultFlowController) handleSecond(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := checkoutState(flow, eventAskedThird)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("2")
}

func (r *DefaultFlowController) handlerThird(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := checkoutState(flow, eventAskedSecond)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("3")
}

func (r *DefaultFlowController) handleLast(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := checkoutState(flow, eventClose)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("finally")
}

func (r *DefaultFlowController) handleClose(c tele.Context) error {
	// todo чтобы не надо было доставать извне до flowManager, нужно
	// сделать хендлеры методами flow handler-a
	// flowManager инжектить в flowHandler
	flow := getCurrentFlow(c)
	if err := r.manager.InvalidateFlow(flow); err != nil {
		return fmt.Errorf("invalidateFlow: %w", err)
	}
	return c.Send(flow.Data())
}

func (r *DefaultFlowController) CheckoutState(flow *lflow.Flow, e string) error {
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
