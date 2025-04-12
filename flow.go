package teleflow

import tele "gopkg.in/telebot.v4"

const cacheKey = "flow"

type Flow interface {
	ID() int64
	State() string
	SetState(string)
	IsCurrentFlow(key string) bool
	Data() interface{}
	SetData(data interface{})
}

type SimpleFlow struct {
	id    int64
	state string
	key   string
	data  interface{}
}

func NewSimpleFlow(id int64, initialState string, key string) *SimpleFlow {
	return &SimpleFlow{id: id, state: initialState, key: key}
}

func (s *SimpleFlow) ID() int64 {
	return s.id
}

func (s *SimpleFlow) State() string {
	return s.state
}

func (s *SimpleFlow) SetState(state string) {
	s.state = state
}

func (s *SimpleFlow) IsCurrentFlow(key string) bool {
	return s.key == key
}

func (s *SimpleFlow) Data() interface{} {
	return s.data
}

func (s *SimpleFlow) SetData(data interface{}) {
	s.data = data
}

// helpers

func GetCurrentFlow(c tele.Context) Flow {
	return GetFromContext(c)
}

func GetFromContext(c tele.Context) Flow {
	return c.Get(cacheKey).(Flow)
}

func SaveToCtx(c tele.Context, flow Flow) {
	c.Set(cacheKey, flow)
}
