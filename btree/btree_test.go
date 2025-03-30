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
	node := NewBNode(WithHeader(BNODE_LEAF, 2))

	k1, v1 := []byte("k1"), []byte("hi")
	k2, v2 := []byte("k2"), []byte(`{"name":"joe"}`)

	appendKV(node, 0, 0, k1, v1)
	appendKV(node, 1, 0, k2, v2)

	t.Log("node:", node[:80])

	t.Log("k1:", k1, "v1:", v1)
	t.Log("k2:", k2, "v2:", v2)
	t.Log("node size:", node.nbytes())

	assert.Equal(t, k1, node.getKey(0))
	assert.Equal(t, v1, node.getVal(0))
	assert.Equal(t, k2, node.getKey(1))
	assert.Equal(t, v2, node.getVal(1))
}

func TestLeafInsert(t *testing.T) {
	node := NewBNode(WithHeader(BNODE_LEAF, 2))

	k1, v1 := []byte("k1"), []byte("hi")
	k2, v2 := []byte("k2"), []byte(`{"name":"joe"}`)

	appendKV(node, 0, 0, k1, v1)
	appendKV(node, 1, 0, k2, v2)

	t.Log("node:", node[:80])
	t.Log("k1:", k1, "v1:", v1)
	t.Log("k2:", k2, "v2:", v2)
	t.Log("node size:", node.nbytes())

	new := NewBNode(WithHeader(BNODE_LEAF, 3))
	k3, v3 := []byte("k3"), []byte("hello")
	leafInsert(new, node, 1, k3, v3)
	t.Log("new:", new[:80])
	t.Log("new size:", new.nbytes())
	assert.Equal(t, k1, new.getKey(0))
	assert.Equal(t, v1, new.getVal(0))
	assert.Equal(t, k3, new.getKey(1))
	assert.Equal(t, v3, new.getVal(1))
	assert.Equal(t, k2, new.getKey(2))
	assert.Equal(t, v2, new.getVal(2))

	new_update := NewBNode(WithHeader(BNODE_LEAF, 3))
	leafUpdate(new_update, new, 2, k3, []byte("hello world"))
	t.Log("new_update:", new_update[:80])
	assert.Equal(t, k3, new_update.getKey(2))
	assert.Equal(t, []byte("hello world"), new_update.getVal(2))
}
