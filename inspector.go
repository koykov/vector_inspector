package vector_inspector

import (
	"strconv"

	"github.com/koykov/dyntpl"
	"github.com/koykov/inspector"
	"github.com/koykov/jsonvector"
	"github.com/koykov/vector"
	"github.com/koykov/x2bytes"
)

type VectorInspector struct {
	inspector.BaseInspector
}

func (i VectorInspector) TypeName() string {
	return "vector"
}

func (i VectorInspector) Get(src any, path ...string) (any, error) {
	var buf any
	err := i.GetTo(src, &buf, path...)
	return buf, err
}

func (i VectorInspector) GetTo(src any, buf *any, path ...string) (err error) {
	if src == nil {
		return
	}
	var (
		node *vector.Node
	)
	if vec, ok := src.(vector.Interface); ok && vec != nil {
		node = vec.Get(path...)
	} else if vec, ok = src.(*vector.Vector); ok {
		node = vec.Get(path...)
	} else if root, ok := src.(*vector.Node); ok {
		node = root.Get(path...)
	} else {
		return
	}
	*buf = node
	return
}

func (i VectorInspector) Set(_, _ any, _ ...string) error {
	// Vector is read-only struct.
	return nil
}

func (i VectorInspector) SetWithBuffer(_, _ any, _ inspector.AccumulativeBuffer, _ ...string) error {
	// Vector is read-only struct.
	return nil
}

func (i VectorInspector) Compare(src any, cond inspector.Op, right string, result *bool, path ...string) error {
	var (
		node *vector.Node
	)
	if vec, ok := src.(vector.Interface); ok && vec != nil {
		node = vec.Get(path...)
	} else if vec, ok = src.(*vector.Vector); ok {
		node = vec.Get(path...)
	} else if root, ok := src.(*vector.Node); ok {
		node = root.Get(path...)
	} else {
		*result = false
		return nil
	}

	switch node.Type() {
	case vector.TypeStr, vector.TypeAttr, vector.TypeBool:
		*result = i.cmpStr(node.String(), cond, right)
	case vector.TypeNum:
		if r, err := strconv.ParseInt(right, 0, 0); err == nil {
			n, _ := node.Int()
			*result = i.cmpInt(n, cond, r)
		} else if r, err := strconv.ParseUint(right, 0, 0); err == nil {
			n, _ := node.Uint()
			*result = i.cmpUint(n, cond, r)
		} else if r, err := strconv.ParseFloat(right, 0); err == nil {
			n, _ := node.Float()
			*result = i.cmpFloat(n, cond, r)
		}
	default:
		*result = false
	}
	return nil
}

func (i VectorInspector) Loop(src any, l inspector.Iterator, buf *[]byte, path ...string) error {
	// todo cover me with test/bench
	var node *vector.Node
	if vec, ok := src.(vector.Interface); ok && vec != nil {
		node = vec.Get(path...)
	} else if vec, ok = src.(*vector.Vector); ok {
		node = vec.Get(path...)
	} else if root, ok := src.(*vector.Node); ok {
		node = root.Get(path...)
	} else {
		return nil
	}

	node.Each(func(idx int, child *vector.Node) {
		if l.RequireKey() {
			*buf = strconv.AppendInt((*buf)[:0], int64(idx), 10)
			l.SetKey(buf, inspector.StaticInspector{})
		}
		l.SetVal(child, VectorInspector{})
		ctl := l.Iterate()
		if ctl == inspector.LoopCtlBrk || ctl == inspector.LoopCtlCnt {
			return
		}
	})

	return nil
}

func (i VectorInspector) DeepEqual(l, r any) bool {
	return i.DeepEqualWithOptions(l, r, nil)
}

func (i VectorInspector) DeepEqualWithOptions(l, r any, _ *inspector.DEQOptions) bool {
	var a, b *vector.Node
	if vec, ok := l.(vector.Interface); ok && vec != nil {
		a = vec.Root()
	} else if vec, ok = l.(*vector.Vector); ok {
		a = vec.Root()
	} else if root, ok := l.(*vector.Node); ok {
		a = root
	} else {
		return false
	}

	if vec, ok := r.(vector.Interface); ok && vec != nil {
		b = vec.Root()
	} else if vec, ok = r.(*vector.Vector); ok {
		b = vec.Root()
	} else if root, ok := r.(*vector.Node); ok {
		b = root
	} else {
		return false
	}

	_, _ = a, b
	ok := a.EqualWith(b)
	return ok
}

func (i VectorInspector) Unmarshal(p []byte, typ inspector.Encoding) (any, error) {
	switch typ {
	case inspector.EncodingJSON:
		vec := jsonvector.NewVector()
		err := vec.Parse(p)
		return vec, err
	default:
		return nil, inspector.ErrUnknownEncodingType
	}
}

func (i VectorInspector) Copy(x any) (any, error) {
	// Vector/node copy is senseless.
	return x, nil
}

func (i VectorInspector) CopyTo(src, dst any, _ inspector.AccumulativeBuffer) error {
	_, _ = src, dst
	// Vector/node copy is senseless.
	return nil
}

func (i VectorInspector) Length(src any, result *int, path ...string) error {
	var node *vector.Node
	if vec, ok := src.(vector.Interface); ok && vec != nil {
		node = vec.Get(path...)
	} else if vec, ok = src.(*vector.Vector); ok {
		node = vec.Get(path...)
	} else if root, ok := src.(*vector.Node); ok {
		node = root.Get(path...)
	} else {
		return nil
	}

	*result = node.Get(path...).Limit()
	return nil
}

func (i VectorInspector) Capacity(src any, result *int, path ...string) error {
	return i.Length(src, result, path...)
}

func (i VectorInspector) Reset(x any) error {
	if vec, ok := x.(vector.Interface); ok && vec != nil {
		vec.Reset()
		return nil
	} else if vec, ok = x.(*vector.Vector); ok {
		vec.Reset()
		return nil
	} else if root, ok := x.(*vector.Node); ok {
		root.Reset()
		return nil
	}
	return inspector.ErrUnsupportedType
}

func (i VectorInspector) cmpInt(left int64, cond inspector.Op, right int64) bool {
	switch cond {
	case inspector.OpEq:
		return left == right
	case inspector.OpNq:
		return left != right
	case inspector.OpGt:
		return left > right
	case inspector.OpGtq:
		return left >= right
	case inspector.OpLt:
		return left < right
	case inspector.OpLtq:
		return left <= right
	}
	return false
}

func (i VectorInspector) cmpUint(left uint64, cond inspector.Op, right uint64) bool {
	switch cond {
	case inspector.OpEq:
		return left == right
	case inspector.OpNq:
		return left != right
	case inspector.OpGt:
		return left > right
	case inspector.OpGtq:
		return left >= right
	case inspector.OpLt:
		return left < right
	case inspector.OpLtq:
		return left <= right
	}
	return false
}

func (i VectorInspector) cmpFloat(left float64, cond inspector.Op, right float64) bool {
	switch cond {
	case inspector.OpEq:
		return left == right
	case inspector.OpNq:
		return left != right
	case inspector.OpGt:
		return left > right
	case inspector.OpGtq:
		return left >= right
	case inspector.OpLt:
		return left < right
	case inspector.OpLtq:
		return left <= right
	}
	return false
}

func (i VectorInspector) cmpStr(left string, cond inspector.Op, right string) bool {
	switch cond {
	case inspector.OpEq:
		return left == right
	case inspector.OpNq:
		return left != right
	case inspector.OpGt:
		return left > right
	case inspector.OpGtq:
		return left >= right
	case inspector.OpLt:
		return left < right
	case inspector.OpLtq:
		return left <= right
	}
	return false
}

func VectorNodeToBytes(dst []byte, val any) ([]byte, error) {
	if node, ok := val.(*vector.Node); ok {
		dst = append(dst, node.Bytes()...)
	} else {
		return dst, x2bytes.ErrUnknownType
	}

	return dst, nil
}

func VectorNodeEmptyCheck(_ *dyntpl.Ctx, val any) bool {
	if node, ok := val.(*vector.Node); ok {
		return node.Type() == vector.TypeNull
	}
	return false
}
