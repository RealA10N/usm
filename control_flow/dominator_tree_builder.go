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
	SemiDom []uint
	Parent  []uint
}

func (b *dominatorTreeBuilder) isRoot(v uint) bool {
	return v == b.Parent[v]
}

func (b *dominatorTreeBuilder) isDirectAncestorOfRoot(v uint) bool {
	return b.isRoot(b.Parent[v])
}

// eval(v) returns v if and only if it is the root in a tree of the forest.
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
// merge it with its parent. Then, the second case in the eval definition
// applies, and we iterate over the ancestors of the node in the tree until
// reaching the root, returning the node with the minimal semidominator.
// On the way, we compress the path to the root of the tree, to speed up future
// evaluations.
//
// For reference, see Henrik Knakkegaard Christensen's master's thesis, section
// 3.3, page 27:
// https://users-cs.au.dk/gerth/advising/thesis/henrik-knakkegaard-christensen.pdf
func (b *dominatorTreeBuilder) eval(v uint) (minNode uint, root uint) {
	if b.isRoot(v) {
		return v, v
	}

	if b.isDirectAncestorOfRoot(v) {
		return v, b.Parent[v]
	}

	minNode, root = b.eval(b.Parent[v])
	b.Parent[v] = root
	if b.SemiDom[v] < b.SemiDom[minNode] {
		minNode = v
	}

	return minNode, root
}

func (b *dominatorTreeBuilder) link(p uint, v uint) {
	b.Parent[v] = p
}
