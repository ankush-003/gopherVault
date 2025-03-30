package btree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSize(t *testing.T) {
	node1Max := 4 + 1*8 + 1*2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	t.Log("node1Max:", node1Max)
	assert.GreaterOrEqual(t, BTREE_PAGE_SIZE, node1Max)
}

func TestNewBNode(t *testing.T) {
	node := NewBNode(WithHeader(BNODE_LEAF, 5))

	k1, v1 := []byte("k1"), []byte("hi")
	k2, v2 := []byte("k2"), []byte("hello")

	appendKV(node, 0, 0, k1, v1)
	appendKV(node, 1, 0, k2, v2)

	t.Log("node:", node[:80])

	t.Log("k1:", k1, "v1:", v1)
	t.Log("k2:", k2, "v2:", v2)

	assert.Equal(t, k1, node.getKey(0))
	assert.Equal(t, v1, node.getVal(0))
	assert.Equal(t, k2, node.getKey(1))
	assert.Equal(t, v2, node.getVal(1))
}
