package control_flow_test

import (
	"testing"

	"alon.kr/x/usm/control_flow"
	"github.com/stretchr/testify/assert"
)

func TestLineDfs(t *testing.T) {
	// 0 -> 1 -> 2

	g := control_flow.NewGraph(3, [][]uint{{1}, {2}, {}})
	dfs := g.Dfs(0)

	assert.EqualValues(t, []uint{0, 1, 2}, dfs.PreOrder)
	assert.EqualValues(t, []uint{2, 1, 0}, dfs.PostOrder)
	assert.EqualValues(t, []uint{0, 0, 1}, dfs.Parent)
	assert.EqualValues(t, []uint{0, 1, 2}, dfs.Depth)
}

func TestBinaryTree(t *testing.T) {
	//    0
	//  1   2
	// 3 4 5 6

	g := control_flow.NewGraph(7, [][]uint{{1, 2}, {3, 4}, {5, 6}, {}, {}, {}, {}})
	dfs := g.Dfs(0)

	//                    index: 0  1  2  3  4  5  6
	assert.EqualValues(t, []uint{0, 1, 4, 2, 3, 5, 6}, dfs.PreOrder)
	assert.EqualValues(t, []uint{6, 2, 5, 0, 1, 3, 4}, dfs.PostOrder)
	assert.EqualValues(t, []uint{0, 0, 0, 1, 1, 2, 2}, dfs.Parent)
	assert.EqualValues(t, []uint{0, 1, 1, 2, 2, 2, 2}, dfs.Depth)
}

func TestSpecialEdges(t *testing.T) {
	//    0      with an additional cross edge 5 -> 1
	//  1   2    a back edge 4 -> 0
	// 3 4 5 6   and a forward edge 0 -> 6

	g := control_flow.NewGraph(7, [][]uint{
		{1, 2, 6}, {3, 4}, {5, 6}, {}, {0}, {1}, {},
	})
	dfs := g.Dfs(0)

	//                    index: 0  1  2  3  4  5  6
	assert.EqualValues(t, []uint{0, 1, 4, 2, 3, 5, 6}, dfs.PreOrder)
	assert.EqualValues(t, []uint{6, 2, 5, 0, 1, 3, 4}, dfs.PostOrder)
	assert.EqualValues(t, []uint{0, 0, 0, 1, 1, 2, 2}, dfs.Parent)
	assert.EqualValues(t, []uint{0, 1, 1, 2, 2, 2, 2}, dfs.Depth)
}
