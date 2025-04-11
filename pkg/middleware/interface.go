package middleware

import (
	tele "gopkg.in/telebot.v4"
)

type Middleware interface {
	Handle(next tele.HandlerFunc) tele.HandlerFunc
}
