package state

import (
	tele "gopkg.in/telebot.v4"
)

//а вот эта штука уйдет в flowHandler. Потому что стейты вне flow не имеют точности

type Manager struct {
	stateHandlers map[string]tele.HandlerFunc
}

func NewStateManager() *Manager {
	return &Manager{
		stateHandlers: make(map[string]tele.HandlerFunc),
	}
}

func (sm *Manager) GetHandler(state string) tele.HandlerFunc {
	return sm.stateHandlers[state]
}

func (sm *Manager) AddHandler(state string, handler tele.HandlerFunc) {
	sm.stateHandlers[state] = handler
}

func (sm *Manager) GetHandlerForCurrentState(currentState string) (tele.HandlerFunc, error) {
	return sm.GetHandler(currentState), nil
}
