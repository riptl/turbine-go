package gossip

import (
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

const (
	DataIDContactInfo = uint32(iota)
	DataIDVote
	DataIDLowestSlot
	DataIDSnapshotHashes
	DataIDAccountsHashes
	DataIDEpochSlots
	DataIDLegacyVersion
	DataIDVersion
	DataIDNodeInstance
	DataIDDuplicateShred
	DataIDIncrementalSnapshotHashes
)

type DataFilter struct {
	Filter   Bloom
	Mask     uint64
	MaskBits uint32
}

type Value struct {
	Signature [64]byte
	Data      DataEnum
}

type Data interface {
	DataID() uint32
}

type DataEnum struct {
	DataID uint32
	Data   Data
}

func (d *DataEnum) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	if d.DataID, err = dec.ReadUint32(bin.LE); err != nil {
		return err
	}
	switch d.DataID {
	case DataIDContactInfo:
		d.Data = new(ContactInfo)
	case DataIDVote:
		d.Data = new(VoteData)
	case DataIDLowestSlot:
		d.Data = new(LowestSlotData)
	case DataIDSnapshotHashes:
		d.Data = new(SnapshotHashes)
	case DataIDAccountsHashes:
		d.Data = new(AccountsHashes)
	case DataIDEpochSlots:
		d.Data = new(EpochSlots)
	case DataIDLegacyVersion:
		d.Data = new(LegacyVersion)
	case DataIDVersion:
		d.Data = new(Version)
	case DataIDNodeInstance:
		d.Data = new(NodeInstance)
	case DataIDDuplicateShred:
		d.Data = new(DuplicateShredData)
	case DataIDIncrementalSnapshotHashes:
		d.Data = new(IncrementalSnapshotHashes)
	default:
		return fmt.Errorf("unsupported data type %#x", d.DataID)
	}
	return dec.Decode(d.Data)
}

type ContactInfo struct {
	ID           [32]byte
	Gossip       SocketAddr
	TVU          SocketAddr
	TVUFwd       SocketAddr
	Repair       SocketAddr
	TPU          SocketAddr
	TPUFwd       SocketAddr
	TPUVote      SocketAddr
	RPC          SocketAddr
	RPCPubSub    SocketAddr
	ServeRepair  SocketAddr
	Wallclock    uint64
	ShredVersion uint16
}

func (*ContactInfo) DataID() uint32 {
	return DataIDContactInfo
}

type VoteData struct {
	Index uint8
	Vote
}

type Vote struct {
	From        [32]byte
	Transaction solana.Transaction
	Wallclock   uint64
	Slot        uint64 `bin:"optional"`
}

func (*VoteData) DataID() uint32 {
	return DataIDVote
}

type LowestSlotData struct {
	Index uint8 // deprecated
	LowestSlot
}

type LowestSlot struct {
	From      [32]byte
	Root      uint64 // deprecated
	Lowest    uint64
	Slots     []uint64               // sorted
	Stash     []EpochIncompleteSlots // deprecated
	Wallclock uint64
}

type EpochIncompleteSlots struct {
	First           uint64
	CompressionType uint32
	CompressedList  []byte
}

func (*LowestSlotData) DataID() uint32 {
	return DataIDLowestSlot
}

type SnapshotHashes struct {
	From      [32]byte
	Hashes    []HashEvent
	Wallclock uint64
}

func (*SnapshotHashes) DataID() uint32 {
	return DataIDSnapshotHashes
}

type AccountsHashes struct {
	From      [32]byte
	Hashes    []HashEvent
	Wallclock uint64
}

func (*AccountsHashes) DataID() uint32 {
	return DataIDAccountsHashes
}

type EpochSlots struct {
	Index     uint8
	From      [32]byte
	Slots     []SlotsVecEnum
	Wallclock uint64
}

func (*EpochSlots) DataID() uint32 {
	return DataIDEpochSlots
}

type LegacyVersion struct {
	From      [32]byte
	Wallclock uint64
	Major     uint16
	Minor     uint16
	Patch     uint16
	Commit    uint32 `bin:"optional"`
}

func (v *LegacyVersion) DataID() uint32 {
	return DataIDLegacyVersion
}

type Version struct {
	From      [32]byte
	Wallclock uint64
	Major     uint16
	Minor     uint16
	Patch     uint16
	Commit    uint32 `bin:"optional"`
	Features  uint32
}

func (v *Version) DataID() uint32 {
	return DataIDVersion
}

type NodeInstance struct {
	From      [32]byte
	Wallclock uint64
	Timestamp uint64
	Token     uint64
}

func (i *NodeInstance) DataID() uint32 {
	return DataIDNodeInstance
}

type DuplicateShredData struct {
	Index uint16
	DuplicateShred
}

type DuplicateShred struct {
	From       [32]byte
	Wallclock  uint64
	Slot       uint64
	ShredIndex uint32
	ShredType  uint8
	NumChunks  uint8
	ChunkIndex uint8
	Chunk      []byte
}

func (d *DuplicateShredData) DataID() uint32 {
	return DataIDDuplicateShred
}

type IncrementalSnapshotHashes struct {
	From      [32]byte
	Base      HashEvent
	Hashes    []HashEvent
	Wallclock uint64
}

func (h *IncrementalSnapshotHashes) DataID() uint32 {
	return DataIDIncrementalSnapshotHashes
}
