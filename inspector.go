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

func (i *VectorInspector) TypeName() string {
	return "vector"
}

func (i *VectorInspector) Get(src interface{}, path ...string) (interface{}, error) {
	var buf interface{}
	err := i.GetTo(src, &buf, path...)
	return buf, err
}

func (i *VectorInspector) GetTo(src interface{}, buf *interface{}, path ...string) (err error) {
	if src == nil {
		return
	}
	var (
		node *vector.Node
	)
	if vec, ok := src.(*vector.Vector); ok {
		node = vec.Get(path...)
	} else if root, ok := src.(*vector.Node); ok {
		node = root.Get(path...)
	} else {
		return
	}
	*buf = node
	return
}

func (i *VectorInspector) Set(_, _ interface{}, _ ...string) error {
	return nil
}

func (i *VectorInspector) SetWB(_, _ interface{}, _ inspector.AccumulativeBuffer, _ ...string) error {
	return nil
}

func (i *VectorInspector) Cmp(src interface{}, cond inspector.Op, right string, result *bool, path ...string) error {
	var (
		node *vector.Node
	)
	if vec, ok := src.(*vector.Vector); ok {
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

// todo cover me with test/bench
func (i *VectorInspector) Loop(src interface{}, l inspector.Looper, buf *[]byte, path ...string) error {
	var (
		node *vector.Node
	)
	if vec, ok := src.(*vector.Vector); ok {
		node = vec.Get(path...)
	} else if root, ok := src.(*vector.Node); ok {
		node = root.Get(path...)
	} else {
		return nil
	}

	node.Each(func(idx int, child *vector.Node) {
		if l.RequireKey() {
			*buf = strconv.AppendInt((*buf)[:0], int64(idx), 10)
			l.SetKey(buf, &inspector.StaticInspector{})
		}
		l.SetVal(child, &VectorInspector{})
		ctl := l.Iterate()
		if ctl == inspector.LoopCtlBrk || ctl == inspector.LoopCtlCnt {
			return
		}
	})

	return nil
}

func (i *VectorInspector) DeepEqual(l, r interface{}) bool {
	return i.DeepEqualWithOptions(l, r, nil)
}

func (i *VectorInspector) DeepEqualWithOptions(l, r interface{}, opts *inspector.DEQOptions) bool {
	_, _, _ = l, r, opts
	// todo implement me; cover with test/bench
	return true
}

func (i *VectorInspector) Unmarshal(p []byte, typ inspector.Encoding) (interface{}, error) {
	switch typ {
	case inspector.EncodingJSON:
		vec := jsonvector.NewVector()
		err := vec.Parse(p)
		return vec, err
	default:
		return nil, inspector.ErrUnknownEncodingType
	}
}

func (i *VectorInspector) Copy(x interface{}) (interface{}, error) {
	// Vector/node copy is senseless.
	return x, nil
}

func (i *VectorInspector) cmpInt(left int64, cond inspector.Op, right int64) bool {
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

func (i *VectorInspector) cmpUint(left uint64, cond inspector.Op, right uint64) bool {
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

func (i *VectorInspector) cmpFloat(left float64, cond inspector.Op, right float64) bool {
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

func (i *VectorInspector) cmpStr(left string, cond inspector.Op, right string) bool {
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

func VectorNodeToBytes(dst []byte, val interface{}) ([]byte, error) {
	if node, ok := val.(*vector.Node); ok {
		dst = append(dst, node.Bytes()...)
	} else {
		return dst, x2bytes.ErrUnknownType
	}

	return dst, nil
}

func VectorNodeEmptyCheck(_ *dyntpl.Ctx, val interface{}) bool {
	if node, ok := val.(*vector.Node); ok {
		return node.Type() == vector.TypeNull
	}
	return false
}
