package middleware

import (
	"fmt"

	lflow "github.com/ulngollm/msg-constructor/internal/flow"
	tele "gopkg.in/telebot.v4"
)

type flowName = string

type FlowRouter struct {
	flowManager *lflow.Manager
	handlers    map[flowName]FlowGroup
}

func NewFlowFinder(flowManager *lflow.Manager) *FlowRouter {
	return &FlowRouter{
		flowManager: flowManager,
		handlers:    make(map[flowName]FlowGroup),
	}
}

func (f *FlowRouter) Group(flow string) FlowGroup {
	f.handlers[flow] = make(FlowGroup)
	return f.handlers[flow]
}

func (f *FlowRouter) Handle(initial tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		flowID := c.Sender().ID
		existsFlow, err := f.flowManager.GetFlow(flowID)
		if err != nil {
			return fmt.Errorf("getFlow: %w", err)
		}
		if existsFlow == nil {
			return initial(c)
		}

		var g FlowGroup
		for flow, group := range f.handlers {
			if existsFlow.IsCurrentFlow(flow) {
				g = group
				break
			}
		}

		state := existsFlow.GetCurrentState()
		handler, ok := g[state]
		if !ok {
			return initial(c)
		}

		existsFlow.SaveToCtx(c)
		return handler(c)
	}
}
