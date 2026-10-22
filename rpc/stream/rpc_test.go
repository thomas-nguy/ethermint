package stream

import (
	"fmt"
	"math/big"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	proto "google.golang.org/protobuf/proto"
)

func TestBlockTimeTracker(t *testing.T) {
	newLogs := func(n int) []*ethtypes.Log {
		logs := make([]*ethtypes.Log, n)
		for i := range logs {
			logs[i] = &ethtypes.Log{}
		}
		return logs
	}
	requireTimestamps := func(t *testing.T, logs []*ethtypes.Log, blockTime uint64) {
		t.Helper()
		for _, l := range logs {
			require.Equal(t, blockTime, l.BlockTimestamp)
		}
	}

	t.Run("logs after block are timestamped immediately", func(t *testing.T) {
		tracker := newBlockTimeTracker()
		require.Empty(t, tracker.addBlock(10, 100))
		ready := tracker.addLogs(10, newLogs(2))
		require.Len(t, ready, 2)
		requireTimestamps(t, ready, 100)
	})

	t.Run("logs before block are held until the block arrives", func(t *testing.T) {
		tracker := newBlockTimeTracker()
		require.Empty(t, tracker.addLogs(10, newLogs(2)))
		require.Empty(t, tracker.addLogs(10, newLogs(1)))
		ready := tracker.addBlock(10, 100)
		require.Len(t, ready, 3)
		requireTimestamps(t, ready, 100)
		require.Empty(t, tracker.pending)
	})

	t.Run("logs of a missed block are released by a later block", func(t *testing.T) {
		tracker := newBlockTimeTracker()
		require.Empty(t, tracker.addLogs(10, newLogs(1)))
		require.Empty(t, tracker.addLogs(12, newLogs(1)))
		ready := tracker.addBlock(11, 110)
		require.Len(t, ready, 1)
		requireTimestamps(t, ready, 0)
		ready = tracker.addBlock(12, 120)
		require.Len(t, ready, 1)
		requireTimestamps(t, ready, 120)
	})

	t.Run("pending logs are released in block order", func(t *testing.T) {
		tracker := newBlockTimeTracker()
		for h := uint64(15); h >= 10; h-- {
			tracker.addLogs(int64(h), []*ethtypes.Log{{BlockNumber: h}})
		}
		ready := tracker.addBlock(13, 130)
		require.Len(t, ready, 4)
		for i, l := range ready {
			require.Equal(t, uint64(10+i), l.BlockNumber)
		}
		requireTimestamps(t, ready[3:], 130)

		ready = tracker.flush()
		require.Len(t, ready, 2)
		require.Equal(t, uint64(14), ready[0].BlockNumber)
		require.Equal(t, uint64(15), ready[1].BlockNumber)
	})

	t.Run("logs of an unknown past block are not held back", func(t *testing.T) {
		tracker := newBlockTimeTracker()
		for h := int64(1); h <= blockTimeRetention+10; h++ {
			tracker.addBlock(h, uint64(h))
		}
		// pruned block time
		ready := tracker.addLogs(1, newLogs(1))
		require.Len(t, ready, 1)
		requireTimestamps(t, ready, 0)
		// retained block time
		ready = tracker.addLogs(blockTimeRetention, newLogs(1))
		require.Len(t, ready, 1)
		requireTimestamps(t, ready, blockTimeRetention)
		require.Empty(t, tracker.pending)
	})

	t.Run("old block times are pruned", func(t *testing.T) {
		tracker := newBlockTimeTracker()
		for h := int64(1); h <= 2*blockTimeRetention; h++ {
			tracker.addBlock(h, uint64(h))
		}
		require.Len(t, tracker.times, blockTimeRetention)
	})

	t.Run("flush releases all pending logs", func(t *testing.T) {
		tracker := newBlockTimeTracker()
		tracker.addLogs(10, newLogs(1))
		tracker.addLogs(11, newLogs(2))
		require.Len(t, tracker.flush(), 3)
		require.Empty(t, tracker.pending)
	})

	t.Run("logs after flush are not held back and use known block times", func(t *testing.T) {
		tracker := newBlockTimeTracker()
		tracker.addBlock(10, 100)
		tracker.flush()
		ready := tracker.addLogs(10, newLogs(1))
		require.Len(t, ready, 1)
		requireTimestamps(t, ready, 100)
		ready = tracker.addLogs(11, newLogs(1))
		require.Len(t, ready, 1)
		requireTimestamps(t, ready, 0)
		require.Empty(t, tracker.pending)
	})
}

func TestEvmTxHashFromEventData(t *testing.T) {
	t.Run("empty block returns EmptyRootHash", func(t *testing.T) {
		data := tmtypes.EventDataNewBlock{
			Block:               &tmtypes.Block{},
			ResultFinalizeBlock: abci.ResponseFinalizeBlock{},
		}
		txDecoder := func([]byte) (sdk.Tx, error) { return nil, nil }
		got, err := evmTxHashFromEventData(data, txDecoder)
		require.NoError(t, err)
		require.Equal(t, ethtypes.EmptyRootHash, got)
	})

	t.Run("all txs failed returns EmptyRootHash", func(t *testing.T) {
		data := tmtypes.EventDataNewBlock{
			Block: &tmtypes.Block{
				Data: tmtypes.Data{Txs: tmtypes.Txs{[]byte("tx")}},
			},
			ResultFinalizeBlock: abci.ResponseFinalizeBlock{
				TxResults: []*abci.ExecTxResult{{Code: 1}},
			},
		}
		txDecoder := func([]byte) (sdk.Tx, error) { return nil, nil }
		got, err := evmTxHashFromEventData(data, txDecoder)
		require.NoError(t, err)
		require.Equal(t, ethtypes.EmptyRootHash, got)
	})

	t.Run("txDecoder error returns error", func(t *testing.T) {
		data := tmtypes.EventDataNewBlock{
			Block: &tmtypes.Block{
				Data: tmtypes.Data{Txs: tmtypes.Txs{[]byte("invalid")}},
			},
			ResultFinalizeBlock: abci.ResponseFinalizeBlock{
				TxResults: []*abci.ExecTxResult{{Code: 0}},
			},
		}
		txDecoder := func([]byte) (sdk.Tx, error) {
			return nil, fmt.Errorf("cannot decode tx")
		}
		_, err := evmTxHashFromEventData(data, txDecoder)
		require.Error(t, err)
	})

	t.Run("tx/result count mismatch returns error", func(t *testing.T) {
		data := tmtypes.EventDataNewBlock{
			Block: &tmtypes.Block{
				Data: tmtypes.Data{Txs: tmtypes.Txs{[]byte("tx1"), []byte("tx2")}},
			},
			ResultFinalizeBlock: abci.ResponseFinalizeBlock{
				TxResults: []*abci.ExecTxResult{{Code: 0}}, // only 1 result for 2 txs
			},
		}
		txDecoder := func([]byte) (sdk.Tx, error) { return nil, nil }
		_, err := evmTxHashFromEventData(data, txDecoder)
		require.Error(t, err)
	})

	t.Run("no EVM messages returns EmptyRootHash", func(t *testing.T) {
		data := tmtypes.EventDataNewBlock{
			Block: &tmtypes.Block{
				Data: tmtypes.Data{Txs: tmtypes.Txs{[]byte("cosmos-tx")}},
			},
			ResultFinalizeBlock: abci.ResponseFinalizeBlock{
				TxResults: []*abci.ExecTxResult{{Code: 0}},
			},
		}
		txDecoder := func([]byte) (sdk.Tx, error) {
			return &mockCosmosOnlyTx{}, nil
		}
		got, err := evmTxHashFromEventData(data, txDecoder)
		require.NoError(t, err)
		require.Equal(t, ethtypes.EmptyRootHash, got)
	})

	t.Run("EVM tx produces non-empty trie root matching EvmTxHashFromMsgs", func(t *testing.T) {
		to := common.HexToAddress("0x1234567890123456789012345678901234567890")
		msg := evmtypes.NewTx(big.NewInt(9000), 0, &to, big.NewInt(0), 21000, big.NewInt(1), nil, nil, nil, nil)
		data := tmtypes.EventDataNewBlock{
			Block: &tmtypes.Block{
				Data: tmtypes.Data{Txs: tmtypes.Txs{[]byte("evm-tx")}},
			},
			ResultFinalizeBlock: abci.ResponseFinalizeBlock{
				TxResults: []*abci.ExecTxResult{{Code: 0}},
			},
		}
		txDecoder := func([]byte) (sdk.Tx, error) {
			return &mockEvmTx{msgs: []*evmtypes.MsgEthereumTx{msg}}, nil
		}
		got, err := evmTxHashFromEventData(data, txDecoder)
		require.NoError(t, err)
		require.NotEqual(t, ethtypes.EmptyRootHash, got)
		require.Equal(t, rpctypes.EvmTxHashFromMsgs([]*evmtypes.MsgEthereumTx{msg}), got)
	})
}

// mockCosmosOnlyTx is an sdk.Tx with no EVM messages.
type mockCosmosOnlyTx struct{}

func (m *mockCosmosOnlyTx) GetMsgs() []sdk.Msg                  { return []sdk.Msg{} }
func (m *mockCosmosOnlyTx) ValidateBasic() error                { return nil }
func (m *mockCosmosOnlyTx) GetMsgsV2() ([]proto.Message, error) { return nil, nil }

// mockEvmTx is an sdk.Tx wrapping one or more EVM messages.
type mockEvmTx struct{ msgs []*evmtypes.MsgEthereumTx }

func (m *mockEvmTx) GetMsgs() []sdk.Msg {
	out := make([]sdk.Msg, len(m.msgs))
	for i, msg := range m.msgs {
		out[i] = msg
	}
	return out
}
func (m *mockEvmTx) ValidateBasic() error                { return nil }
func (m *mockEvmTx) GetMsgsV2() ([]proto.Message, error) { return nil, nil }
