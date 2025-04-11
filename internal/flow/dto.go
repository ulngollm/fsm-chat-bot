package flow

// flow еще могут быть разных назначений и могут быть разные fsm
type Flow struct {
	id    int64
	state string
	key   string
	//	fsm будет снаружи. Не надо столько в памяти хранить
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

func (s *Flow) IsCurrentFlow(key string) bool {
	return s.key == key
}
