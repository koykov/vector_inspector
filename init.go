package vector_inspector

import (
	"github.com/koykov/dyntpl"
	"github.com/koykov/inspector"
	"github.com/koykov/x2bytes"
)

func init() {
	inspector.RegisterInspector("vector", VectorInspector{})

	dyntpl.RegisterEmptyCheckFn("vector_node", VectorNodeEmptyCheck)
	_ = dyntpl.RegisterPool("jsonvector", &jsonPool{})
	_ = dyntpl.RegisterPool("urlvector", &urlPool{})
	_ = dyntpl.RegisterPool("halvector", &halPool{})
	_ = dyntpl.RegisterPool("xmlvector", &xmlPool{})
	// todo: register yamlvector

	x2bytes.RegisterToBytesFn(VectorNodeToBytes)
}
