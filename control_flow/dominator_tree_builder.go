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
	ControlFlowGraph
	LinkEvalForest

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

	// The original control flow graph nodes are numbered from 0 to n-1.
	// In addition, the semi-dominator algorithm and the link-eval forest
	// requires the nodes to be numbered in a preorder traversal, from 0 to n-1.
	// Thus, we use those two representations interchangeably.
	// We refer to the original node numbers with the "orig" prefix and to the
	// dfs preorder numbers with the "preo" prefix.
	//
	// I admit, this is a bit confusing.

	dfsResult := cfg.Dfs(CfgEntryBlock)
	builder.OriginalToPreorder = dfsResult.Preorder
	builder.PreorderToOriginal = reversePermutation(builder.OriginalToPreorder)

	for preoCurrent := uint(n - 1); preoCurrent > 0; preoCurrent-- {
		origCurrent := builder.PreorderToOriginal[preoCurrent]
		currentBlock := cfg.BasicBlocks[origCurrent]

		for _, origPredecessor := range currentBlock.BackwardEdges {
			preoPredecessor := builder.OriginalToPreorder[origPredecessor]

			candidate := builder.LinkEvalForest.Eval(preoPredecessor)
			candidateSemidom := builder.SemiDom[candidate]

			mySemidom := builder.SemiDom[preoCurrent]
			if candidateSemidom < mySemidom {
				builder.SemiDom[preoCurrent] = candidateSemidom
			}
		}

		origDfsParent := dfsResult.Parent[origCurrent]
		preoDfsParent := builder.OriginalToPreorder[origDfsParent]
		builder.Link(preoCurrent, preoDfsParent)
	}

	return builder
}
