package middleware

import (
	"fmt"

	"github.com/ulngollm/msg-constructor/internal/flow"
	"github.com/ulngollm/msg-constructor/internal/state"
	tele "gopkg.in/telebot.v4"
)

type FlowHandler struct {
	flowName     string
	initialState string
	flowManager  *flow.Manager
	stateManager *state.Manager
}

func NewFlowHandler(flowName string, initialState string, flowManager *flow.Manager, stateManager *state.Manager) *FlowHandler {
	return &FlowHandler{flowName: flowName, initialState: initialState, flowManager: flowManager, stateManager: stateManager}
}

// хендлеры сами переключают стейты
func (m *FlowHandler) AddStateHandler(state string, handler tele.HandlerFunc) {
	m.stateManager.AddHandler(state, handler)
}

func (m *FlowHandler) Handle(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		flowID := c.Sender().ID
		fl, err := m.flowManager.GetFlow(flowID)
		if err != nil {
			return fmt.Errorf("getFlow: %w", err)
		}
		if fl != nil && !fl.IsCurrentFlow(m.flowName) {
			return nil // тогда другой обработчик этим должен заниматься
		}
		if fl == nil {
			fl, err = m.flowManager.InitFlow(flowID, m.initialState, m.flowName)
			if err != nil {
				return fmt.Errorf("initFlow: %w", err)
			}
		}
		c.Set("flow", fl) // see для того чтобы не искать flow заново в хендлерах

		handler, err := m.stateManager.GetHandlerForCurrentState(fl.GetCurrentState())
		if err != nil {
			return err
		}
		return handler(c)
	}
}
