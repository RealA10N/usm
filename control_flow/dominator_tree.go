// This file contains the implementation of the dominator tree data structure.
// The construction of the data structure is implemented in
// `dominator_tree_builder.go`.

package control_flow

type DominatorTree[InstT SupportsControlFlow] struct {
	ControlFlowGraph ControlFlowGraph[InstT]

	// ImmDom[node] is the immediate dominator of the node `node`.
	// It is assumed that ImmDom[entryNode] = entryNode.
	ImmDom []uint

	// InTime[node] is the location of the node in a pre-order traversal of the
	// DFS tree. It is assumed to be a number in [0, n).
	InTime []uint

	// InTime[node] is the index of the node in a post-order traversal of the
	// DFS tree. It is assumed to be a number in [0, n).
	OutTime []uint
}

func (t *DominatorTree[InstT]) IsDominatorOf(dominator uint, dominated uint) bool {
	return (t.InTime[dominator] <= t.InTime[dominated] &&
		t.OutTime[dominator] >= t.OutTime[dominated])
}

func (t *DominatorTree[InstT]) IsStrictDominatorOf(dominator uint, dominated uint) bool {
	return (t.InTime[dominator] < t.InTime[dominated] &&
		t.OutTime[dominator] > t.OutTime[dominated])
}

func (t *DominatorTree[InstT]) Dominators(node uint) []uint {
	dominators := []uint{}
	for ; node != entryNode; node = t.ImmDom[node] {
		dominators = append(dominators, node)
	}
	return dominators
}

func (t *DominatorTree[InstT]) StrictDominators(node uint) []uint {
	dominators := []uint{}

	// It is OK to not check here if node == entryNode since we assume that
	// ImmDom[entryNode] = entryNode.
	node = t.ImmDom[node]

	for ; node != entryNode; node = t.ImmDom[node] {
		dominators = append(dominators, node)
	}
	return dominators
}
