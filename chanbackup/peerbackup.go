package chanbackup

import (
	"bytes"
	"github.com/lightningnetwork/lnd/keychain"
	"sync"
)

type PeerBackup struct {
	BkupStateMtx sync.RWMutex
	BackupState  map[string][]byte
	keychain     keychain.KeyRing
	Swapper
}

type PeerBackupConfig struct {
	Keychain keychain.KeyRing
	Swapper  Swapper
}

func NewPeerBackup(c *PeerBackupConfig) (*PeerBackup, error) {
	peerBackup, err := c.Swapper.ExtractPeerBackupMulti(c.Keychain)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(peerBackup)
	bs, err := UnPackPartialPeerBackupMulti(reader)
	if err != nil {
		return nil, err
	}

	return &PeerBackup{
		BackupState: bs,
		keychain:    c.Keychain,
		Swapper:     c.Swapper,
	}, nil
}

func (p *PeerBackup) VerifyBackup(data []byte, latestPeerBackup []byte) (error,
	bool) {

	// read data to Multi structure
	newReader := bytes.NewReader(data)
	var decodedData Multi

	err := decodedData.UnpackFromReader(newReader,
		p.keychain)

	if err != nil {
		return err, false
	}

	// read latestPeerBackup to Multi structure
	newReader = bytes.NewReader(latestPeerBackup)
	var decodedLatestPeerBackup Multi

	err = decodedLatestPeerBackup.UnpackFromReader(newReader,
		p.keychain)

	if err != nil {
		return err, false
	}

	return nil, decodedData.DeepEqual(
		&decodedLatestPeerBackup)

}

func (p *PeerBackup) StorePeerBackup(data []byte, peerAddr string) ([]byte,
	error) {

	p.BkupStateMtx.RLock()
	p.BackupState[peerAddr] = data
	p.BkupStateMtx.RUnlock()

	newBackup := p.BackupState
	err := p.UpdateAndSwapPeerBackup(newBackup)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *PeerBackup) RetrieveBackupForPeer(peerAddr string) []byte {
	p.BkupStateMtx.RLock()
	defer p.BkupStateMtx.RUnlock()
	return p.BackupState[peerAddr]
}

func (p *PeerBackup) PrunePeerStorage(peerAddrs []string) error {

	p.BkupStateMtx.Lock()
	for _, addr := range peerAddrs {
		delete(p.BackupState, addr)
	}
	p.BkupStateMtx.Unlock()

	newBackup := p.BackupState
	err := p.UpdateAndSwapPeerBackup(newBackup)
	if err != nil {
		return err
	}

	return nil
}
