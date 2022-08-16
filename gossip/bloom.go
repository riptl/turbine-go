package gossip

import "golang.org/x/exp/constraints"

type Bloom struct {
	Keys       []uint64
	Bits       BitVec[uint64]
	NumBitsSet uint64
}

type BitVec[T constraints.Unsigned] struct {
	Bits []T `bin:"optional"`
	Len  uint64
}
