package middleware

import (
	lflow "github.com/ulngollm/msg-constructor/internal/flow"
	tele "gopkg.in/telebot.v4"
)

type Middleware interface {
	Handle(next tele.HandlerFunc) tele.HandlerFunc
}

type FlowController interface {
	CheckoutState(flow *lflow.Flow, e string) error
}
