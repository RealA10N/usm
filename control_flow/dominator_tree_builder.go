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
	*ControlFlowGraph
	LinkEvalForest

	// ImmDom[preo] is the immediate dominator (index in preorder) of the
	// node 'preo', after it has been computed.
	ImmDom []uint

	// SemiDomBuckets[preo] contains a slice of all nodes (represented by index
	// in the preorder) that have their semidominator set to 'preo'.
	SemiDomBuckets [][]uint

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

func newDominatorTreeBuilder(cfg *ControlFlowGraph) dominatorTreeBuilder {
	n := cfg.Size()
	dfsResult := cfg.Dfs(CfgEntryBlock)

	builder := dominatorTreeBuilder{
		ControlFlowGraph:   cfg,
		LinkEvalForest:     NewLinkEvalForest(n),
		ImmDom:             make([]uint, n),
		SemiDomBuckets:     make([][]uint, n),
		OriginalToPreorder: dfsResult.Preorder,
		PreorderToOriginal: reversePermutation(dfsResult.Preorder),
		DfsParent:          dfsResult.Parent,
	}

	return builder
}

func reversePermutation(p []uint) []uint {
	n := len(p)
	q := make([]uint, n)
	for i, v := range p {
		q[v] = uint(i)
	}
	return q
}

// "Step 2" of the Lengauer-Tarjan algorithm, as it is described in the original
// paper.
//
// Implementation is mainly based on the description of Knakkegaard's Thesis,
// section 3.3.
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

	return b.SemiDom[preoCurrent]
}

func (b *dominatorTreeBuilder) linkToDfsParent(preoCurrent uint) uint {
	origCurrent := b.PreorderToOriginal[preoCurrent]
	origDfsParent := b.DfsParent[origCurrent]
	preoDfsParent := b.OriginalToPreorder[origDfsParent]
	b.Link(preoCurrent, preoDfsParent)
	return preoDfsParent
}

// "Step 3" of the Lengauer-Tarjan algorithm.
func (b *dominatorTreeBuilder) consumeSemiDominatorsBucket(preoParent uint) {
	for _, v := range b.SemiDomBuckets[preoParent] {
		u := b.Eval(v)
		if b.SemiDom[u] < b.SemiDom[v] {
			b.ImmDom[v] = u
		} else {
			b.ImmDom[v] = b.SemiDom[v]
		}
	}

	// The paper explicitly states that in the loop above, we should remove
	// each value from the bucket after we process it.
	// If I understand it correctly, not removing them does not hurt
	// the correctness, but hurts efficiency.
	// Instead of removing each value one by one, we empty the bucket after
	// iterating over it finishes.
	b.SemiDomBuckets[preoParent] = []uint{}
}

// "Step 4" of the Lengauer-Tarjan algorithm.
func (b *dominatorTreeBuilder) explicitlyDefineImmediateDominators() {
	n := b.Size()
	for preoCurrent := uint(1); preoCurrent < n; preoCurrent++ {
		if b.ImmDom[preoCurrent] != b.SemiDom[preoCurrent] {
			b.ImmDom[preoCurrent] = b.ImmDom[b.ImmDom[preoCurrent]]
		}
	}
}

// Converts the current internal representation of the builder into an outfacing
// DominatorTree type.
func (b *dominatorTreeBuilder) toDominatorTree() DominatorTree {
	n := b.Size()

	// Convert the internal immediate dominator mapping which uses preorder
	// addressing, to the original vertex indices.
	origImmDom := make([]uint, n)
	for i := uint(0); i < n; i++ {
		origImmDom[i] = b.PreorderToOriginal[b.ImmDom[b.OriginalToPreorder[i]]]
	}

	return DominatorTree{
		ImmDom: origImmDom,
	}
}

func (b *dominatorTreeBuilder) LengauerTarjan() DominatorTree {
	n := b.Size()

	// "Step 1" of the algorithm has been already computed "on the fly" when
	// we created the 'dominatorTreeBuilder' instance.

	for preoCurrent := uint(n - 1); preoCurrent > 0; preoCurrent-- {
		// "Step 2": Compute the semidominators of all vertices. Carry out the
		// computation vertex by vertex in decreasing preorder.
		semiDominator := b.calculateSemidominator(preoCurrent)

		// The Semidominator of the current vertex is now computed.
		// Add current vertex to it's semidominator bucket.
		// > "add w to bucket(vertex(semi(w)))"
		b.SemiDomBuckets[semiDominator] = append(
			b.SemiDomBuckets[semiDominator],
			preoCurrent,
		)

		// Link the current vertex to it's w.r.t the DFS tree, in the link-eval
		// forest.
		// > "LINK(parent(w), w)"
		preoParent := b.linkToDfsParent(preoCurrent)

		// "Step 3": Implicitly define the immediate dominator of each vertex.
		b.consumeSemiDominatorsBucket(preoParent)
	}

	// "Step 4": Explicitly define the immediate dominator of each vertex.
	b.explicitlyDefineImmediateDominators()

	// b.ImmDom now contains the immediate dominators of each vertex,
	// stored as a map from and to preorder indices.

	return b.toDominatorTree()
}
