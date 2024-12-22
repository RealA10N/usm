package graph_test

import (
	"testing"

	"alon.kr/x/usm/graph"
	"github.com/stretchr/testify/assert"
)

func TestSimpleGraphEquality(t *testing.T) {
	g1 := graph.NewGraph([][]uint{{1, 2}, {}, {}})
	g2 := graph.NewGraph([][]uint{{2, 1}, {}, {}})
	assert.True(t, g1.Equal(&g2))
	assert.True(t, g2.Equal(&g1))
}
