// This file implements the construction of the dominator tree data structure.
//
// Construction algorithm is based on the Lengauer-Tarjan algorithm:
// https://doi.org/10.1145/357062.357071
//
// A resource I found (really!) helpful is Henrik Knakkegaard Christensen's
// master's thesis on "Algorithms for Finding Dominators in Directed Graphs".
// https://users-cs.au.dk/gerth/advising/thesis/henrik-knakkegaard-christensen.pdf
// Especially:
// - Section 2.6 (page 14): Dominator tree properties
// - Section 3.3 (page 27): Semidominators and how to compute them
// - Section 3.4 (page 30): Lengauer-Tarjan algorithm (besides the semidominators)
//
// I've also used the "Static Single Assignment Book" extensively:
// https://pfalcon.github.io/ssabook/latest/book-full.pdf

package control_flow

type dominatorTreeBuilder struct {
	ControlFlowGraph ControlFlowGraph
	LinkEvalForest   LinkEvalForest

	OriginalToPreorder []uint
	PreorderToOriginal []uint
}

func reversePermutation(p []uint) []uint {
	n := len(p)
	q := make([]uint, n)
	for i, v := range p {
		q[v] = uint(i)
	}
	return q
}

func newDominatorTreeBuilder(cfg ControlFlowGraph) dominatorTreeBuilder {
	n := cfg.Size()

	builder := dominatorTreeBuilder{
		ControlFlowGraph: cfg,
		LinkEvalForest:   NewLinkEvalForest(n),
	}

	builder.OriginalToPreorder = cfg.PreOrderDfs()
	builder.PreorderToOriginal = reversePermutation(builder.OriginalToPreorder)

	return builder
}
