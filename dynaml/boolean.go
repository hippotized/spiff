package dynaml

import (
	"fmt"

	"github.com/hippotized/spiff/yaml"
)

type BooleanExpr struct {
	Value bool
}

func (e BooleanExpr) Evaluate(Binding) (yaml.Node, EvaluationInfo, bool) {
	return node(e.Value), DefaultInfo(), true
}

func (e BooleanExpr) String() string {
	return fmt.Sprintf("%v", e.Value)
}
