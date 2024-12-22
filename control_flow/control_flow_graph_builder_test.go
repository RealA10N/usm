package control_flow_test

import (
	"testing"

	"alon.kr/x/usm/control_flow"
	"github.com/stretchr/testify/assert"
)

// TODO: the following tests relay on the current implementation and
// assume the order of the basic blocks in the graph, and the exact representation
// and implementation details (for example, comparison to nil slice and not
// empty slice). This is not ideal, and should be refactored to be more robust.

func TestBuildSingleton(t *testing.T) {
	g := control_flow.NewGraph(1, [][]uint{{}})
	cfg := g.ControlFlowGraph(0)

	assert.EqualValues(t, 1, cfg.Graph.Size())
	assert.Empty(t, cfg.Nodes[0].ForwardEdges)
	assert.Empty(t, cfg.Nodes[0].BackwardEdges)
}

func TestBuildSelfLoop(t *testing.T) {
	g := control_flow.NewGraph(1, [][]uint{{0}})
	cfg := g.ControlFlowGraph(0)

	assert.EqualValues(t, 1, cfg.Graph.Size())
	assert.EqualValues(t, []uint{0}, cfg.Nodes[0].ForwardEdges)
	assert.EqualValues(t, []uint{0}, cfg.Nodes[0].BackwardEdges)
}

func TestSingleBlock(t *testing.T) {
	// input graph:
	// 0 -> 1 -> 2
	//
	// control flow graph:
	// 0 (single node, no edges)

	g := control_flow.NewGraph(3, [][]uint{{1}, {2}, {}})

	cfg := g.ControlFlowGraph(0)
	assert.EqualValues(t, 1, cfg.Size())
	assert.EqualValues(t, []uint{0, 1, 2}, cfg.BasicBlockToNodes[0])
	assert.EqualValues(t, []uint{0, 0, 0}, cfg.NodeToBasicBlock)
	assert.Empty(t, cfg.Nodes[0].ForwardEdges)
	assert.Empty(t, cfg.Nodes[0].BackwardEdges)
}

func TestSimpleLoop(t *testing.T) {
	// input graph:
	// 0 -> 1 -> 2
	//      ^----+
	//
	// control flow graph:
	// 0 -> 1-+
	//      ^-+

	g := control_flow.NewGraph(3, [][]uint{{1}, {2}, {1}})

	cfg := g.ControlFlowGraph(0)
	assert.EqualValues(t, 2, cfg.Size())
	assert.EqualValues(t, [][]uint{{0}, {1, 2}}, cfg.BasicBlockToNodes)
	assert.EqualValues(t, []uint{0, 1, 1}, cfg.NodeToBasicBlock)

	assert.EqualExportedValues(t, control_flow.Node{
		ForwardEdges:  []uint{1},
		BackwardEdges: nil,
	}, cfg.Nodes[0])

	assert.EqualValues(t, control_flow.Node{
		ForwardEdges:  []uint{1},
		BackwardEdges: []uint{1, 0},
	}, cfg.Nodes[1])
}

func TestLoopToEntry(t *testing.T) {
	// input graph:
	// 0 -> 1 -> 2
	// ^---------+
	//
	// control flow graph:
	// 0-+  (single node, self loop)
	// ^-+

	g := control_flow.NewGraph(3, [][]uint{{1}, {2}, {0}})

	cfg := g.ControlFlowGraph(0)
	assert.EqualValues(t, cfg.Size(), 1)
	assert.EqualValues(t, [][]uint{{0, 1, 2}}, cfg.BasicBlockToNodes)
	assert.EqualValues(t, []uint{0, 0, 0}, cfg.NodeToBasicBlock)

	assert.EqualValues(t, control_flow.Node{
		ForwardEdges:  []uint{0},
		BackwardEdges: []uint{0},
	}, cfg.Nodes[0])
}

func TestLoopAndJumpOverLoop(t *testing.T) {
	// input graph:
	// 0 1 2 3
	// +-----^ (jump over loop)
	//   +-^   (loop)
	//   ^-+
	//
	// control flow graph:
	// 0 (singleton with no edges)

	g := control_flow.NewGraph(4, [][]uint{{3}, {2}, {1}, {}})
	cfg := g.ControlFlowGraph(0)

	assert.EqualValues(t, cfg.Size(), 1)
	assert.EqualValues(t, [][]uint{{0, 3}}, cfg.BasicBlockToNodes)
	assert.EqualValues(t, []uint{
		0,
		control_flow.Unreachable,
		control_flow.Unreachable,
		0,
	}, cfg.NodeToBasicBlock)

	assert.EqualValues(t, control_flow.Node{
		ForwardEdges:  nil,
		BackwardEdges: nil,
	}, cfg.Nodes[0])
}

func TestIfElse(t *testing.T) {
	// input graph:
	// 0 1 2 3
	// +-^     if
	// +---^   else
	//   +---^ continue
	//     +-^ continue
	//
	// control flow graph:
	//     0
	//  v--+--v
	//  1     3
	//  +--v--+
	//     2

	g := control_flow.NewGraph(4, [][]uint{{1, 2}, {3}, {3}, {}})

	cfg := g.ControlFlowGraph(0)
	assert.EqualValues(t, 4, cfg.Size())
	assert.EqualValues(t, [][]uint{{0}, {1}, {3}, {2}}, cfg.BasicBlockToNodes)
	assert.EqualValues(t, []uint{0, 1, 3, 2}, cfg.NodeToBasicBlock)

	assert.EqualValues(t, control_flow.Node{
		ForwardEdges:  []uint{1, 3},
		BackwardEdges: nil,
	}, cfg.Nodes[0])

	assert.EqualValues(t, control_flow.Node{
		ForwardEdges:  []uint{2},
		BackwardEdges: []uint{0},
	}, cfg.Nodes[1])

	assert.EqualValues(t, control_flow.Node{
		ForwardEdges:  nil,
		BackwardEdges: []uint{1, 3},
	}, cfg.Nodes[2])

	assert.EqualValues(t, control_flow.Node{
		ForwardEdges:  []uint{2},
		BackwardEdges: []uint{0},
	}, cfg.Nodes[3])
}
