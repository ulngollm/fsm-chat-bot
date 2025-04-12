package _default

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

func (r *DefaultFlowController) HandleInit(c tele.Context) error {
	id := c.Sender().ID
	flow := teleflow.NewSimpleFlow(id, StateFirst, "default")
	if err := r.manager.InitFlow(flow); err != nil {
		return err
	}
	teleflow.SaveToCtx(c, flow)
	return c.Send("hello! let's get started")
}

func (r *DefaultFlowController) HandleFirst(c tele.Context) error {
	flow := teleflow.GetCurrentFlow(c)
	err := r.checkoutState(flow, EventAskedFirst)
	flow.SetData(c.Message().Text)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("got it! Next")
}

func (r *DefaultFlowController) HandleSecond(c tele.Context) error {
	f := teleflow.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	err := r.checkoutState(f, EventAskedThird)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("next")
}

func (r *DefaultFlowController) HandlerThird(c tele.Context) error {
	f := teleflow.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	err := r.checkoutState(f, EventAskedSecond)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("next")
}

func (r *DefaultFlowController) HandleLast(c tele.Context) error {
	f := teleflow.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	err := r.checkoutState(f, EventClose)
	if err != nil {
		return fmt.Errorf("checkoutState: %v", err)
	}
	return c.Send("and final")
}

func (r *DefaultFlowController) HandleClose(c tele.Context) error {
	f := teleflow.GetCurrentFlow(c)
	f.SetData(fmt.Sprintf("%s %s", f.Data(), c.Message().Text))
	//not require to invalidate
	if err := r.manager.InvalidateFlow(f); err != nil {
		return fmt.Errorf("invalidateFlow: %w", err)
	}
	return c.Send(f.Data())
}

func (r *DefaultFlowController) checkoutState(flow teleflow.Flow, event string) error {
	sm := fsm.NewFSM(
		StateFirst,
		fsm.Events{
			{Name: EventAskedFirst, Src: []string{StateFirst}, Dst: StateSecond},
			{Name: EventAskedThird, Src: []string{StateSecond}, Dst: StateThird},
			{Name: EventAskedSecond, Src: []string{StateThird}, Dst: StateLast},
			{Name: EventClose, Src: []string{StateLast}, Dst: StateClosed},
		},
		fsm.Callbacks{},
	)
	sm.SetState(flow.State())
	if err := sm.Event(context.Background(), event); err != nil {
		return fmt.Errorf("event: %w", err)
	}
	flow.SetState(sm.Current())
	return nil
}
