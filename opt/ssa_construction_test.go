package opt_test

import (
	"testing"

	usmssa "alon.kr/x/usm/usm/ssa"
)

func TestSsaConstruction(t *testing.T) {
	RunOptimizationTests(t, "ssa", usmssa.FunctionToSsaForm)
}
