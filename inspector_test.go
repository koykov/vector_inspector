package vector_inspector

import (
	"testing"

	"github.com/koykov/byteconv"
	"github.com/koykov/inspector"
	"github.com/koykov/jsonvector"
	"github.com/koykov/vector"
)

var (
	vecSrc0  = []byte(`{"color":{"value":"#c3c3c3"},"bgcolor":{"value":"#ffffff"},"inner_margin":{"value":15},"inner_padding":{"value":5},"border":{"value":"1 px solid #ccc"},"need_desc":{"value":true}}`)
	vecSrc1  = []byte(`{"color":{"value":"#c3c3c3"},"bgcolor":{"value":"#ffffff"},"inner_margin":{"value":15},"inner_padding":{"value":1500},"border":{"value":"1 px solid #ccc"},"need_desc":{"value":true}}`)
	p0color  = []string{"color", "value"}
	p0margin = []string{"inner_margin", "value"}
	p0desc   = []string{"need_desc", "value"}
	loopSrc  = []byte(`{"a":{"b":{"c":["foo","bar","string"]}}}`)
	loopExp  = []testExp{{"0", "foo"}, {"1", "bar"}, {"2", "string"}}
)

type testVI struct {
	val any
	ins inspector.Inspector
}

type testExp struct{ key, val string }

type testIterator struct {
	tb  testing.TB
	key testVI
	val testVI
	exp []testExp
	c   int
}

func (i *testIterator) RequireKey() bool                        { return true }
func (i *testIterator) SetKey(val any, ins inspector.Inspector) { i.key = testVI{val: val, ins: ins} }
func (i *testIterator) SetVal(val any, ins inspector.Inspector) { i.val = testVI{val: val, ins: ins} }
func (i *testIterator) Iterate() inspector.LoopCtl {
	exp := i.exp[i.c]
	key := byteconv.B2S(*i.key.val.(*[]byte))
	val := i.val.val.(*vector.Node).String()
	if exp.key != key {
		i.tb.Errorf("key mismatch: need '%s' got '%s'", exp.key, key)
	}
	if exp.val != val {
		i.tb.Errorf("val mismatch: need '%s' got '%s'", exp.val, val)
	}
	i.c++
	return inspector.LoopCtlNone
}
func (i *testIterator) reset() { i.c = 0 }

func TestVectorInspector(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		ins := VectorInspector{}
		vec := jsonvector.NewVector()
		if err := vec.Parse(vecSrc0); err != nil {
			t.Error(err)
		}
		var buf any

		_ = ins.GetTo(vec, &buf, p0color...)
		if buf.(*vector.Node).String() != "#c3c3c3" {
			t.Error("color.value mismatch: need #c3c3c3, got", buf)
		}

		_ = ins.GetTo(vec, &buf, p0margin...)
		if n, _ := buf.(*vector.Node).Int(); n != 15 {
			t.Error("inner_margin.value mismatch: need 15, got", buf)
		}

		_ = ins.GetTo(vec, &buf, p0desc...)
		if !buf.(*vector.Node).Bool() {
			t.Error("need_desc.value mismatch: need 15, got", buf)
		}
	})
	t.Run("compare", func(t *testing.T) {
		ins := VectorInspector{}
		vec := jsonvector.NewVector()
		if err := vec.Parse(vecSrc0); err != nil {
			t.Error(err)
		}
		var ok bool

		_ = ins.Compare(vec, inspector.OpLt, "18", &ok, p0margin...)
		if !ok {
			t.Error("inner_margin.value >= 18")
		}

		_ = ins.Compare(vec, inspector.OpGt, "13", &ok, p0margin...)
		if !ok {
			t.Error("inner_margin.value <= 13")
		}

		_ = ins.Compare(vec, inspector.OpEq, "15", &ok, p0margin...)
		if !ok {
			t.Error("inner_margin.value != 15")
		}
	})
	t.Run("compare", func(t *testing.T) {
		a, b := jsonvector.NewVector(), jsonvector.NewVector()
		_ = a.Parse(vecSrc0)
		_ = b.Parse(vecSrc0)
		var ins VectorInspector
		if !ins.DeepEqual(a, b) {
			t.FailNow()
		}
	})
	t.Run("compare", func(t *testing.T) {
		a, b := jsonvector.NewVector(), jsonvector.NewVector()
		_ = a.Parse(vecSrc0)
		_ = b.Parse(vecSrc1)
		var ins VectorInspector
		if ins.DeepEqual(a, b) {
			t.FailNow()
		}
	})
	t.Run("loop", func(t *testing.T) {
		ins := VectorInspector{}
		vec := jsonvector.NewVector()
		var buf []byte
		_ = vec.Parse(loopSrc)
		it := testIterator{tb: t, exp: loopExp}
		_ = ins.Loop(vec, &it, &buf, "a", "b", "c")
	})
}

func BenchmarkVectorInspector(b *testing.B) {
	b.Run("get", func(b *testing.B) {
		ins := VectorInspector{}
		vec := jsonvector.NewVector()
		if err := vec.Parse(vecSrc0); err != nil {
			b.Error(err)
		}
		var buf any

		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = ins.GetTo(vec, &buf, p0color...)
			if buf.(*vector.Node).String() != "#c3c3c3" {
				b.Error("color.value mismatch: need #c3c3c3, got", buf)
			}

			_ = ins.GetTo(vec, &buf, p0margin...)
			if n, _ := buf.(*vector.Node).Int(); n != 15 {
				b.Error("inner_margin.value mismatch: need 15, got", buf)
			}

			_ = ins.GetTo(vec, &buf, p0desc...)
			if !buf.(*vector.Node).Bool() {
				b.Error("need_desc.value mismatch: need 15, got", buf)
			}
		}
	})
	b.Run("compare", func(b *testing.B) {
		ins := VectorInspector{}
		vec := jsonvector.NewVector()
		if err := vec.Parse(vecSrc0); err != nil {
			b.Error(err)
		}
		var ok bool

		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = ins.Compare(vec, inspector.OpLt, "18", &ok, p0margin...)
			if !ok {
				b.Error("inner_margin.value >= 18")
			}

			_ = ins.Compare(vec, inspector.OpGt, "13", &ok, p0margin...)
			if !ok {
				b.Error("inner_margin.value <= 13")
			}

			_ = ins.Compare(vec, inspector.OpEq, "15", &ok, p0margin...)
			if !ok {
				b.Error("inner_margin.value != 15")
			}
		}
	})
	b.Run("compare", func(b *testing.B) {
		b.ReportAllocs()
		a, b_ := jsonvector.NewVector(), jsonvector.NewVector()
		_ = a.Parse(vecSrc0)
		_ = b_.Parse(vecSrc0)
		var ins VectorInspector
		for i := 0; i < b.N; i++ {
			if !ins.DeepEqual(a, b_) {
				b.FailNow()
			}
		}
	})
	b.Run("compare", func(b *testing.B) {
		b.ReportAllocs()
		a, b_ := jsonvector.NewVector(), jsonvector.NewVector()
		_ = a.Parse(vecSrc0)
		_ = b_.Parse(vecSrc1)
		var ins VectorInspector
		for i := 0; i < b.N; i++ {
			if ins.DeepEqual(a, b_) {
				b.FailNow()
			}
		}
	})
	b.Run("loop", func(b *testing.B) {
		b.ReportAllocs()
		ins := VectorInspector{}
		vec := jsonvector.NewVector()
		var buf []byte
		path := []string{"a", "b", "c"}
		_ = vec.Parse(loopSrc)
		it := testIterator{tb: b, exp: loopExp}
		for i := 0; i < b.N; i++ {
			buf = buf[:0]
			_ = ins.Loop(vec, &it, &buf, path...)
			it.reset()
		}
	})
}
