package lnd

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/lightningnetwork/lnd/chainntnfs"
	"github.com/lightningnetwork/lnd/kvdb"
	"github.com/prometheus/common/log"
	"sync"
)

const (
	maxWipeWindow = 2016
)

var (
	peerStorage = []byte("peer-storage")
)

type blockViewer interface {
	BestHeight() int32
}

type kvdbPeerStorage struct {
	*chainntnfs.BlockEpochEvent
	db             kvdb.Backend
	awaitingDelete map[int32]*btcec.PublicKey
	mapMtx         sync.RWMutex
	quit           chan interface{}
	blockViewer    chainntnfs.BestBlockView
	wg             sync.WaitGroup
}

type peerStorageConfig struct {
	chainNotifier chainntnfs.ChainNotifier
	db            kvdb.Backend
	blockViewer   chainntnfs.BestBlockView
}

func newKvdbPeerStorage(c *peerStorageConfig) (*kvdbPeerStorage, error) {
	epochEvent, err := c.chainNotifier.RegisterBlockEpochNtfn(nil)

	if err != nil {
		return nil, err
	}

	return &kvdbPeerStorage{
		db:              c.db,
		BlockEpochEvent: epochEvent,
		blockViewer:     c.blockViewer,
	}, nil
}

func (k *kvdbPeerStorage) start() {
	errChan := make(chan error)

	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		err := k.garbageCollector()
		select {
		case errChan <- err:
		case <-k.quit:
		}

	}()

	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		select {
		case err := <-errChan:
			if err != nil {
				log.Errorf("kvdbPeerStorage experienced "+
					"error during operation %v", err)
			}
		case <-k.quit:
		}
	}()
}

func (k *kvdbPeerStorage) stop() {
	close(k.quit)
	k.wg.Wait()

	return
}

func (k *kvdbPeerStorage) StorePeerData(data []byte,
	peerPub *btcec.PublicKey) error {
	// Check if to delete

	return kvdb.Update(k.db, func(tx kvdb.RwTx) error {
		bucket, err := tx.CreateTopLevelBucket(peerStorage)
		if err != nil {
			return err
		}

		peerPubBytes := peerPub.SerializeCompressed()

		return bucket.Put(peerPubBytes, data)
	}, func() {})
}

func (k *kvdbPeerStorage) RetrievePeerData(peerPub *btcec.PublicKey) ([]byte,
	error) {

	var data []byte
	err := kvdb.View(k.db, func(tx kvdb.RTx) error {
		bucket := tx.ReadBucket(peerStorage)

		peerPubBytes := peerPub.SerializeCompressed()

		data = bucket.Get(peerPubBytes)
		return nil
	}, func() {})

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (k *kvdbPeerStorage) MarkForDelete(peerPub *btcec.PublicKey) error {
	k.mapMtx.Lock()
	defer k.mapMtx.Unlock()

	height, err := k.blockViewer.BestHeight()
	if err != nil {
		return err
	}

	k.awaitingDelete[int32(height+maxWipeWindow)] = peerPub

	return nil
}

func (k *kvdbPeerStorage) garbageCollector() error {

	for {
		select {
		case e := <-k.Epochs:
			k.mapMtx.RLock()
			peerPub, ok := k.awaitingDelete[e.Height]
			k.mapMtx.RUnlock()

			if ok {
				err := kvdb.Update(k.db, func(tx kvdb.RwTx) error {
					bucket, err := tx.CreateTopLevelBucket(
						peerStorage)
					if err != nil {
						return err
					}

					peerPubBytes := peerPub.
						SerializeCompressed()

					return bucket.Delete(peerPubBytes)
				}, func() {})

				if err != nil {
					return err
				}

			}
		case <-k.quit:
			k.Cancel()
			return nil
		}
	}
}
