package gossip

import (
	"crypto/ed25519"
)

type Ping struct {
	From      [32]byte
	Token     [32]byte
	Signature [64]byte
}

const PackedIDPing uint32 = 0x04
const PacketIDPong uint32 = 0x05

// PingSize is the size of a serialized ping message.
const PingSize = 128

// SignPing creates a new ping message.
//
// Panics if the provided private key is invalid.
func SignPing(token [32]byte, key ed25519.PrivateKey) *Ping {
	sig := ed25519.Sign(key, token[:])
	p := new(Ping)
	copy(p.From[:], key.Public().(ed25519.PublicKey))
	copy(p.Token[:], token[:])
	copy(p.Signature[:], sig[:])
	return p
}

// PingFromBytes deserializes a ping message.
//
// Returns nil if ping cannot be deserialized.
func PingFromBytes(b []byte) *Ping {
	if len(b) != PingSize {
		return nil
	}
	p := new(Ping)
	copy(p.From[:], b[0:32])
	copy(p.Token[:], b[32:64])
	copy(p.Signature[:], b[64:128])
	return p
}

// ToBytes serializes a ping message.
func (p *Ping) ToBytes(b []byte) {
	copy(b[0:32], p.From[:])
	copy(b[32:64], p.Token[:])
	copy(b[64:128], p.Signature[:])
}

// Verify checks the Ping's signature.
func (p *Ping) Verify() bool {
	return ed25519.Verify(p.From[:], p.Token[:], p.Signature[:])
}
