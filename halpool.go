package vector_inspector

import "github.com/koykov/halvector"

type halPool struct{}

func (p *halPool) Get() any {
	return halvector.Acquire()
}

func (p *halPool) Reset(x any) {
	vec, ok := x.(*halvector.Vector)
	if !ok {
		return
	}
	vec.Reset()
}

func (p *halPool) Put(x any) {
	vec, ok := x.(*halvector.Vector)
	if !ok {
		return
	}
	halvector.Release(vec)
}
