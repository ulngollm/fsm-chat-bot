package middleware

import tele "gopkg.in/telebot.v4"

type FlowGroup map[string]tele.HandlerFunc

// если тут нет pointer, он будет копироваться или один и тот же передаваться?
func (g FlowGroup) AddHandler(state string, handler tele.HandlerFunc) {
	g[state] = handler
}

func (g FlowGroup) GetHandlerForCurrentState(currentState string) tele.HandlerFunc {
	return g[currentState] // todo нужно вернуть ошибку, если нет обработчика?
}
