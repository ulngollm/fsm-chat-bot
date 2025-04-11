package flow

type Pool struct {
	data map[int64]*Flow
}

func NewPool() *Pool {
	return &Pool{
		data: make(map[int64]*Flow),
	}
}

func (p *Pool) AddToPool(flow *Flow) {
	p.data[flow.id] = flow
}

func (p *Pool) GetFromPool(id int64) (*Flow, bool) {
	v, ok := p.data[id]
	if ok {
		return v, true
	}
	return &Flow{}, false
}

func (p *Pool) RemoveFromPool(flow *Flow) error {
	delete(p.data, flow.id)
	return nil
}
