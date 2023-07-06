package vector_inspector

import (
	"testing"

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
)

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
}

func BenchmarkVectorInspector_Get(b *testing.B) {
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
}
