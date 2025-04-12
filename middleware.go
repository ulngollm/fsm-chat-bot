package teleflow

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type flowName = string

type FlowRouter struct {
	flowManager *FlowManager
	handlers    map[flowName]FlowGroup
}

func NewFlowRouter(flowManager *FlowManager) *FlowRouter {
	return &FlowRouter{
		flowManager: flowManager,
		handlers:    make(map[flowName]FlowGroup),
	}
}

func (f *FlowRouter) Group(flow string) FlowGroup {
	f.handlers[flow] = make(FlowGroup)
	return f.handlers[flow]
}

func (f *FlowRouter) Middleware() tele.MiddlewareFunc {
	return func(initial tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			flowID := c.Sender().ID // todo maybe need customize
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

			state := existsFlow.State()
			handler, ok := g[state]
			if !ok {
				return initial(c)
			}

			SaveToCtx(c, existsFlow)
			return handler(c)
		}
	}
}
