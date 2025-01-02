// This file contains the implementation of the dominator tree data structure.
// The construction of the data structure is implemented in
// `dominator_tree_builder.go`.

package graph

type DominatorTree struct {
	*Graph
	Dfs

	// ImmDom[v] is the immediate dominator of v, and by definition of the
	// dominator tree, ImmDom[v] is the parent of v in the dominator tree.
	//
	// It is assumed that ImmDom[root] = root.
	ImmDom []uint
}

func (t *DominatorTree) IsDominatorOf(dominator uint, dominated uint) bool {
	return t.IsAncestor(dominator, dominated)
}

func (t *DominatorTree) IsStrictDominatorOf(dominator uint, dominated uint) bool {
	return t.IsStrictAncestor(dominator, dominated)
}
