package vector_inspector

import (
	"github.com/koykov/dyntpl"
	"github.com/koykov/inspector"
	"github.com/koykov/x2bytes"
)

func init() {
	inspector.RegisterInspector("vector", VectorInspector{})
	dyntpl.RegisterEmptyCheckFn("vector_node", VectorNodeEmptyCheck)
	x2bytes.RegisterToBytesFn(VectorNodeToBytes)
}
