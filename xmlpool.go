package vector_inspector

import "github.com/koykov/xmlvector"

type xmlPool struct{}

func (p *xmlPool) Get() any {
	return xmlvector.Acquire()
}

func (p *xmlPool) Reset(x any) {
	vec, ok := x.(*xmlvector.Vector)
	if !ok {
		return
	}
	vec.Reset()
}

func (p *xmlPool) Put(x any) {
	vec, ok := x.(*xmlvector.Vector)
	if !ok {
		return
	}
	xmlvector.Release(vec)
}
