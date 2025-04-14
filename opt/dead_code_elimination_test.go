package opt_test

import (
	"testing"

	"alon.kr/x/usm/opt"
)

func TestDeadCodeElimination(t *testing.T) {
	RunOptimizationTests(t, "dead_code_elimination", opt.DeadCodeElimination)
}
