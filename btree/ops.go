package btree

func leafInsert(new, old BNode, idx uint16, key, val []byte) {
	if idx >= old.nkeys() {
		panic("index out of range")
	}
	appendRange(new, old, 0, 0, idx)
	appendKV(new, idx, 0, key, val)
	appendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

func leafUpdate(new, old BNode, idx uint16, key, val []byte) {
	if idx >= old.nkeys() {
		panic("index out of range")
	}
	appendRange(new, old, 0, 0, idx)
	appendKV(new, idx, 0, key, val)
	appendRange(new, old, idx+1, idx+1, old.nkeys()-idx-1)
}
