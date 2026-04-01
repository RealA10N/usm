package opt_test

import (
	"testing"

	"alon.kr/x/usm/opt"
)

func TestConstantPropagation(t *testing.T) {
	RunOptimizationTests(t, "constant_propagation", opt.ConstantPropagation)
}
