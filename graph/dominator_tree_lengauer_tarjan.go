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
//
// TODO: It seems like there is a linear time construction algorithm for
// dominator trees:
// https://doi.org/10.1145/22145.22166

package graph

// A set of variables that the algorithm uses in it's runtime.
// For code readability and performance, instead of defining them as local
// variables and passing them around, we define a struct containing them all
// and a set of methods that manipulate it.
type lengauerTarjanContext struct {
	*Graph
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
	// The Original -> Preorder mapping is stored in Dfs.Preorder.
	// The Preorder -> Original mapping is stored in Dfs.PreorderReversed.
	//
	// I admit, this is a bit confusing.
	Dfs
}

func newLengauerTarjanContext(g *Graph, entry uint) lengauerTarjanContext {
	n := g.Size()

	// "Step 1" of the Lengauer-Tarjan algorithm: compute and store the DFS
	// tree, and especially the preorder traversal.
	dfs := g.Dfs(entry)

	builder := lengauerTarjanContext{
		Graph:          g,
		LinkEvalForest: NewLinkEvalForest(n),
		ImmDom:         make([]uint, n),
		SemiDomBuckets: make([][]uint, n),
		Dfs:            dfs,
	}

	return builder
}

// "Step 2" of the Lengauer-Tarjan algorithm, as it is described in the original
// paper.
//
// Implementation is mainly based on the description of Knakkegaard's Thesis,
// section 3.3.
func (c *lengauerTarjanContext) calculateSemidominator(preoCurrent uint) uint {
	origCurrent := c.PreOrderReversed[preoCurrent]
	currentBlock := c.Nodes[origCurrent]

	for _, origPredecessor := range currentBlock.BackwardEdges {
		preoPredecessor := c.PreOrder[origPredecessor]

		candidate := c.LinkEvalForest.Eval(preoPredecessor)
		candidateSemidom := c.SemiDom[candidate]

		mySemidom := c.SemiDom[preoCurrent]
		if candidateSemidom < mySemidom {
			c.SemiDom[preoCurrent] = candidateSemidom
		}
	}

	return c.SemiDom[preoCurrent]
}

func (c *lengauerTarjanContext) linkToDfsParent(preoCurrent uint) uint {
	origCurrent := c.PreOrderReversed[preoCurrent]
	origDfsParent := c.Dfs.Parent[origCurrent]
	preoDfsParent := c.PreOrder[origDfsParent]
	c.Link(preoCurrent, preoDfsParent)
	return preoDfsParent
}

// "Step 3" of the Lengauer-Tarjan algorithm.
func (c *lengauerTarjanContext) consumeSemiDominatorsBucket(preoParent uint) {
	for _, v := range c.SemiDomBuckets[preoParent] {
		u := c.Eval(v)
		if c.SemiDom[u] < c.SemiDom[v] {
			c.ImmDom[v] = u
		} else {
			c.ImmDom[v] = c.SemiDom[v]
		}
	}

	// The paper explicitly states that in the loop above, we should remove
	// each value from the bucket after we process it.
	// If I understand it correctly, not removing them does not hurt
	// the correctness, but hurts efficiency.
	// Instead of removing each value one by one, we empty the bucket after
	// iterating over it finishes.
	c.SemiDomBuckets[preoParent] = []uint{}
}

// "Step 4" of the Lengauer-Tarjan algorithm.
func (c *lengauerTarjanContext) explicitlyDefineImmediateDominators() {
	n := c.Size()
	for preoCurrent := uint(1); preoCurrent < n; preoCurrent++ {
		if c.ImmDom[preoCurrent] != c.SemiDom[preoCurrent] {
			c.ImmDom[preoCurrent] = c.ImmDom[c.ImmDom[preoCurrent]]
		}
	}
}

// Converts the current internal representation of the builder into an outfacing
// DominatorTree type.
func (c *lengauerTarjanContext) getOriginalImmediateDominators() []uint {
	n := c.Size()

	// Convert the internal immediate dominator mapping which uses preorder
	// addressing, to the original vertex indices.
	origImmDom := make([]uint, n)
	for i := uint(0); i < n; i++ {
		origImmDom[i] = c.PreOrderReversed[c.ImmDom[c.PreOrder[i]]]
	}

	return origImmDom
}

func (c *lengauerTarjanContext) LengauerTarjan() []uint {
	n := c.Size()

	// "Step 1" of the algorithm has been already computed "on the fly" when
	// we created the 'dominatorTreeBuilder' instance.

	for preoCurrent := uint(n - 1); preoCurrent > 0; preoCurrent-- {
		// "Step 2": Compute the semidominators of all vertices. Carry out the
		// computation vertex by vertex in decreasing preorder.
		semiDominator := c.calculateSemidominator(preoCurrent)

		// The Semidominator of the current vertex is now computed.
		// Add current vertex to it's semidominator bucket.
		// > "add w to bucket(vertex(semi(w)))"
		c.SemiDomBuckets[semiDominator] = append(
			c.SemiDomBuckets[semiDominator],
			preoCurrent,
		)

		// Link the current vertex to it's w.r.t the DFS tree, in the link-eval
		// forest.
		// > "LINK(parent(w), w)"
		preoParent := c.linkToDfsParent(preoCurrent)

		// "Step 3": Implicitly define the immediate dominator of each vertex.
		c.consumeSemiDominatorsBucket(preoParent)
	}

	// "Step 4": Explicitly define the immediate dominator of each vertex.
	c.explicitlyDefineImmediateDominators()

	// b.ImmDom now contains the immediate dominators of each vertex,
	// stored as a map from and to preorder indices.
	// We only need to transform the representation back to use the original
	// indices.

	return c.getOriginalImmediateDominators()
}
