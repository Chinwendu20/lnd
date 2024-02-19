package lnwire

import (
	"bytes"
	"encoding/binary"
	"io"
)

// YourPeerStorageBlob is the type of the data stored by peers as backup for
// other peers. This message is sent in response to the PeerStorage message.
type YourPeerStorageBlob []byte

// YourPeerStorage message is message sent by a peer in response to a
// PeerStorage message. It contains the data stored as backup for the peer
// sending the per storage message.
type YourPeerStorage struct {
	// Blob is the data backed up by a peer for another peer.
	Blob YourPeerStorageBlob
}

// NewYourPeerStorageMsg creates new instance of YourPeerStorage message object.
func NewYourPeerStorageMsg(data YourPeerStorageBlob) *YourPeerStorage {
	return &YourPeerStorage{
		Blob: data,
	}
}

// A compile time check to ensure Init implements the lnwire.Message
// interface.
var _ Message = (*YourPeerStorage)(nil)

// Decode deserializes a serialized YourPeerStorage message stored in the passed
// io.Reader observing the specified protocol version.
//
// This is part of the lnwire.Message interface.
func (msg *YourPeerStorage) Decode(r io.Reader, _ uint32) error {
	return ReadElements(r,
		&msg.Blob,
	)
}

// Encode serializes the target YourPeerStorage into the passed io.Writer
// observing the protocol version specified.
//
// This is part of the lnwire.Message interface.
func (msg *YourPeerStorage) Encode(w *bytes.Buffer, _ uint32) error {
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
func (msg *YourPeerStorage) MsgType() MessageType {
	return MsgYourPeerStorage
}
