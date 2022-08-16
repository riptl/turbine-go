package gossip

import (
	"fmt"

	bin "github.com/gagliardetto/binary"
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

func (c *ContactInfo) DataID() uint32 {
	return DataIDContactInfo
}
