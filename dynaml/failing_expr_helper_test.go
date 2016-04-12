package dynaml

import (
	"github.com/hippotized/spiff/yaml"
)

type FailingExpr struct{}

func (FailingExpr) Evaluate(Binding) (yaml.Node, EvaluationInfo, bool) {
	return nil, DefaultInfo(), false
}
