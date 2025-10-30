package backend

import (
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	tmrpcclient "github.com/cometbft/cometbft/rpc/client"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/rpc/backend/mocks"
	rpc "github.com/evmos/ethermint/rpc/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/gogo/protobuf/proto"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Client defines a mocked object that implements the Tendermint JSON-RPC Client
// interface. It allows for performing Client queries without having to run a
// Tendermint RPC Client server.
//
// To use a mock method it has to be registered in a given test.
var _ tmrpcclient.Client = &mocks.Client{}

// Tx Search
func RegisterTxSearch(client *mocks.Client, query string, txBz []byte) {
	resulTxs := []*tmrpctypes.ResultTx{{Tx: txBz}}
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(&tmrpctypes.ResultTxSearch{Txs: resulTxs, TotalCount: 1}, nil).
		Maybe()
}

func RegisterTxSearchEmpty(client *mocks.Client, query string) {
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(&tmrpctypes.ResultTxSearch{}, nil).
		Maybe()
}

func RegisterTxSearchError(client *mocks.Client, query string) {
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

// Broadcast Tx
func RegisterBroadcastTx(client *mocks.Client, tx types.Tx) {
	client.On("BroadcastTxSync", context.Background(), tx).
		Return(&tmrpctypes.ResultBroadcastTx{}, nil).
		Maybe()
}

func RegisterBroadcastTxError(client *mocks.Client, tx types.Tx) {
	client.On("BroadcastTxSync", context.Background(), tx).
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

// Unconfirmed Transactions
func RegisterUnconfirmedTxs(client *mocks.Client, limit *int, txs []types.Tx) {
	client.On("UnconfirmedTxs", rpc.ContextWithHeight(1), limit).
		Return(&tmrpctypes.ResultUnconfirmedTxs{Txs: txs}, nil).
		Maybe()
}

func RegisterUnconfirmedTxsEmpty(client *mocks.Client, limit *int) {
	client.On("UnconfirmedTxs", rpc.ContextWithHeight(1), limit).
		Return(&tmrpctypes.ResultUnconfirmedTxs{
			Txs: make([]types.Tx, 2),
		}, nil).
		Maybe()
}

func RegisterUnconfirmedTxsError(client *mocks.Client, limit *int) {
	client.On("UnconfirmedTxs", rpc.ContextWithHeight(1), limit).
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

// Status
func RegisterStatus(client *mocks.Client) {
	client.On("Status", rpc.ContextWithHeight(1)).
		Return(&tmrpctypes.ResultStatus{}, nil).
		Maybe()
}

func RegisterStatusError(client *mocks.Client) {
	client.On("Status", rpc.ContextWithHeight(1)).
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

// Block
func RegisterBlockMultipleTxs(
	client *mocks.Client,
	height int64,
	txs []types.Tx,
) (*tmrpctypes.ResultBlock, error) {
	block := types.MakeBlock(height, txs, nil, nil)
	block.ChainID = ChainID
	resBlock := &tmrpctypes.ResultBlock{Block: block}
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).Return(resBlock, nil).
		Maybe()
	return resBlock, nil
}

func RegisterBlock(
	client *mocks.Client,
	height int64,
	tx []byte,
) (*tmrpctypes.ResultBlock, error) {
	// without tx
	if tx == nil {
		emptyBlock := types.MakeBlock(height, []types.Tx{}, nil, nil)
		emptyBlock.ChainID = ChainID
		resBlock := &tmrpctypes.ResultBlock{Block: emptyBlock}
		client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).Return(resBlock, nil).
			Maybe()
		return resBlock, nil
	}

	// with tx
	block := types.MakeBlock(height, []types.Tx{tx}, nil, nil)
	block.ChainID = ChainID
	resBlock := &tmrpctypes.ResultBlock{Block: block}
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).Return(resBlock, nil).
		Maybe()
	return resBlock, nil
}

// Block returns error
func RegisterBlockError(client *mocks.Client, height int64) {
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

// Block not found
func RegisterBlockNotFound(
	client *mocks.Client,
	height int64,
) (*tmrpctypes.ResultBlock, error) {
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(&tmrpctypes.ResultBlock{Block: nil}, nil).
		Maybe()

	return &tmrpctypes.ResultBlock{Block: nil}, nil
}

// Block panic
func RegisterBlockPanic(client *mocks.Client, height int64) {
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(func(context.Context, *int64) *tmrpctypes.ResultBlock {
			panic("Block call panic")
		}, nil).
		Maybe()
}

func TestRegisterBlock(t *testing.T) {
	client := mocks.NewClient(t)
	height := rpc.BlockNumber(1).Int64()
	RegisterBlock(client, height, nil)

	res, err := client.Block(rpc.ContextWithHeight(height), &height)

	emptyBlock := types.MakeBlock(height, []types.Tx{}, nil, nil)
	emptyBlock.ChainID = ChainID
	resBlock := &tmrpctypes.ResultBlock{Block: emptyBlock}
	require.Equal(t, resBlock, res)
	require.NoError(t, err)
}

// ConsensusParams
func RegisterConsensusParams(client *mocks.Client, height int64) {
	consensusParams := types.DefaultConsensusParams()
	client.On("ConsensusParams", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(&tmrpctypes.ResultConsensusParams{ConsensusParams: *consensusParams}, nil).
		Maybe()
}

func RegisterConsensusParamsError(client *mocks.Client, height int64) {
	client.On("ConsensusParams", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

func TestRegisterConsensusParams(t *testing.T) {
	client := mocks.NewClient(t)
	height := int64(1)
	RegisterConsensusParams(client, height)

	res, err := client.ConsensusParams(rpc.ContextWithHeight(height), &height)
	consensusParams := types.DefaultConsensusParams()
	require.Equal(t, &tmrpctypes.ResultConsensusParams{ConsensusParams: *consensusParams}, res)
	require.NoError(t, err)
}

// BlockResults

func RegisterBlockResultsWithEventLog(client *mocks.Client, height int64) (*tmrpctypes.ResultBlockResults, error) {
	any, err := codectypes.NewAnyWithValue(&evmtypes.MsgEthereumTxResponse{
		Logs: []*evmtypes.Log{
			{Data: []byte("data")},
		},
	})
	if err != nil {
		return nil, err
	}
	data, err := proto.Marshal(&sdk.TxMsgData{MsgResponses: []*codectypes.Any{any}})
	if err != nil {
		return nil, err
	}

	res := &tmrpctypes.ResultBlockResults{
		Height: height,
		TxsResults: []*abci.ExecTxResult{
			{Code: 0, GasUsed: 0, Data: data},
		},
	}

	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(res, nil).
		Maybe()
	return res, nil
}

func RegisterBlockResults(
	client *mocks.Client,
	height int64,
) (*tmrpctypes.ResultBlockResults, error) {
	res := &tmrpctypes.ResultBlockResults{
		Height:     height,
		TxsResults: []*abci.ExecTxResult{{Code: 0, GasUsed: 0}},
	}

	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(res, nil).
		Maybe()
	return res, nil
}

func RegisterBlockResultsError(client *mocks.Client, height int64) {
	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

func TestRegisterBlockResults(t *testing.T) {
	client := mocks.NewClient(t)
	height := int64(1)
	RegisterBlockResults(client, height)

	res, err := client.BlockResults(rpc.ContextWithHeight(height), &height)
	expRes := &tmrpctypes.ResultBlockResults{
		Height:     height,
		TxsResults: []*abci.ExecTxResult{{Code: 0, GasUsed: 0}},
	}
	require.Equal(t, expRes, res)
	require.NoError(t, err)
}

// BlockByHash
func RegisterBlockByHash(
	client *mocks.Client,
	hash common.Hash,
	tx []byte,
) (*tmrpctypes.ResultBlock, error) {
	block := types.MakeBlock(1, []types.Tx{tx}, nil, nil)
	resBlock := &tmrpctypes.ResultBlock{Block: block}

	client.On("BlockByHash", rpc.ContextWithHeight(1), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(resBlock, nil).
		Maybe()
	return resBlock, nil
}

func RegisterBlockByHashError(client *mocks.Client, hash common.Hash, tx []byte) {
	client.On("BlockByHash", rpc.ContextWithHeight(1), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

func RegisterBlockByHashNotFound(client *mocks.Client, hash common.Hash, tx []byte) {
	client.On("BlockByHash", rpc.ContextWithHeight(1), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(nil, nil).
		Maybe()
}

// HeaderByHash
func RegisterHeaderByHash(
	client *mocks.Client,
	hash common.Hash,
	tx []byte,
) (*tmrpctypes.ResultHeader, error) {
	block := types.MakeBlock(1, []types.Tx{tx}, nil, nil)
	resHeader := &tmrpctypes.ResultHeader{Header: &block.Header}

	client.On("HeaderByHash", rpc.ContextWithHeight(1), bytes.HexBytes(hash.Bytes())).
		Return(resHeader, nil).
		Maybe()
	return resHeader, nil
}

func RegisterHeaderByHashError(client *mocks.Client, hash common.Hash, tx []byte) {
	client.On("HeaderByHash", rpc.ContextWithHeight(1), bytes.HexBytes(hash.Bytes())).
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

func RegisterHeaderByHashNotFound(client *mocks.Client, hash common.Hash, tx []byte) {
	client.On("HeaderByHash", rpc.ContextWithHeight(1), bytes.HexBytes(hash.Bytes())).
		Return(&tmrpctypes.ResultHeader{Header: nil}, nil).
		Maybe()
}

func RegisterHeaderByHashAny(client *mocks.Client) {
	client.On("HeaderByHash", mock.Anything, mock.Anything).
		Return(&tmrpctypes.ResultHeader{Header: &types.Header{Height: 1}}, nil).Maybe()
}

// Header
func RegisterHeader(client *mocks.Client, height *int64, tx []byte) (*tmrpctypes.ResultHeader, error) {
	block := types.MakeBlock(*height, []types.Tx{tx}, nil, nil)
	resHeader := &tmrpctypes.ResultHeader{Header: &block.Header}
	client.On("Header", rpc.ContextWithHeight(*height), height).Return(resHeader, nil).
		Maybe()
	return resHeader, nil
}

func RegisterHeaderError(client *mocks.Client, height *int64) {
	client.On("Header", rpc.ContextWithHeight(*height), height).Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

// Header not found
func RegisterHeaderNotFound(client *mocks.Client, height int64) {
	client.On("Header", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(&tmrpctypes.ResultHeader{Header: nil}, nil).
		Maybe()
}

func RegisterABCIQueryWithOptions(client *mocks.Client, height int64, path string, data bytes.HexBytes, opts tmrpcclient.ABCIQueryOptions) {
	client.On("ABCIQueryWithOptions", context.Background(), path, data, opts).
		Return(&tmrpctypes.ResultABCIQuery{
			Response: abci.ResponseQuery{
				Value:  []byte{2}, // TODO replace with data.Bytes(),
				Height: height,
			},
		}, nil).
		Maybe()
}

func RegisterABCIQueryWithOptionsError(clients *mocks.Client, path string, data bytes.HexBytes, opts tmrpcclient.ABCIQueryOptions) {
	clients.On("ABCIQueryWithOptions", context.Background(), path, data, opts).
		Return(nil, errortypes.ErrInvalidRequest).
		Maybe()
}

func RegisterABCIQueryAny(client *mocks.Client, height int64) {
	client.On("ABCIQueryWithOptions", context.Background(), mock.Anything, mock.Anything, mock.Anything).
		Return(&tmrpctypes.ResultABCIQuery{
			Response: abci.ResponseQuery{
				Value:  []byte{2}, // TODO replace with data.Bytes(),
				Height: height,
			},
		}, nil).
		Maybe()

}

func RegisterABCIQueryAccount(clients *mocks.Client, data bytes.HexBytes, opts tmrpcclient.ABCIQueryOptions, acc client.Account) {
	baseAccount := authtypes.NewBaseAccount(acc.GetAddress(), acc.GetPubKey(), acc.GetAccountNumber(), acc.GetSequence())
	accAny, _ := codectypes.NewAnyWithValue(baseAccount)
	accResponse := authtypes.QueryAccountResponse{Account: accAny}
	respBz, _ := accResponse.Marshal()
	clients.On("ABCIQueryWithOptions", context.Background(), "/cosmos.auth.v1beta1.Query/Account", data, opts).
		Return(&tmrpctypes.ResultABCIQuery{
			Response: abci.ResponseQuery{
				Value:  respBz,
				Height: 1,
			},
		}, nil).
		Maybe()
}
