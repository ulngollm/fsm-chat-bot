package middleware

import (
	"github.com/ulngollm/msg-constructor/internal/flow"
	"github.com/ulngollm/msg-constructor/internal/state"
	tele "gopkg.in/telebot.v4"
)

type FlowHandler struct {
	flowName     string
	stateManager *state.Manager
}

// по сути это роутер
// у flow есть steps (state)
// у шагов есть обработчик (строго 1-1)

func NewFlowHandler(name string, stateManager *state.Manager) *FlowHandler {
	return &FlowHandler{flowName: name, stateManager: stateManager}
}

// хендлеры сами переключают стейты
func (m *FlowHandler) AddStateHandler(state string, handler tele.HandlerFunc) {
	m.stateManager.AddHandler(state, handler)
}

// этот метод не уникален для flowGroup
// здесь ничего не используется от flowGroup
// flowGroup должно быть много, а такой метод на всех один. Может вынести его в FlowHandler ?
func (m *FlowHandler) Handle(defaultHandler tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		fl := c.Get("flow").(*flow.Flow)
		if fl == nil {
			return defaultHandler(c)
		}
		// todo изменить механизм поиска хендлера для стейта
		// искать в стейтах хендлера
		handler, err := m.stateManager.GetHandlerForCurrentState(fl.GetCurrentState())
		if err != nil {
			return err
		}
		if handler == nil {
			return defaultHandler(c)
		}
		return handler(c)
	}
}

//todo вообще именно middleware должен быть один
// и искать соответсвующий найденному flow flowGroup
// главная задача flowGroup - вернуть tele.Handler для стейта

//итак, разделение ответственности:
// middleware ищет нужный flowGroup по flow
// нужный flowGroup ищет нужный обработчик для стейта

// todo только вопрос, откуда стартует этот процесс? откуда возьмется первый flow? откуда запустится?
// что если он будет не найден - какой инициализируется?
// возможное решение: middleware, который ищет flow + defaultHandler для эндпоинта:
// на обычный эндпоинт маппишь инициализатор
// на входе запрос перехватит middleware и если найдет инициализированный flow - продолжит его
// если не найдет - тогда вызовется инициализатор стейта и все будет ок
