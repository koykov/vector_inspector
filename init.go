package vector_inspector

import (
	"github.com/koykov/inspector"
	"github.com/koykov/x2bytes"
)

func init() {
	inspector.RegisterInspector("vector", &VectorInspector{})
	x2bytes.RegisterToBytesFn(VectorNodeToBytes)
}
