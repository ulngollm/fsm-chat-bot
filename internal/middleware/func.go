package middleware

import (
	"github.com/ulngollm/msg-constructor/internal/flow"
	tele "gopkg.in/telebot.v4"
)

func GetCurrentFlow(c tele.Context) *flow.Flow {
	return flow.GetFromContext(c)
}
