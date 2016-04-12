package dynaml

import (
	"fmt"

	"github.com/hippotized/spiff/yaml"
)

type DivisionExpr struct {
	A Expression
	B Expression
}

func (e DivisionExpr) Evaluate(binding Binding) (yaml.Node, EvaluationInfo, bool) {
	resolved := true

	aint, info, ok := ResolveIntegerExpressionOrPushEvaluation(&e.A, &resolved, nil, binding)
	if !ok {
		return nil, info, false
	}

	bint, info, ok := ResolveIntegerExpressionOrPushEvaluation(&e.B, &resolved, &info, binding)
	if !ok {
		return nil, info, false
	}

	if !resolved {
		return node(e), info, true
	}

	if bint == 0 {
		info.Issue = "division by zero"
		return nil, info, false
	}
	return node(aint / bint), info, true
}

func (e DivisionExpr) String() string {
	return fmt.Sprintf("%s / %s", e.A, e.B)
}
