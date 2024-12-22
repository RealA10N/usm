package graph_test

import (
	"slices"
	"testing"

	"alon.kr/x/usm/graph"
	"github.com/stretchr/testify/assert"
)

func TestSreedharGaoDJGraphExample(t *testing.T) {
	// Taken from the example provided in Figures 1 & 2 in the paper by Sreedhar
	// & Gao about the DJ-Graph data structure:
	// https://doi.org/10.1145/199448.199464

	cfg := SreedharGaoGraphExample()
	dj := cfg.DominatorJoinGraph(0)

	expectedJoinGraph := graph.NewGraph([][]uint{
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

	frontier := dj.DominatorFrontier(3)
	slices.Sort(frontier)
	assert.EqualValues(t, []uint{3, 15}, frontier)

	frontier = dj.DominatorFrontier(9)
	slices.Sort(frontier)
	assert.EqualValues(t, []uint{3, 15}, frontier)

	frontier = dj.DominatorFrontier(12)
	slices.Sort(frontier)
	assert.EqualValues(t, []uint{3, 12, 15}, frontier)
}

func TestSreedharGaoDominatorIteratedFrontierExample(t *testing.T) {
	// Examples taken the paper by Sreedhar & Gao about the DJ-Graph data
	// structure: https://doi.org/10.1145/199448.199464

	cfg := SreedharGaoGraphExample()
	dj := cfg.DominatorJoinGraph(0)

	// Section 3, example of Algorithm 3.1
	frontier := dj.IteratedDominatorFrontier([]uint{9, 12})
	slices.Sort(frontier)
	assert.EqualValues(t, []uint{3, 12, 15}, frontier)

	// Section 5, example of Algorithm 4.1 (Also presented in Figure 3)
	frontier = dj.IteratedDominatorFrontier([]uint{5, 13})
	slices.Sort(frontier)
	assert.EqualValues(t, []uint{2, 3, 8, 12, 15}, frontier)
}
