package gossip

type Bloom struct {
	Keys       []uint64
	Bits       BitVec
	NumBitsSet uint64
}

type BitVec struct {
	Bits []uint64 `bin:"optional"`
	Len  uint64
}
