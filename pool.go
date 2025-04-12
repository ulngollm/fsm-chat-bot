package teleflow

type Pool interface {
	Add(flow Flow)
	Get(id int64) (Flow, error)
	Remove(flow Flow) error
}

type MemoryPool struct {
	data map[int64]Flow
}

func NewMemoryPool() *MemoryPool {
	return &MemoryPool{
		data: make(map[int64]Flow),
	}
}

func (p *MemoryPool) Add(flow Flow) {
	p.data[flow.ID()] = flow
}

func (p *MemoryPool) Get(id int64) (Flow, error) {
	return p.data[id], nil
}

func (p *MemoryPool) Remove(flow Flow) error {
	delete(p.data, flow.ID())
	return nil
}
