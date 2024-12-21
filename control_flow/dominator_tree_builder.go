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

	// The original control flow graph nodes are numbered from 0 to n-1.
	// In addition, the semi-dominator algorithm and the link-eval forest
	// requires the nodes to be numbered in a preorder traversal, from 0 to n-1.
	// Thus, we use those two representations interchangeably.
	// We refer to the original node numbers with the "orig" prefix and to the
	// dfs preorder numbers with the "preo" prefix.
	//
	// I admit, this is a bit confusing.
	OriginalToPreorder []uint
	PreorderToOriginal []uint

	// DfsParent[preo] is the parent of preo in the DFS spanning tree.
	DfsParent []uint
}

func reversePermutation(p []uint) []uint {
	n := len(p)
	q := make([]uint, n)
	for i, v := range p {
		q[v] = uint(i)
	}
	return q
}

// Implementation is mainly based on the description of Henrik Thesis,
// found here (section 3.3):
// https://users-cs.au.dk/gerth/advising/thesis/henrik-knakkegaard-christensen.pdf
func (b *dominatorTreeBuilder) calculateSemidominator(preoCurrent uint) uint {
	origCurrent := b.PreorderToOriginal[preoCurrent]
	currentBlock := b.BasicBlocks[origCurrent]

	for _, origPredecessor := range currentBlock.BackwardEdges {
		preoPredecessor := b.OriginalToPreorder[origPredecessor]

		candidate := b.LinkEvalForest.Eval(preoPredecessor)
		candidateSemidom := b.SemiDom[candidate]

		mySemidom := b.SemiDom[preoCurrent]
		if candidateSemidom < mySemidom {
			b.SemiDom[preoCurrent] = candidateSemidom
		}
	}

	origDfsParent := b.DfsParent[origCurrent]
	preoDfsParent := b.OriginalToPreorder[origDfsParent]
	b.Link(preoCurrent, preoDfsParent)

	return b.SemiDom[preoCurrent]
}

func NewDominatorTreeBuilder(cfg ControlFlowGraph) dominatorTreeBuilder {
	n := cfg.Size()

	builder := dominatorTreeBuilder{
		ControlFlowGraph: cfg,
		LinkEvalForest:   NewLinkEvalForest(n),
	}

	dfsResult := cfg.Dfs(CfgEntryBlock)
	builder.OriginalToPreorder = dfsResult.Preorder
	builder.PreorderToOriginal = reversePermutation(builder.OriginalToPreorder)
	builder.DfsParent = dfsResult.Parent

	return builder
}

func (b *dominatorTreeBuilder) Build() {
	n := b.Size()

	for preoCurrent := uint(n - 1); preoCurrent > 0; preoCurrent-- {
		b.calculateSemidominator(preoCurrent)
	}
}
