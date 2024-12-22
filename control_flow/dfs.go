package control_flow

type Dfs struct {
	// PreOrder[v] is the index of v in the DFS pre-order traversal.
	PreOrder []uint

	// PreOrderReversed[i] is the node with the i-th index in the preorder.
	PreOrderReversed []uint

	// PostOrder[v] is the index of v in the DFS post-order traversal.
	PostOrder []uint

	// Parent[v] is the parent of v in the DFS spanning tree.
	// The parent of the root of the tree is itself.
	Parent []uint

	// Depth[v] is the depth of v in the DFS spanning tree.
	Depth []uint

	// SubtreeSize[v] is the size of the subtree of v in the DFS spanning tree,
	// including v itself.
	SubtreeSize []uint
}

func (d *Dfs) IsAncestor(ancestor uint, descendant uint) bool {
	return d.PreOrder[ancestor] <= d.PreOrder[descendant] &&
		d.PostOrder[ancestor] >= d.PostOrder[descendant]
}

func (d *Dfs) IsStrictAncestor(ancestor uint, descendant uint) bool {
	return d.PreOrder[ancestor] < d.PreOrder[descendant] &&
		d.PostOrder[ancestor] > d.PostOrder[descendant]
}

func (d *Dfs) IsDeeper(deeper uint, shallower uint) bool {
	return d.Depth[deeper] >= d.Depth[shallower]
}

func (d *Dfs) IsStrictlyDeeper(deeper uint, shallower uint) bool {
	return d.Depth[deeper] > d.Depth[shallower]
}

func (d *Dfs) Subtree(v uint) []uint {
	start := d.PreOrder[v]
	end := start + d.SubtreeSize[v]
	return d.PreOrderReversed[start:end]
}
