package main

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
	"github.com/ulngollm/teleflow"
	tele "gopkg.in/telebot.v4"
)

type DefaultFlowController struct {
	manager *teleflow.FlowManager
}

func NewDefaultFlowController(manager *teleflow.FlowManager) *DefaultFlowController {
	return &DefaultFlowController{manager: manager}
}

func (r *DefaultFlowController) handleInit(c tele.Context) error {
	id := c.Sender().ID
	flow := teleflow.NewSimpleFlow(id, stateFirst, "default")
	if err := r.manager.InitFlow(flow); err != nil {
		return err
	}
	teleflow.SaveToCtx(c, flow)
	return c.Send("hello! let's get started")
}

func (r *DefaultFlowController) handleFirst(c tele.Context) error {
	flow := teleflow.GetCurrentFlow(c)
	err := r.CheckoutState(flow, eventAskedFirst)
	flow.SetData(c.Message().Text)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("got it! Next")
}

func (r *DefaultFlowController) handleSecond(c tele.Context) error {
	f := teleflow.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	err := r.CheckoutState(f, eventAskedThird)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("next")
}

func (r *DefaultFlowController) handlerThird(c tele.Context) error {
	f := teleflow.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	err := r.CheckoutState(f, eventAskedSecond)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("next")
}

func (r *DefaultFlowController) handleLast(c tele.Context) error {
	f := teleflow.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	err := r.CheckoutState(f, eventClose)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("and final")
}

func (r *DefaultFlowController) handleClose(c tele.Context) error {
	f := teleflow.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	//not require to invalidate
	if err := r.manager.InvalidateFlow(f); err != nil {
		return fmt.Errorf("invalidateFlow: %w", err)
	}
	return c.Send(f.Data())
}

func (r *DefaultFlowController) CheckoutState(flow teleflow.Flow, e string) error {
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
	sm.SetState(flow.State())
	if err := sm.Event(context.Background(), e); err != nil {
		return fmt.Errorf("event: %w", err)
	}
	flow.SetState(sm.Current())
	return nil
}
