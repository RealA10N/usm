package control_flow

type Dfs struct {
	// PreOrder[v] is the index of v in the DFS pre-order traversal.
	PreOrder []uint

	// PostOrder[v] is the index of v in the DFS post-order traversal.
	PostOrder []uint

	// Parent[v] is the parent of v in the DFS spanning tree.
	// The parent of the root of the tree is itself.
	Parent []uint

	// Depth[v] is the depth of v in the DFS spanning tree.
	Depth []uint
}

func (d *Dfs) IsAncestor(ancestor uint, descendant uint) bool {
	return d.PreOrder[ancestor] <= d.PreOrder[descendant] &&
		d.PostOrder[ancestor] >= d.PostOrder[descendant]
}

func (d *Dfs) IsStrictAncestor(ancestor uint, descendant uint) bool {
	return d.PreOrder[ancestor] < d.PreOrder[descendant] &&
		d.PostOrder[ancestor] > d.PostOrder[descendant]
}
