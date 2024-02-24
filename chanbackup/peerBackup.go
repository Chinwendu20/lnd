package chanbackup

import (
	"bytes"
	"github.com/lightningnetwork/lnd/keychain"
)

type PeerBackup struct {
	BackupState map[string][]byte
	keychain    keychain.KeyRing
	Swapper
}

type PeerBackupConfig struct {
	Keychain keychain.KeyRing
	Swapper  Swapper
}

func NewPeerBackup(c *PeerBackupConfig) *PeerBackup {
	p := PeerBackup{}
	return &p
}

func (p *PeerBackup) VerifyBackUp(data []byte, latestPeerBackup *Multi) (error,
	bool, *Multi) {

	newReader := bytes.NewReader(data)
	var decodedPeerBackup Multi

	err := decodedPeerBackup.UnpackFromReader(newReader,
		p.keychain)

	if err != nil {
		return err, false, nil
	}

	return nil, decodedPeerBackup.DeepEqual(
		latestPeerBackup), &decodedPeerBackup

}

func (p *PeerBackup) StorePeerBackup(data []byte, peerAddr string) ([]byte,
	error) {
	storedBackup, err := p.UpdateAndSwapPeerBackup(peerAddr, data,
		p.keychain)
	if err != nil {
		return nil, err
	}
	p.BackupState[peerAddr] = storedBackup

	return storedBackup, nil
}

func (p *PeerBackup) RetrieveBackupForPeer(peerAddr string) []byte {
	return p.BackupState[peerAddr]
}
