package dynaml

import (
	"strconv"

	"github.com/hippotized/spiff/yaml"
)

type IntegerExpr struct {
	Value int64
}

func (e IntegerExpr) Evaluate(Binding) (yaml.Node, EvaluationInfo, bool) {
	return node(e.Value), DefaultInfo(), true
}

func (e IntegerExpr) String() string {
	return strconv.FormatInt(e.Value, 10)
}
