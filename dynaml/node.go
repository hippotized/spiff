package dynaml

import (
	"github.com/hippotized/spiff/yaml"
)

func node(val interface{}) yaml.Node {
	return yaml.NewNode(val, "dynaml")
}
