package gossip

import (
	"fmt"
	"net/netip"

	bin "github.com/gagliardetto/binary"
)

type SocketAddr struct {
	netip.AddrPort
}

func (s *SocketAddr) UnmarshalWithDecoder(d *bin.Decoder) (err error) {
	ipType, err := d.ReadUint32(bin.LE)
	if err != nil {
		return err
	}
	var ipBytes []byte
	switch ipType {
	case 0:
		ipBytes, err = d.ReadNBytes(4)
	case 1:
		ipBytes, err = d.ReadNBytes(16)
	default:
		err = fmt.Errorf("invalid SocketAddr type %#x", ipType)
	}
	if err != nil {
		return err
	}
	ipAddr, _ := netip.AddrFromSlice(ipBytes)
	if ipAddr == netip.AddrFrom4([4]byte{0, 0, 0, 0}) {
		// All zero IP serves as a placeholder
		ipAddr = netip.Addr{}
	}
	port, err := d.ReadUint16(bin.LE)
	if err != nil {
		return err
	}
	s.AddrPort = netip.AddrPortFrom(ipAddr, port)
	return nil
}

func (s *SocketAddr) MarshalWithEncoder(e *bin.Encoder) (err error) {
	ipBytes := s.Addr().AsSlice()
	var ipType uint32
	switch len(ipBytes) {
	case 4:
		ipType = 0
	case 16:
		ipType = 1
	default:
		return fmt.Errorf("silly length SocketAddr address: %d", len(ipBytes))
	}
	if err = e.WriteUint32(ipType, bin.LE); err != nil {
		return err
	}
	if err = e.WriteBytes(ipBytes, false); err != nil {
		return err
	}
	return e.WriteUint16(s.Port(), bin.LE)
}

const (
	SlotsVecIDFlate2 = uint32(iota) // TODO rename, "flate2" is a crate, not a compression format
	SlotsVecIDUncompressed
)

type SlotsVec interface {
	SlotsVecID() uint32
}

type SlotsVecEnum struct {
	SlotsVecID uint32
	SlotsVec
}

type SlotsVecFlate2 struct {
	FirstSlot  uint64
	Num        uint64
	Compressed []byte
}

func (s *SlotsVecFlate2) SlotsVecID() uint32 {
	return SlotsVecIDFlate2
}

type SlotsVecUncompressed struct {
	FirstSlot uint64
	Num       uint64
	Slots     BitVec[uint8]
}

type HashEvent struct {
	Slot uint64
	Hash [32]byte
}
