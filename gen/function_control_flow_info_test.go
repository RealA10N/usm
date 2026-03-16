package gen_test

import (
	"testing"

	"alon.kr/x/usm/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFunctionControlFlowInfo_SingleBlock checks that a function whose only
// instruction is a return produces a CFG with one node and no forward edges.
func TestNewFunctionControlFlowInfo_SingleBlock(t *testing.T) {
	src := `func $32 @f $32 %n {
.entry
	ret
}
`
	fn, results := generateFunctionFromSource(t, src)
	require.True(t, results.IsEmpty())

	info := gen.NewFunctionControlFlowInfo(fn)

	assert.Len(t, info.BasicBlocks, 1)
	assert.Len(t, info.BasicBlocksToIndex, 1)
	require.NotNil(t, info.ControlFlowGraph)
	assert.Len(t, info.ControlFlowGraph.Nodes, 1)
	assert.Empty(t, info.ControlFlowGraph.Nodes[0].ForwardEdges)
}

// TestNewFunctionControlFlowInfo_BlockToIndexMapping checks that every block
// in BasicBlocks has an entry in BasicBlocksToIndex whose value equals its
// position in the slice.
func TestNewFunctionControlFlowInfo_BlockToIndexMapping(t *testing.T) {
	// jz creates a fall-through block AND a branch target block, so we get
	// three blocks: entry, fall-through (ret), and .end (ret).
	src := `func $32 @f $32 %n {
.entry
	jz %n .end
	ret
.end
	ret
}
`
	fn, results := generateFunctionFromSource(t, src)
	require.True(t, results.IsEmpty())

	info := gen.NewFunctionControlFlowInfo(fn)

	require.Equal(t, len(info.BasicBlocks), len(info.BasicBlocksToIndex))
	for i, block := range info.BasicBlocks {
		idx, ok := info.BasicBlocksToIndex[block]
		assert.True(t, ok, "block at position %d missing from BasicBlocksToIndex", i)
		assert.Equal(t, uint(i), idx, "block at position %d has wrong index", i)
	}
}

// TestNewFunctionControlFlowInfo_UnconditionalJump checks that an
// unconditional jump (j) produces exactly one forward edge from the jumping
// block to its target.
func TestNewFunctionControlFlowInfo_UnconditionalJump(t *testing.T) {
	src := `func $32 @f $32 %n {
.entry
	j .end
.end
	ret
}
`
	fn, results := generateFunctionFromSource(t, src)
	require.True(t, results.IsEmpty())

	info := gen.NewFunctionControlFlowInfo(fn)

	require.Len(t, info.BasicBlocks, 2)
	require.NotNil(t, info.ControlFlowGraph)

	// Block 0 (.entry) must have exactly one forward edge pointing to block 1 (.end).
	entryEdges := info.ControlFlowGraph.Nodes[0].ForwardEdges
	require.Len(t, entryEdges, 1)
	assert.Equal(t, uint(1), entryEdges[0])

	// Block 1 (.end) has no successors.
	assert.Empty(t, info.ControlFlowGraph.Nodes[1].ForwardEdges)
}

// TestNewFunctionControlFlowInfo_ConditionalBranch checks that a conditional
// branch (jz) produces two forward edges: one fall-through and one branch.
// The jz instruction causes an implicit block split, so there are three blocks:
//
//	block 0 – .entry containing "jz %n .end"
//	block 1 – implicit fall-through block containing "ret"
//	block 2 – .end containing "ret"
func TestNewFunctionControlFlowInfo_ConditionalBranch(t *testing.T) {
	src := `func $32 @f $32 %n {
.entry
	jz %n .end
	ret
.end
	ret
}
`
	fn, results := generateFunctionFromSource(t, src)
	require.True(t, results.IsEmpty())

	info := gen.NewFunctionControlFlowInfo(fn)

	require.Len(t, info.BasicBlocks, 3)
	require.NotNil(t, info.ControlFlowGraph)

	// The entry block should have two forward edges: fall-through (1) and
	// branch target (2).
	entryEdges := info.ControlFlowGraph.Nodes[0].ForwardEdges
	assert.Len(t, entryEdges, 2)

	n := uint(len(info.BasicBlocks))
	for _, target := range entryEdges {
		assert.Less(t, target, n, "forward edge index %d out of range", target)
	}

	// The two leaf blocks (fall-through ret and .end ret) have no successors.
	assert.Empty(t, info.ControlFlowGraph.Nodes[1].ForwardEdges)
	assert.Empty(t, info.ControlFlowGraph.Nodes[2].ForwardEdges)
}
