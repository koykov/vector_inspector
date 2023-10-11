package vector_inspector

import "github.com/koykov/jsonvector"

type jsonPool struct{}

func (p *jsonPool) Get() any {
	return jsonvector.Acquire()
}

func (p *jsonPool) Reset(x any) {
	vec, ok := x.(*jsonvector.Vector)
	if !ok {
		return
	}
	vec.Reset()
}

func (p *jsonPool) Put(x any) {
	vec, ok := x.(*jsonvector.Vector)
	if !ok {
		return
	}
	jsonvector.Release(vec)
}
