package flow

type Manager struct {
	pool *Pool
}

func New(pool *Pool) *Manager {
	return &Manager{pool: pool}
}

func (h *Manager) GetFlow(id int64) (*Flow, error) {
	if s, ok := h.pool.GetFromPool(id); ok {
		return s, nil
	}
	return nil, nil
}

func (h *Manager) InitFlow(id int64, initialState string, key string) (*Flow, error) {
	s := &Flow{
		id:    id,
		state: initialState,
		key:   key,
	}
	h.pool.AddToPool(s)
	return s, nil
}

func (h *Manager) InvalidateFlow(flow *Flow) error {
	return h.pool.RemoveFromPool(flow)
}
