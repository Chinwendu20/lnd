package chanbackup

import (
	"bytes"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/lightningnetwork/lnd/lnwire"
	"io"
)

// PeerBackupMulti is the structure used to store backups for peer.
type PeerBackupMulti map[string]*Multi

// Pack encrypts and encodes PeerBackupMulti to bytes.
func (p PeerBackupMulti) Pack(b *bytes.Buffer, keyRing keychain.KeyRing,
) error {
	err := lnwire.WriteElements(b, uint32(len(p)))
	if err != nil {
		return err
	}

	for k, v := range p {
		err = lnwire.WriteElements(b, k)
		if err != nil {
			return err
		}

		err = v.PackToWriter(b, keyRing)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p PeerBackupMulti) PartiallyPack(keyRing keychain.KeyRing) (
	map[string][]byte, error) {
	var data map[string][]byte
	for k, v := range p {
		var b bytes.Buffer
		err := v.PackToWriter(&b, keyRing)
		if err != nil {
			return nil, err
		}
		data[k] = b.Bytes()
	}
	return data, nil
}

// UnpackPeerBackupMulti decrypts and decodes a reader to `PeerBackupMulti`.
func (p PeerBackupMulti) UnpackPeerBackupMulti(r io.Reader,
	keyRing keychain.KeyRing) error {

	var numPeerBackupMulti uint32
	err := lnwire.ReadElements(r, &numPeerBackupMulti)
	if err != nil {
		return err
	}

	for i := uint32(0); i < numPeerBackupMulti; i++ {
		var key string
		err := lnwire.ReadElements(r, &key)
		if err != nil {
			return err
		}
		var multiBytes []byte
		pMulti := PackedMulti(multiBytes)

		multi, err := pMulti.Unpack(keyRing)
		if err != nil {
			return err
		}

		p[key] = multi
	}

	return nil
}

func PackPartialPeerBackupMulti(data map[string][]byte, b *bytes.Buffer) error {
	err := lnwire.WriteElements(b, uint32(len(data)))
	if err != nil {
		return err
	}

	for k, v := range data {
		err = lnwire.WriteElements(b, k)
		if err != nil {
			return err
		}

		err = lnwire.WriteElements(b, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnPackPartialPeerBackupMulti(r io.Reader) (map[string][]byte, error) {

	var (
		numPeerBackupMulti uint32
		p                  map[string][]byte
	)
	err := lnwire.ReadElements(r, &numPeerBackupMulti)
	if err != nil {
		return nil, err
	}

	for i := uint32(0); i < numPeerBackupMulti; i++ {
		var key string
		err = lnwire.ReadElements(r, &key)
		if err != nil {
			return nil, err
		}
		var b []byte
		err = lnwire.ReadElements(r, &b)
		if err != nil {
			return nil, err
		}

		p[key] = b
	}

	return p, nil
}
