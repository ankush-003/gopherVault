package btree

import (
	"encoding/binary"
)

const (
	BNODE_NODE         = 1 // internal nodes with pointers
	BNODE_LEAF         = 2 // leaf nodes with values
	BTREE_PAGE_SIZE    = 4096
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
)

type Node struct {
	keys [][]byte
	// vals are only used for leaf nodes
	vals [][]byte
	// children are only used for non-leaf nodes
	kids []*Node
}

type BNode []byte

// BNode Binary Format
// | type | nkeys | pointers   | offsets    | key-values | unused |
// | 2B   | 2B    | nkeys × 8B | nkeys × 2B | ...        |        |
//
// KV Format
// | key size | value size | key        | value     |
// | 2B       | 2B         | key        | value     |

type nodeOpts func(BNode) BNode

func WithHeader(btype, nkeys uint16) nodeOpts {
	return func(node BNode) BNode {
		node.setHeader(btype, nkeys)
		return node
	}
}

func NewBNode(opts ...nodeOpts) BNode {
	node := BNode(make([]byte, BTREE_PAGE_SIZE))

	for _, opt := range opts {
		node = opt(node)
	}

	return node
}

func (node BNode) btype() uint16 { // 2 bytes -> node type
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nkeys() uint16 { // 2 bytes -> number of keys
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) setHeader(btype, nkeys uint16) {
	if len(node) < 4 {
		panic("node too small")
	}
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

func (node BNode) getPtr(idx uint16) uint64 { // 8 bytes -> ptr to child
	if idx >= node.nkeys() {
		panic("index out of range")
	}
	pos := 4 + idx*8
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, ptr uint64) {
	if idx >= node.nkeys() {
		panic("index out of range")
	}
	pos := 4 + idx*8
	binary.LittleEndian.PutUint64(node[pos:], ptr)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	if idx > node.nkeys() {
		panic("index out of range")
	}
	pos := 4 + 8*node.nkeys() + (idx-1)*2
	return binary.LittleEndian.Uint16(node[pos:])
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	if idx > node.nkeys() {
		panic("index out of range")
	}

	pos := 4 + 8*node.nkeys() + (idx-1)*2
	binary.LittleEndian.PutUint16(node[pos:], offset)
}

func (node BNode) kvPos(idx uint16) uint16 {
	if idx > node.nkeys() {
		panic("index out of range")
	}
	pos := 4 + 8*node.nkeys() + 2*(node.nkeys()) + node.getOffset(idx)
	return pos
}

func (node BNode) getKey(idx uint16) []byte {
	if idx >= node.nkeys() {
		panic("index out of range")
	}
	pos := node.kvPos(idx)
	keySize := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:keySize]
}

func (node BNode) getVal(idx uint16) []byte {
	if idx >= node.nkeys() {
		panic("index out of range")
	}
	pos := node.kvPos(idx)
	keysize := binary.LittleEndian.Uint16(node[pos+0:])
	valSize := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+keysize:][:valSize]
}

func (node BNode) nbytes() uint16 {
	return node.getOffset(node.nkeys()) // offset of last key (0index at 1, 1index at n)
}

// appends a key-value pair to the node at the given index
func appendKV(node BNode, idx uint16, ptr uint64, key, val []byte) {
	if idx >= node.nkeys() {
		panic("index out of range")
	}

	node.setPtr(idx, ptr)

	pos := node.kvPos(idx)
	keySize := uint16(len(key))
	valSize := uint16(len(val))
	binary.LittleEndian.PutUint16(node[pos+0:], keySize)
	binary.LittleEndian.PutUint16(node[pos+2:], valSize)
	copy(node[pos+4:], key)
	copy(node[pos+4+keySize:], val)

	node.setOffset(idx+1, node.getOffset(idx)+4+keySize+valSize)
}

// appends n key-value pairs from old to new, n is the number of key-value pairs to append
func appendRange(new, old BNode, dst, src, n uint16) {
	for i := uint16(0); i < n; i++ {
		appendKV(new, dst+i,
			old.getPtr(src+i), old.getKey(src+i), old.getVal(src+i),
		)
	}
}
