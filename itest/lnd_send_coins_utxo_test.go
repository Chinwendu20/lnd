package itest

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/lightningnetwork/lnd/lntest"
	"github.com/lightningnetwork/lnd/lnwallet"
	"github.com/stretchr/testify/require"
)

type sendCoinTestCase struct {
	// name is the name of the target test case.
	name string

	// amt specifies the sum to be sent.
	amt int64

	// sendCoinShouldFail denotes if we expect the channel opening to fail.
	sendCoinShouldFail bool

	// sendAll specifies all funds in the wallet or select utxos should
	// be sent.
	sendAll bool

	// selectedCoins are coins selected to be used as inputs for the
	// transaction.
	selectedCoins []btcutil.Amount
}

// testSendCoinSelectUtxo tests sendCoins with selected utxos as input.
func testSendCoinSelectUtxo(ht *lntest.HarnessTest) {
	alice := ht.NewNode("Alice", nil)
	defer ht.Shutdown(alice)

	initialCoins := []btcutil.Amount{
		100_000,
		50_000,
		300_000,
		20_000,
		1_000,
	}

	var tcs = []*sendCoinTestCase{
		{
			name: "sendCoins with selected utxos",
			selectedCoins: []btcutil.Amount{
				50_000,
				300_000,
			},
			amt: 220_000,
		},
		{
			name: "sendCoins with selected utxos, sendAll",
			selectedCoins: []btcutil.Amount{
				1_000,
				300_000,
			},
			sendAll: true,
		},
		{
			name:    "sendAll coins in wallet",
			sendAll: true,
		},
		{
			name: "send coins, amount specified, " +
				"no select outpoints",
			amt: 200_000,
		},
		{
			name:               "sendAll specified with amount",
			sendAll:            true,
			amt:                250_000,
			sendCoinShouldFail: true,
		},
		{
			name: "send coins amount specfied greater than " +
				"available amount",
			amt:                1_000_000,
			sendCoinShouldFail: true,
		},
		{
			name: "sendCoins with selected utxos, amt specified " +
				"greater that select utxos funds",
			selectedCoins: []btcutil.Amount{
				50_000,
				1_000,
			},
			amt:                200_000,
			sendCoinShouldFail: true,
		},
	}

	for _, tc := range tcs {
		ht.Run(tc.name, func(t *testing.T) {
			// Fund alice and compute the expected wallet balance.
			totalinitialCoinSum := 0
			for _, initialCoin := range initialCoins {
				ht.FundCoins(initialCoin, alice)
				totalinitialCoinSum += int(initialCoin)
			}

			ht.AssertWalletAccountBalance(
				alice, lnwallet.DefaultAccountName,
				int64(totalinitialCoinSum), 0,
			)

			defer func() {
				// Fund additional coins to sweep in case the
				// wallet contains dust.
				ht.FundCoins(100_000, alice)

				// Remove all funds from Alice.
				sweepNodeWalletAndAssert(ht, alice)
			}()

			// Create an outpoint lookup for each unique amount.
			lookup := make(map[int64]*lnrpc.OutPoint)
			res := alice.RPC.ListUnspent(
				&walletrpc.ListUnspentRequest{},
			)
			for _, utxo := range res.Utxos {
				lookup[utxo.AmountSat] = utxo.Outpoint
			}

			// Map the selected coin to the respective outpoint.
			var selectedOutpoints []*lnrpc.OutPoint
			selectedOpLookup := make(map[string]struct{})
			for _, selectedCoin := range tc.selectedCoins {
				if o, ok := lookup[int64(selectedCoin)]; ok {
					selectedOpLookup[o.TxidStr] = struct{}{}
					selectedOutpoints = append(
						selectedOutpoints, o,
					)
				}
			}

			// Generate address for the request.
			p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(
				make([]byte, 20), harnessNetParams,
			)

			require.NoError(t, err, "error generating "+
				"target address for request")

			sendCoinReq := &lnrpc.SendCoinsRequest{
				Addr:      p2wkhAddr.String(),
				SendAll:   tc.sendAll,
				Outpoints: selectedOutpoints,
				Amount:    tc.amt,
			}

			// If we don't expect the send coin to work, we should
			// simply check for an error.
			if tc.sendCoinShouldFail {
				alice.RPC.SendCoinsAssertErr(sendCoinReq)
				return
			}

			alice.RPC.SendCoins(sendCoinReq)
			block := ht.MineBlocksAndAssertNumTxes(
				1, 1,
			)[0]
			tx := block.Transactions[1]

			if len(selectedOutpoints) > 0 {
				require.Equal(
					t, len(tx.TxIn), len(selectedOutpoints),
				)

				for _, in := range tx.TxIn {
					Op := in.PreviousOutPoint.Hash.String()
					_, ok := selectedOpLookup[Op]
					require.True(t, ok)
				}
			}

			if sendCoinReq.SendAll && tc.selectedCoins == nil {
				ht.AssertWalletAccountBalance(
					alice, lnwallet.DefaultAccountName,
					0, 0,
				)
			} else {
				ht.AssertConfWalletBalanceLT(
					alice, lnwallet.DefaultAccountName,
					int64(totalinitialCoinSum)-tc.amt,
				)
			}

			// When re-selecting a spent output for funding we
			// expect the respective an error.
			if len(selectedOutpoints) > 0 {
				alice.RPC.SendCoinsAssertErr(sendCoinReq)
			}
		})
	}
}
