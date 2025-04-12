package main

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
	lflow "github.com/ulngollm/msg-constructor/internal/flow"
	tele "gopkg.in/telebot.v4"
)

type DefaultFlowController struct {
	manager *lflow.Manager
}

func NewDefaultFlowController(manager *lflow.Manager) *DefaultFlowController {
	return &DefaultFlowController{manager: manager}
}

func (r *DefaultFlowController) handleInit(c tele.Context) error {
	id := c.Sender().ID
	flow, err := r.manager.InitFlow(id, stateFirst, "default")
	if err != nil {
		return err
	}
	c.Set("flow", flow)
	return c.Send("hello! let's get started")
}

func (r *DefaultFlowController) handleFirst(c tele.Context) error {
	flow := getCurrentFlow(c)
	err := r.CheckoutState(flow, eventAskedFirst)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("got it! Next")
}

func (r *DefaultFlowController) handleSecond(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := r.CheckoutState(flow, eventAskedThird)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("next")
}

func (r *DefaultFlowController) handlerThird(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := r.CheckoutState(flow, eventAskedSecond)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("next")
}

func (r *DefaultFlowController) handleLast(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	err := r.CheckoutState(flow, eventClose)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("and final")
}

func (r *DefaultFlowController) handleClose(c tele.Context) error {
	flow := getCurrentFlow(c)
	flow.SetData(fmt.Sprintf("%s %s", flow.Data(), c.Message().Text))
	//not require to invalidate
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
