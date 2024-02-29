package lnwire

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// MaxPeerStorageBytes is maximum size in bytes of blob in peer storage message.
const MaxPeerStorageBytes = 65531

var ErrPeerStorageBytesExceeded = fmt.Errorf("peer storage bytes exceede" +
	"d")

// PeerStorageBlob is the type of the data sent by peers to other peers for
// backup.
type PeerStorageBlob []byte

// PeerStorage is the message in which the peer sends to store its data with
// another peer.
type PeerStorage struct {
	// Blob is the data that a peer wants to back up with another peer.
	Blob PeerStorageBlob
}

// NewPeerStorageMsg creates new instance of PeerStorage message object.
func NewPeerStorageMsg(data PeerStorageBlob) (*PeerStorage, error) {
	if len(data) > MaxPeerStorageBytes {
		return nil, ErrPeerStorageBytesExceeded
	}

	return &PeerStorage{
		Blob: data,
	}, nil
}

// A compile time check to ensure Init implements the lnwire.Message
// interface.
var _ Message = (*PeerStorage)(nil)

// Decode deserializes a serialized PeerStorage message stored in the passed
// io.Reader observing the specified protocol version.
//
// This is part of the lnwire.Message interface.
func (msg *PeerStorage) Decode(r io.Reader, _ uint32) error {
	return ReadElements(r,
		&msg.Blob,
	)
}

// Encode serializes the target PeerStorage into the passed io.Writer observing
// the protocol version specified.
//
// This is part of the lnwire.Message interface.
func (msg *PeerStorage) Encode(w *bytes.Buffer, _ uint32) error {
	var l [2]byte
	length := len(msg.Blob)
	binary.BigEndian.PutUint16(l[:], uint16(length))
	if _, err := w.Write(l[:]); err != nil {
		return err
	}

	return WriteBytes(w, msg.Blob)
}

// MsgType returns the integer uniquely identifying this message type on the
// wire.
//
// This is part of the lnwire.Message interface.
func (msg *PeerStorage) MsgType() MessageType {
	return MsgPeerStorage
}
