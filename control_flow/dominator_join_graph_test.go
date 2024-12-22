package control_flow_test

import (
	"testing"

	"alon.kr/x/usm/control_flow"
	"github.com/stretchr/testify/assert"
)

func TestSreedharGaoDJGraphExample(t *testing.T) {
	// Taken from the example provided in Figures 1 & 2 in the paper by Sreedhar
	// & Gao about the DJ-Graph data structure:
	// https://doi.org/10.1145/199448.199464

	cfg := SreedharGaoGraphExample()
	dj := cfg.DominatorJoinGraph(0)

	expectedJoinGraph := control_flow.NewGraph([][]uint{
		{},      // 0 (START)
		{},      // 1
		{4, 7},  // 2
		{},      // 3
		{},      // 4
		{},      // 5
		{2, 8},  // 6
		{8},     // 7
		{7, 15}, // 8
		{},      // 9
		{12},    // 10
		{12},    // 11
		{},      // 12
		{3, 15}, // 13
		{12},    // 14
		{16},    // 15
		{},      // 16 (END)
	})

	assert.True(t, dj.JoinGraph.Equal(&expectedJoinGraph))
}

func TestSreedharGaoDominatorFrontierExample(t *testing.T) {
	// Example taken from the example of Algorithm 3.1 in the paper by Sreedhar
	// & Gao about the DJ-Graph data structure:
	// https://doi.org/10.1145/199448.199464

	cfg := SreedharGaoGraphExample()
	dj := cfg.DominatorJoinGraph(0)

	assert.EqualValues(t, []uint{3, 15}, dj.DominatorFrontier(3))
	assert.EqualValues(t, []uint{3, 15}, dj.DominatorFrontier(9))
	assert.EqualValues(t, []uint{3, 12, 15}, dj.DominatorFrontier(12))
}
