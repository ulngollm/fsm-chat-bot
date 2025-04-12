package teleflow

type FlowManager struct {
	pool Pool
}

func NewFlowManager(pool Pool) *FlowManager {
	return &FlowManager{pool: pool}
}

func (h *FlowManager) GetFlow(id int64) (Flow, error) {
	s, err := h.pool.Get(id)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (h *FlowManager) InitFlow(flow Flow) error {
	h.pool.Add(flow)
	return nil
}

func (h *FlowManager) InvalidateFlow(flow Flow) error {
	return h.pool.Remove(flow)
}
