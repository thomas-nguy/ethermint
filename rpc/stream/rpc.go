package stream

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"cosmossdk.io/log/v2"
	cmtquery "github.com/cometbft/cometbft/libs/pubsub/query"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/rpc/types"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"google.golang.org/grpc"
)

const (
	streamSubscriberName = "ethermint-json-rpc"
	subscribBufferSize   = 1024

	headerStreamSegmentSize = 128
	headerStreamCapacity    = 128 * 32
	txStreamSegmentSize     = 1024
	txStreamCapacity        = 1024 * 32
	logStreamSegmentSize    = 2048
	logStreamCapacity       = 2048 * 32
)

var (
	evmEvents = cmtquery.MustCompile(fmt.Sprintf("%s='%s' AND %s.%s='%s'",
		tmtypes.EventTypeKey,
		tmtypes.EventTx,
		sdk.EventTypeMessage,
		sdk.AttributeKeyModule, evmtypes.ModuleName)).String()
	blockEvents  = tmtypes.QueryForEvent(tmtypes.EventNewBlock).String()
	evmTxHashKey = fmt.Sprintf("%s.%s", evmtypes.TypeMsgEthereumTx, evmtypes.AttributeKeyEthereumTxHash)
)

type RPCHeader struct {
	EthHeader *ethtypes.Header
	Hash      common.Hash
}

type validatorAccountFunc func(
	ctx context.Context, in *evmtypes.QueryValidatorAccountRequest, opts ...grpc.CallOption,
) (*evmtypes.QueryValidatorAccountResponse, error)

// RPCStream provides data streams for newHeads, logs, and pendingTransactions.
// it's only started on demand, so there's no overhead if the filter apis are not called at all.
type RPCStream struct {
	evtClient rpcclient.EventsClient
	logger    log.Logger
	txDecoder sdk.TxDecoder

	// headerStream/logStream are backed by cometbft event subscription
	headerStream *Stream[RPCHeader]
	logStream    *Stream[*ethtypes.Log]

	// pendingTxStream is backed by check-tx ante handler
	pendingTxStream *Stream[common.Hash]

	wg               sync.WaitGroup
	validatorAccount validatorAccountFunc
}

func NewRPCStreams(
	evtClient rpcclient.EventsClient,
	logger log.Logger,
	txDecoder sdk.TxDecoder,
	validatorAccount validatorAccountFunc,
) *RPCStream {
	return &RPCStream{
		evtClient:        evtClient,
		logger:           logger,
		txDecoder:        txDecoder,
		validatorAccount: validatorAccount,
		pendingTxStream:  NewStream[common.Hash](txStreamSegmentSize, txStreamCapacity),
	}
}

func (s *RPCStream) initSubscriptions() {
	if s.headerStream != nil {
		// already initialized
		return
	}

	s.headerStream = NewStream[RPCHeader](headerStreamSegmentSize, headerStreamCapacity)
	s.logStream = NewStream[*ethtypes.Log](logStreamSegmentSize, logStreamCapacity)

	ctx := context.Background()

	chBlocks, err := s.evtClient.Subscribe(ctx, streamSubscriberName, blockEvents, subscribBufferSize)
	if err != nil {
		panic(err)
	}

	chLogs, err := s.evtClient.Subscribe(ctx, streamSubscriberName, evmEvents, subscribBufferSize)
	if err != nil {
		if err := s.evtClient.UnsubscribeAll(context.Background(), streamSubscriberName); err != nil {
			s.logger.Error("failed to unsubscribe", "err", err)
		}
		panic(err)
	}

	go s.start(&s.wg, chBlocks, chLogs)
}

func (s *RPCStream) Close() error {
	if s.headerStream == nil {
		// not initialized
		return nil
	}

	if err := s.evtClient.UnsubscribeAll(context.Background(), streamSubscriberName); err != nil {
		return err
	}
	s.wg.Wait()
	return nil
}

func (s *RPCStream) HeaderStream() *Stream[RPCHeader] {
	s.initSubscriptions()
	return s.headerStream
}

func (s *RPCStream) PendingTxStream() *Stream[common.Hash] {
	return s.pendingTxStream
}

func (s *RPCStream) LogStream() *Stream[*ethtypes.Log] {
	s.initSubscriptions()
	return s.logStream
}

// ListenPendingTx is a callback passed to application to listen for pending transactions in CheckTx.
func (s *RPCStream) ListenPendingTx(hash common.Hash) {
	s.PendingTxStream().Add(hash)
}

func (s *RPCStream) start(
	wg *sync.WaitGroup,
	chBlocks <-chan coretypes.ResultEvent,
	chLogs <-chan coretypes.ResultEvent,
) {
	wg.Add(1)
	defer func() {
		wg.Done()
		if err := s.evtClient.UnsubscribeAll(context.Background(), streamSubscriberName); err != nil {
			s.logger.Error("failed to unsubscribe", "err", err)
		}
	}()

	blockTimes := newBlockTimeTracker()
	for {
		select {
		case ev, ok := <-chBlocks:
			if !ok {
				chBlocks = nil
				// no more block times will come, release the logs waiting for them
				s.logStream.Add(blockTimes.flush()...)
				break
			}

			data, ok := ev.Data.(tmtypes.EventDataNewBlock)
			if !ok {
				s.logger.Error("event data type mismatch", "type", fmt.Sprintf("%T", ev.Data))
				continue
			}

			blockTime, err := ethermint.SafeUint64(data.Block.Time.Unix())
			if err != nil {
				blockTime = 0
			}
			s.logStream.Add(blockTimes.addBlock(data.Block.Height, blockTime)...)

			baseFee := types.BaseFeeFromEvents(data.ResultFinalizeBlock.Events)
			res, err := s.validatorAccount(
				types.ContextWithHeight(data.Block.Height),
				&evmtypes.QueryValidatorAccountRequest{
					ConsAddress: sdk.ConsAddress(data.Block.Header.ProposerAddress).String(),
				},
			)
			if err != nil {
				s.logger.Error("failed to get validator account", "err", err)
				continue
			}
			validator, err := sdk.AccAddressFromBech32(res.AccountAddress)
			if err != nil {
				s.logger.Error("failed to convert validator account", "err", err)
				continue
			}
			// TODO: fetch bloom from events
			header := types.EthHeaderFromTendermint(data.Block.Header, ethtypes.Bloom{}, baseFee, validator)
			txHash, err := evmTxHashFromEventData(data, s.txDecoder)
			if err != nil {
				// Drop rather than publish a wrong transactionsRoot; a gap in
				// the stream is better than incorrect data cached by clients.
				s.logger.Error("failed to compute transactionsRoot for newHeads, dropping header",
					"height", data.Block.Height, "err", err)
				continue
			}
			header.TxHash = txHash
			s.headerStream.Add(RPCHeader{EthHeader: header, Hash: common.BytesToHash(data.Block.Header.Hash())})

		case ev, ok := <-chLogs:
			if !ok {
				chLogs = nil
				break
			}

			if _, ok := ev.Events[evmTxHashKey]; !ok {
				// ignore transaction as it's not from the evm module
				continue
			}

			// get transaction result data
			dataTx, ok := ev.Data.(tmtypes.EventDataTx)
			if !ok {
				s.logger.Error("event data type mismatch", "type", fmt.Sprintf("%T", ev.Data))
				continue
			}
			height, err := ethermint.SafeUint64(dataTx.Height)
			if err != nil {
				continue
			}
			// block time is filled by the block time tracker, as tx events don't carry it
			txLogs, err := evmtypes.DecodeTxLogsFromEvents(dataTx.Result.Data, dataTx.Result.Events, height, 0)
			if err != nil {
				s.logger.Error("fail to decode evm tx response", "error", err.Error())
				continue
			}

			if chBlocks == nil {
				s.logStream.Add(txLogs...)
				continue
			}
			s.logStream.Add(blockTimes.addLogs(dataTx.Height, txLogs)...)
		}

		if chBlocks == nil && chLogs == nil {
			break
		}
	}
}

// blockTimeRetention is the number of recent block times kept to timestamp late tx events.
const blockTimeRetention = 128

// blockTimeTracker fills the block timestamp of the logs, the tx events don't carry the block time,
// and they can be received before or after the new block event of the same height,
// logs are held back until the block time of their height is known.
//
// Holding logs back normally adds no noticeable latency, as the new block event is published
// together with the tx events. If the new block event of a height is missed, e.g. dropped by a
// full subscription buffer, its logs are released without timestamp by the next block event,
// so they are delayed until the next block.
type blockTimeTracker struct {
	times   map[int64]uint64
	pending map[int64][]*ethtypes.Log
	// latest is the highest block height seen
	latest int64
}

func newBlockTimeTracker() *blockTimeTracker {
	return &blockTimeTracker{
		times:   make(map[int64]uint64),
		pending: make(map[int64][]*ethtypes.Log),
	}
}

// addBlock records the block time and returns the pending logs that are ready to publish,
// in block order.
func (t *blockTimeTracker) addBlock(height int64, blockTime uint64) []*ethtypes.Log {
	t.times[height] = blockTime
	t.latest = max(t.latest, height)
	for h := range t.times {
		if h <= t.latest-blockTimeRetention {
			delete(t.times, h)
		}
	}

	heights := make([]int64, 0, len(t.pending))
	for h := range t.pending {
		// logs of heights above are still waiting for their block event, the block event of the
		// heights below is missed, publish their logs without timestamp rather than holding them forever.
		if h <= height {
			heights = append(heights, h)
		}
	}
	slices.Sort(heights)

	var ready []*ethtypes.Log
	for _, h := range heights {
		logs := t.pending[h]
		if h == height {
			evmtypes.SetLogsBlockTimestamp(logs, blockTime)
		}
		ready = append(ready, logs...)
		delete(t.pending, h)
	}
	return ready
}

// addLogs returns the logs if they can be published, otherwise holds them back until the block
// time of their height is known.
func (t *blockTimeTracker) addLogs(height int64, logs []*ethtypes.Log) []*ethtypes.Log {
	if blockTime, ok := t.times[height]; ok {
		evmtypes.SetLogsBlockTimestamp(logs, blockTime)
		return logs
	}
	if height <= t.latest {
		// the block time is pruned or the block event is missed, it won't be known anymore.
		return logs
	}
	t.pending[height] = append(t.pending[height], logs...)
	return nil
}

// flush returns all the pending logs, in block order.
func (t *blockTimeTracker) flush() []*ethtypes.Log {
	heights := make([]int64, 0, len(t.pending))
	for h := range t.pending {
		heights = append(heights, h)
	}
	slices.Sort(heights)

	var logs []*ethtypes.Log
	for _, h := range heights {
		logs = append(logs, t.pending[h]...)
		delete(t.pending, h)
	}
	return logs
}

func evmTxHashFromEventData(data tmtypes.EventDataNewBlock, txDecoder sdk.TxDecoder) (common.Hash, error) {
	msgs, err := types.EvmMsgsFromTxs(txDecoder, data.Block.Txs, data.ResultFinalizeBlock.TxResults)
	if err != nil {
		return common.Hash{}, err
	}
	return types.EvmTxHashFromMsgs(msgs), nil
}
