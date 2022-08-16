package gossip

import (
	"fmt"

	bin "github.com/gagliardetto/binary"
)

const (
	MsgIDPullRequest = uint32(iota)
	MsgIDPullResponse
	MsgIDPush
	MsgIDPrune
	MsgIDPing
	MsgIDPong
)

type Message interface {
	MsgID() uint32
}

type MessageEnum struct {
	MsgID uint32
	Message
}

func (m *MessageEnum) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	if m.MsgID, err = dec.ReadUint32(bin.LE); err != nil {
		return err
	}
	switch m.MsgID {
	case MsgIDPullRequest:
		m.Message = new(PullRequest)
	case MsgIDPullResponse:
		m.Message = new(PullResponse)
	case MsgIDPush:
		m.Message = new(PushMessage)
	case MsgIDPrune:
		m.Message = new(PruneMessage)
	case MsgIDPing:
		m.Message = new(PingMessage)
	case MsgIDPong:
		m.Message = new(PongMessage)
	default:
		return fmt.Errorf("unsupported message type %#x", m.MsgID)
	}
	return dec.Decode(m.Message)
}

func UnmarshalMessage(data []byte) (Message, error) {
	var msg MessageEnum
	dec := bin.NewBinDecoder(data)
	err := dec.Decode(&msg)
	if err == nil && dec.HasRemaining() {
		return nil, fmt.Errorf("unexpected %d bytes past end of message", dec.Remaining())
	}
	return msg.Message, err
}

type PullRequest struct {
	DataFilter
	Value
}

func (*PullRequest) MsgID() uint32 {
	return MsgIDPullRequest
}

type PullResponse struct {
	Pubkey [32]byte
	Values []Value
}

func (*PullResponse) MsgID() uint32 {
	return MsgIDPullResponse
}

type PushMessage struct {
	Pubkey [32]byte
	Values []Value
}

func (*PushMessage) MsgID() uint32 {
	return MsgIDPush
}

type PruneMessage struct {
	Pubkey0     [32]byte
	Pubkey1     [32]byte
	Prunes      [][32]byte
	Signature   [64]byte
	Destination [32]byte
	Wallclock   uint64
}

func (*PruneMessage) MsgID() uint32 {
	return MsgIDPrune
}

type PingMessage struct {
	Ping
}

func (*PingMessage) MsgID() uint32 {
	return MsgIDPing
}

type PongMessage struct {
	Ping
}

func (*PongMessage) MsgID() uint32 {
	return MsgIDPong
}
