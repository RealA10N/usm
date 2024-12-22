// This file implements the link-eval data structure, used in the Lengauer-Tarjan
// algorithm for dominator tree construction.
//
// Some details of the data structure are explained the Lengauer-Tarjan paper:
// https://doi.org/10.1145/357062.357071
// especially in section 4.
//
// I found Henrik Knakkegaard Christensen's master's thesis on "Algorithms for
// Finding Dominators in Directed Graphs" to be a helpful resource:
// https://users-cs.au.dk/gerth/advising/thesis/henrik-knakkegaard-christensen.pdf
// especially in section 3.3.
//
// The original paper (Tarjan) on the data structure is:
// https://doi.org/10.1145/322154.322161

package graph

// Nodes are assumed to be numbered [0, n), where n is the number of nodes.
type LinkEvalForest struct {
	// SemiDom[v] is the semidominator of node v in the forest.
	// It is to be set by the user of the data structure.
	// Initially, SemiDom[v] = v for all v.
	SemiDom []uint

	// Parent[v] is the direct parent of node v in the forest.
	// If v is the root of the tree it is in, Parent[v] == v.
	// The data structure is initialized with Parent[v] = v for all v.
	Parent []uint
}

func NewLinkEvalForest(n uint) LinkEvalForest {
	f := LinkEvalForest{
		SemiDom: make([]uint, n),
		Parent:  make([]uint, n),
	}

	for i := uint(0); i < n; i++ {
		f.Parent[i] = i
		f.SemiDom[i] = i
	}

	return f
}

// MARK: Public

// Eval(v) returns v if and only if it is the root in a tree of the forest.
// Otherwise, it returns a node u in the path from v to the root of the tree
// that contains v (not including the root itself), such that Semidominator(u)
// is minimal.
//
// Since we are calculating the semidominators of vertices in reverse order,
// if evaluating a node that has yet to be processed and given a semidominator,
// the tree containing that node will be a singleton and thus it will be the
// root, and the node itself will be returned by definition.
//
// If however, we evaluate a node that has already been processed, it must not
// be the root of the tree it is in, since right after processing a node we
// merge it with its parent. Then, the second case in the Eval definition
// applies, and we iterate over the ancestors of the node in the tree until
// reaching the root, returning the node with the minimal semidominator.
// On the way, we compress the path to the root of the tree, to speed up future
// evaluations.
//
// For reference, see Henrik Knakkegaard Christensen's master's thesis, section
// 3.3, page 27:
// https://users-cs.au.dk/gerth/advising/thesis/henrik-knakkegaard-christensen.pdf
func (f *LinkEvalForest) Eval(v uint) uint {
	if f.isRoot(v) {
		return v
	}

	initial := v
	minSemiDom := f.SemiDom[v]
	minValue := v
	for !f.isRootDirectAncestor(v) {
		v = f.Parent[v]
		if f.SemiDom[v] < minSemiDom {
			minSemiDom = f.SemiDom[v]
			minValue = v
		}
	}

	root := f.Parent[v]
	f.compress(initial, root)

	return minValue
}

func (b *LinkEvalForest) Link(child uint, parent uint) {
	b.Parent[child] = parent
}

// MARK: Private

// Returns true iff v is the root of the tree it is in.
func (f *LinkEvalForest) isRoot(v uint) bool {
	return v == f.Parent[v]
}

// Returns true if the root is a direct ancestor of v.
// Assumes that v is not a root.
func (f *LinkEvalForest) isRootDirectAncestor(v uint) bool {
	return f.isRoot(f.Parent[v])
}

// Compresses the path from v to the root of the tree it is in,
// by linking all nodes on the path to the root directly to the provided
// root node.
//
// There is no validation that the provided root is actually the root of the
// node's tree.
func (f *LinkEvalForest) compress(v uint, root uint) {
	for !f.isRoot(v) {
		f.Parent[v] = root
		v = f.Parent[v]
	}
}
