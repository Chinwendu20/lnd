package sweep

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wtxmgr"
	"github.com/lightningnetwork/lnd/lntest/mock"
	"github.com/lightningnetwork/lnd/lnwallet"
	"github.com/lightningnetwork/lnd/lnwallet/chainfee"
	"github.com/stretchr/testify/require"
)

// TestDetermineFeePerKw tests that given a fee preference, the
// DetermineFeePerKw will properly map it to a concrete fee in sat/kw.
func TestDetermineFeePerKw(t *testing.T) {
	t.Parallel()

	defaultFee := chainfee.SatPerKWeight(999)
	relayFee := chainfee.SatPerKWeight(300)

	feeEstimator := newMockFeeEstimator(defaultFee, relayFee)

	// We'll populate two items in the internal map which is used to query
	// a fee based on a confirmation target: the default conf target, and
	// an arbitrary conf target. We'll ensure below that both of these are
	// properly
	feeEstimator.blocksToFee[50] = 300
	feeEstimator.blocksToFee[defaultNumBlocksEstimate] = 1000

	testCases := []struct {
		// feePref is the target fee preference for this case.
		feePref FeePreference

		// fee is the value the DetermineFeePerKw should return given
		// the FeePreference above
		fee chainfee.SatPerKWeight

		// fail determines if this test case should fail or not.
		fail bool
	}{
		// A fee rate below the floor should error out.
		{
			feePref: FeePreference{
				FeeRate: chainfee.SatPerKWeight(99),
			},
			fail: true,
		},

		// A fee rate below the relay fee should error out.
		{
			feePref: FeePreference{
				FeeRate: chainfee.SatPerKWeight(299),
			},
			fail: true,
		},

		// A fee rate above the floor, should pass through and return
		// the target fee rate.
		{
			feePref: FeePreference{
				FeeRate: 900,
			},
			fee: 900,
		},

		// A specified confirmation target should cause the function to
		// query the estimator which will return our value specified
		// above.
		{
			feePref: FeePreference{
				ConfTarget: 50,
			},
			fee: 300,
		},

		// If the caller doesn't specify any values at all, then we
		// should query for the default conf target.
		{
			feePref: FeePreference{},
			fee:     1000,
		},

		// Both conf target and fee rate are set, we should return with
		// an error.
		{
			feePref: FeePreference{
				ConfTarget: 50,
				FeeRate:    90000,
			},
			fee:  300,
			fail: true,
		},
	}
	for i, testCase := range testCases {
		targetFee, err := DetermineFeePerKw(
			feeEstimator, testCase.feePref,
		)
		switch {
		case testCase.fail && err != nil:
			continue

		case testCase.fail && err == nil:
			t.Fatalf("expected failure for #%v", i)

		case !testCase.fail && err != nil:
			t.Fatalf("unable to estimate fee; %v", err)
		}

		if targetFee != testCase.fee {
			t.Fatalf("#%v: wrong fee: expected %v got %v", i,
				testCase.fee, targetFee)
		}
	}
}

type mockUtxoSource struct {
	outputs []*lnwallet.Utxo
}

func newMockUtxoSource(utxos []*lnwallet.Utxo) *mockUtxoSource {
	return &mockUtxoSource{
		outputs: utxos,
	}
}

func (m *mockUtxoSource) ListUnspentWitnessFromDefaultAccount(minConfs int32,
	maxConfs int32) ([]*lnwallet.Utxo, error) {

	return m.outputs, nil
}

type mockCoinSelectionLocker struct {
	fail bool
}

func (m *mockCoinSelectionLocker) WithCoinSelectLock(f func() error) error {
	if err := f(); err != nil {
		return err
	}

	if m.fail {
		return fmt.Errorf("kek")
	}

	return nil

}

type mockOutputLeaser struct {
	leasedOutputs map[wire.OutPoint]struct{}

	releasedOutputs map[wire.OutPoint]struct{}
}

func newMockOutputLeaser() *mockOutputLeaser {
	return &mockOutputLeaser{
		leasedOutputs: make(map[wire.OutPoint]struct{}),

		releasedOutputs: make(map[wire.OutPoint]struct{}),
	}
}

func (m *mockOutputLeaser) LeaseOutput(_ wtxmgr.LockID, o wire.OutPoint,
	t time.Duration) (time.Time, []byte, btcutil.Amount, error) {

	m.leasedOutputs[o] = struct{}{}

	return time.Now().Add(t), nil, 0, nil
}

func (m *mockOutputLeaser) ReleaseOutput(_ wtxmgr.LockID,
	o wire.OutPoint) error {

	m.releasedOutputs[o] = struct{}{}

	return nil
}

func (m *mockOutputLeaser) Reset() {
	m.releasedOutputs = make(map[wire.OutPoint]struct{})

	m.leasedOutputs = make(map[wire.OutPoint]struct{})
}

var sweepScript = []byte{
	0x0, 0x14, 0x64, 0x3d, 0x8b, 0x15, 0x69, 0x4a, 0x54,
	0x7d, 0x57, 0x33, 0x6e, 0x51, 0xdf, 0xfd, 0x38, 0xe3,
	0xe, 0x6e, 0xf8, 0xef,
}

var deliveryAddr = func() btcutil.Address {
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(
		sweepScript, &chaincfg.TestNet3Params,
	)
	if err != nil {
		panic(err)
	}

	return addrs[0]
}()

var testUtxos = []*lnwallet.Utxo{
	{
		// A p2wkh output.
		AddressType: lnwallet.WitnessPubKey,
		PkScript: []byte{
			0x0, 0x14, 0x64, 0x3d, 0x8b, 0x15, 0x69, 0x4a, 0x54,
			0x7d, 0x57, 0x33, 0x6e, 0x51, 0xdf, 0xfd, 0x38, 0xe3,
			0xe, 0x6e, 0xf7, 0xef,
		},
		Value: 1000,
		OutPoint: wire.OutPoint{
			Index: 1,
		},
	},

	{
		// A np2wkh output.
		AddressType: lnwallet.NestedWitnessPubKey,
		PkScript: []byte{
			0xa9, 0x14, 0x97, 0x17, 0xf7, 0xd1, 0x5f, 0x6f, 0x8b,
			0x7, 0xe3, 0x58, 0x43, 0x19, 0xb9, 0x7e, 0xa9, 0x20,
			0x18, 0xc3, 0x17, 0xd7, 0x87,
		},
		Value: 2000,
		OutPoint: wire.OutPoint{
			Index: 2,
		},
	},

	// A p2wsh output.
	{
		AddressType: lnwallet.UnknownAddressType,
		PkScript: []byte{
			0x0, 0x20, 0x70, 0x1a, 0x8d, 0x40, 0x1c, 0x84, 0xfb, 0x13,
			0xe6, 0xba, 0xf1, 0x69, 0xd5, 0x96, 0x84, 0xe2, 0x7a, 0xbd,
			0x9f, 0xa2, 0x16, 0xc8, 0xbc, 0x5b, 0x9f, 0xc6, 0x3d, 0x62,
			0x2f, 0xf8, 0xc5, 0x8c,
		},
		Value: 3000,
		OutPoint: wire.OutPoint{
			Index: 3,
		},
	},
}

func assertUtxosLeased(t *testing.T, utxoLeaser *mockOutputLeaser,
	utxos []*lnwallet.Utxo) {

	t.Helper()

	for _, utxo := range utxos {
		if _, ok := utxoLeaser.leasedOutputs[utxo.OutPoint]; !ok {
			t.Fatalf("utxo %v was never leased", utxo.OutPoint)
		}
	}

}

func assertNoUtxosReleased(t *testing.T, utxoLeaser *mockOutputLeaser,
	utxos []*lnwallet.Utxo) {

	t.Helper()

	if len(utxoLeaser.releasedOutputs) != 0 {
		t.Fatalf("outputs have been released, but shouldn't have been")
	}
}

func assertUtxosReleased(t *testing.T, utxoLeaser *mockOutputLeaser,
	utxos []*lnwallet.Utxo) {

	t.Helper()

	for _, utxo := range utxos {
		if _, ok := utxoLeaser.releasedOutputs[utxo.OutPoint]; !ok {
			t.Fatalf("utxo %v was never released", utxo.OutPoint)
		}
	}
}

func assertUtxosLeasedAndReleased(t *testing.T, utxoLeaser *mockOutputLeaser,
	utxos []*lnwallet.Utxo) {

	t.Helper()

	for _, utxo := range utxos {
		if _, ok := utxoLeaser.leasedOutputs[utxo.OutPoint]; !ok {
			t.Fatalf("utxo %v was never leased", utxo.OutPoint)
		}

		if _, ok := utxoLeaser.releasedOutputs[utxo.OutPoint]; !ok {
			t.Fatalf("utxo %v was never released", utxo.OutPoint)
		}
	}
}

// TestCraftSweepAllTxCoinSelectFail tests that if coin selection fails, then
// we unlock any outputs we may have locked in the passed closure.
func TestCraftSweepAllTxCoinSelectFail(t *testing.T) {
	t.Parallel()

	utxoSource := newMockUtxoSource(testUtxos)
	coinSelectLocker := &mockCoinSelectionLocker{
		fail: true,
	}
	utxoLeaser := newMockOutputLeaser()

	_, err := CraftSweepAllTx(
		0, 0, 10, nil, nil, coinSelectLocker, utxoSource, utxoLeaser,
		nil, 0, nil,
	)

	// Since we instructed the coin select locker to fail above, we should
	// get an error.
	if err == nil {
		t.Fatalf("sweep tx should have failed: %v", err)
	}

	// At this point, we'll now verify that all outputs were initially
	// locked, and then also unlocked due to the failure.
	assertUtxosLeasedAndReleased(t, utxoLeaser, testUtxos)
}

// TestCraftSweepAllTxUnknownWitnessType tests that if one of the inputs we
// encounter is of an unknown witness type, then we fail and unlock any prior
// locked outputs.
func TestCraftSweepAllTxUnknownWitnessType(t *testing.T) {
	t.Parallel()

	utxoSource := newMockUtxoSource(testUtxos)
	coinSelectLocker := &mockCoinSelectionLocker{}
	utxoLeaser := newMockOutputLeaser()

	_, err := CraftSweepAllTx(
		0, 0, 10, nil, nil, coinSelectLocker, utxoSource, utxoLeaser,
		nil, 0, nil,
	)

	// Since passed in a p2wsh output, which is unknown, we should fail to
	// map the output to a witness type.
	if err == nil {
		t.Fatalf("sweep tx should have failed: %v", err)
	}

	// At this point, we'll now verify that all outputs were initially
	// locked, and then also unlocked since we weren't able to find a
	// witness type for the last output.
	assertUtxosLeasedAndReleased(t, utxoLeaser, testUtxos)
}

// TestCraftSweepAllTx tests that we'll properly lock all available outputs
// within the wallet, and craft a single sweep transaction that pays to the
// target output.
func TestCraftSweepAllTx(t *testing.T) {
	t.Parallel()

	// First, we'll make a mock signer along with a fee estimator, We'll
	// use zero fees to we can assert a precise output value.
	signer := &mock.DummySigner{}

	// For our UTXO source, we'll pass in all the UTXOs that we know of,
	// other than the final one which is of an unknown witness type.
	targetUTXOs := testUtxos[:2]
	utxoSource := newMockUtxoSource(targetUTXOs)
	coinSelectLocker := &mockCoinSelectionLocker{}
	utxoLeaser := newMockOutputLeaser()

	testCases := []struct {
		name        string
		selectUtxos []wire.OutPoint
	}{
		{
			name:        "sweep all utxos in wallet",
			selectUtxos: nil,
		},
		{
			name: "sweep select utxos in wallet",
			selectUtxos: []wire.OutPoint{
				targetUTXOs[0].OutPoint,
			},
		},
	}

	for _, tc := range testCases {
		sweepPkg, err := CraftSweepAllTx(
			0, 0, 10, nil, deliveryAddr, coinSelectLocker,
			utxoSource, utxoLeaser, signer, 0, tc.selectUtxos,
		)
		require.NoError(t, err, "unable to make sweep tx")

		inputUtxos := targetUTXOs
		if len(tc.selectUtxos) != 0 {
			inputUtxos, err = fetchUtxosFromOutpoints(
				targetUTXOs, tc.selectUtxos,
			)
			require.NoError(t, err)
		}

		expectedSweepValue := SumUtxos(inputUtxos)

		// At this point, all the UTXOs that we made above should be
		// locked and none of them unlocked.
		assertUtxosLeased(t, utxoLeaser, inputUtxos)
		assertNoUtxosReleased(t, utxoLeaser, inputUtxos)

		// Now that we have our sweep transaction, we should find that
		// we have a UTXO for each input, and also that our final output
		// value is the sum of all our inputs.
		sweepTx := sweepPkg.SweepTx
		require.Equal(t, len(sweepTx.TxIn), len(inputUtxos))

		lookup := make(map[wire.OutPoint]struct{}, len(sweepTx.TxIn))
		for _, tx := range sweepTx.TxIn {
			lookup[tx.PreviousOutPoint] = struct{}{}
		}

		for _, utxo := range inputUtxos {
			_, ok := lookup[utxo.OutPoint]
			require.True(t, ok)
		}

		// We should have a single output that pays to our sweep script
		// generated above.
		if len(sweepTx.TxOut) != 1 {
			t.Fatalf("should have %v outputs, instead have %v", 1,
				len(sweepTx.TxOut))
		}
		output := sweepTx.TxOut[0]
		switch {
		case output.Value != expectedSweepValue:
			t.Fatalf("expected %v sweep value, instead got %v",
				expectedSweepValue, output.Value)

		case !bytes.Equal(sweepScript, output.PkScript):
			t.Fatalf("expected %x sweep script, instead "+
				"got %x", sweepScript, output.PkScript)
		}

		// If we cancel the sweep attempt,
		// then we should find that all the
		// UTXOs within the sweep transaction are now unlocked.
		sweepPkg.CancelSweepAttempt()
		assertUtxosReleased(t, utxoLeaser, inputUtxos)

		// Reset all utxoleaser values after test.
		utxoLeaser.Reset()
	}
}

func SumUtxos(utxos []*lnwallet.Utxo) int64 {
	sum := int64(0)
	for _, utxo := range utxos {
		sum += int64(utxo.Value)
	}

	return sum
}
