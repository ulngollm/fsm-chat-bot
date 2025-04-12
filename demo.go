package main

import (
	"context"
	"fmt"

	"github.com/ulngollm/msg-constructor/internal/middleware"

	"github.com/looplab/fsm"
	"github.com/ulngollm/msg-constructor/internal/flow"
	tele "gopkg.in/telebot.v4"
)

type DefaultFlowController struct {
	manager *flow.Manager
}

func NewDefaultFlowController(manager *flow.Manager) *DefaultFlowController {
	return &DefaultFlowController{manager: manager}
}

func (r *DefaultFlowController) handleInit(c tele.Context) error {
	id := c.Sender().ID
	f, err := r.manager.InitFlow(id, stateFirst, "default")
	if err != nil {
		return err
	}
	f.SaveToCtx(c)
	return c.Send("hello! let's get started")
}

func (r *DefaultFlowController) handleFirst(c tele.Context) error {
	f := middleware.GetCurrentFlow(c)
	err := r.CheckoutState(f, eventAskedFirst)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("got it! Next")
}

func (r *DefaultFlowController) handleSecond(c tele.Context) error {
	f := middleware.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	err := r.CheckoutState(f, eventAskedThird)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("next")
}

func (r *DefaultFlowController) handlerThird(c tele.Context) error {
	f := middleware.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	err := r.CheckoutState(f, eventAskedSecond)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("next")
}

func (r *DefaultFlowController) handleLast(c tele.Context) error {
	f := middleware.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	err := r.CheckoutState(f, eventClose)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("and final")
}

func (r *DefaultFlowController) handleClose(c tele.Context) error {
	f := middleware.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	//not require to invalidate
	if err := r.manager.InvalidateFlow(f); err != nil {
		return fmt.Errorf("invalidateFlow: %w", err)
	}
	return c.Send(f.Data())
}

func (r *DefaultFlowController) CheckoutState(flow *flow.Flow, e string) error {
	//todo handle null flow
	sm := fsm.NewFSM(
		stateFirst,
		fsm.Events{
			{Name: eventAskedFirst, Src: []string{stateFirst}, Dst: stateSecond},
			{Name: eventAskedThird, Src: []string{stateSecond}, Dst: stateThird},
			{Name: eventAskedSecond, Src: []string{stateThird}, Dst: stateLast},
			{Name: eventClose, Src: []string{stateLast}, Dst: stateClosed},
		},
		fsm.Callbacks{},
	)
	sm.SetState(flow.GetCurrentState())
	if err := sm.Event(context.Background(), e); err != nil {
		return fmt.Errorf("event: %w", err)
	}
	flow.SetState(sm.Current())
	return nil
}
