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

func (p *PullRequest) MsgID() uint32 {
	return MsgIDPullRequest
}
