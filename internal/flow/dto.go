package flow

// flow еще могут быть разных назначений и могут быть разные fsm
type Flow struct {
	id    int64
	state string
	key   string
	data  string // пока строка
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
// испльзовать паттерн цепочка обязанностей?
func (s *Flow) IsCurrentFlow(key string) bool {
	return s.key == key
}

func (s *Flow) Data() string {
	return s.data
}

func (s *Flow) SetData(data string) {
	s.data = data
}
