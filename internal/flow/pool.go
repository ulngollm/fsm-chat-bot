package flow

type Pool struct {
	data map[int64]*Flow
}

// это должно быть в pkg в виде интерфейсов
// потому что должна быть возможность переопределить эту dummy thing и использовать другой пул
// с редисом, например. или с ttl

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
