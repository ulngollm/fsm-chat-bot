package middleware

import (
	"fmt"

	lflow "github.com/ulngollm/msg-constructor/internal/flow"
	tele "gopkg.in/telebot.v4"
)

type flowName = string
type flowState = string

type FlowGroup map[flowState]tele.HandlerFunc

// если тут нет pointer, он будет копироваться или один и тот же передаваться?
func (g FlowGroup) AddHandler(state flowState, handler tele.HandlerFunc) {
	g[state] = handler
}

func (g FlowGroup) GetHandlerForCurrentState(currentState string) tele.HandlerFunc {
	return g[currentState] // todo может вернуть ошибку?
}

type FlowFinder struct {
	flowManager *lflow.Manager
	handlers    map[flowName]FlowGroup
}

func NewFlowFinder(flowManager *lflow.Manager) *FlowFinder {
	return &FlowFinder{
		flowManager: flowManager,
		handlers:    make(map[flowName]FlowGroup),
	}
}

// todo добавлять отдельно для группы
func (f *FlowFinder) AddHandler(flow string, state string, handler tele.HandlerFunc) {
	f.handlers[flow][state] = handler
}

func (f *FlowFinder) Group(flow string) FlowGroup {
	f.handlers[flow] = make(FlowGroup)
	return f.handlers[flow]
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

		var group FlowGroup
		for flow, g := range f.handlers {
			if existsFlow.IsCurrentFlow(flow) {
				group = g
				break
			}
		}

		state := existsFlow.GetCurrentState()
		handler, ok := group[state]
		if !ok {
			return initial(c)
		}
		c.Set("flow", existsFlow)
		return handler(c)
	}
}
