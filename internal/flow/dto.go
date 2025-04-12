package flow

import tele "gopkg.in/telebot.v4"

const cacheKey = "flow"

type Flow struct {
	id    int64
	state string
	key   string
	data  string // пока строка
}

func GetFromContext(c tele.Context) *Flow {
	return c.Get(cacheKey).(*Flow)
}

func (s *Flow) GetCurrentState() string {
	return s.state
}

func (s *Flow) SetState(state string) {
	s.state = state
}

func (s *Flow) InitState(state string) {
	s.state = state
}

// метод нужен для того, чтобы сопоставлять flowHandler и flow
func (s *Flow) IsCurrentFlow(key string) bool {
	return s.key == key
}

func (s *Flow) Data() string {
	return s.data
}

func (s *Flow) SetData(data string) {
	s.data = data
}

func (s *Flow) SaveToCtx(c tele.Context) {
	c.Set(cacheKey, s)
}
