package middleware

import (
	"fmt"

	lflow "github.com/ulngollm/msg-constructor/internal/flow"
	tele "gopkg.in/telebot.v4"
)

type FlowFinder struct {
	flowManager *lflow.Manager
	handlers    map[string]*FlowHandler
}

func NewFlowFinder(flowManager *lflow.Manager) *FlowFinder {
	return &FlowFinder{
		flowManager: flowManager,
		handlers:    make(map[string]*FlowHandler),
	}
}

func (f *FlowFinder) RegisterFlowHandler(flow string, handler *FlowHandler) {
	f.handlers[flow] = handler
}

func (f *FlowFinder) Handle(initial tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		flowID := c.Sender().ID
		existsFlow, err := f.flowManager.GetFlow(flowID)
		if err != nil {
			return fmt.Errorf("getFlow: %w", err)
		}
		if existsFlow == nil {
			return initial(c)
		}
		var found *FlowHandler
		for flow, handler := range f.handlers {
			if existsFlow.IsCurrentFlow(flow) {
				found = handler
				break
			}
		}
		c.Set("flow", existsFlow)
		return found.Handle(initial)(c)
	}
}
