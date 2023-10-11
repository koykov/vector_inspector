package vector_inspector

import "github.com/koykov/urlvector"

type urlPool struct{}

func (p *urlPool) Get() any {
	return urlvector.Acquire()
}

func (p *urlPool) Reset(x any) {
	vec, ok := x.(*urlvector.Vector)
	if !ok {
		return
	}
	vec.Reset()
}

func (p *urlPool) Put(x any) {
	vec, ok := x.(*urlvector.Vector)
	if !ok {
		return
	}
	urlvector.Release(vec)
}
