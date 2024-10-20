package ethapi

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/scroll-tech/go-ethereum"
	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/common/hexutil"
	"github.com/scroll-tech/go-ethereum/common/math"
	"github.com/scroll-tech/go-ethereum/core"
	"github.com/scroll-tech/go-ethereum/core/types"
	"github.com/scroll-tech/go-ethereum/core/vm"
	"github.com/scroll-tech/go-ethereum/log"
	"github.com/scroll-tech/go-ethereum/rpc"
)

type BatchTransactionArgs struct {
	Contract   *common.Address   `json:"contract"`
	CallParams []TransactionArgs `json:"callParams"`
}

func (s *PublicBlockChainAPI) EstimateGasWithLogs(ctx context.Context, args BatchTransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*ethereum.EstimateGasResult, error) {

	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}

	result, logs, err := DoCallWithLogs(ctx, s.b, args, bNrOrHash, nil, 0, s.b.RPCGasCap())
	if err != nil {
		return &ethereum.EstimateGasResult{Gas: 0, Logs: []*types.Log{}}, err
	} else {
		return &ethereum.EstimateGasResult{Gas: hexutil.Uint64(result.UsedGas), Logs: logs, Result: result.ReturnData}, err
	}

}

func DoCallWithLogs(ctx context.Context, b Backend, args BatchTransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride, timeout time.Duration, globalGasCap uint64) (*core.ExecutionResult, []*types.Log, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	state, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, nil, err
	}
	if err := overrides.Apply(state); err != nil {
		return nil, nil, err
	}

	header = types.CopyHeader(header)
	header.Time += 2
	header.Number = big.NewInt(0).Add(header.Number, big.NewInt(2))

	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	// Execute the message.
	gas := uint64(10000000)
	result := &core.ExecutionResult{}
	for i, item := range args.CallParams {

		item.Gas = (*hexutil.Uint64)(&gas)
		gp := new(core.GasPool).AddGas(math.MaxUint64)

		msg, err := item.ToMessage(globalGasCap, header.BaseFee)
		if err != nil {
			return nil, nil, err
		}

		evm, vmError, err := b.GetBatchEVM(ctx, msg, i, state, header, &vm.Config{NoBaseFee: true})
		if err != nil {
			return nil, nil, fmt.Errorf("err: load evm: %v ", err)
		}

		data, err := core.ApplyMessage(evm, msg, gp, common.Big0)
		if err != nil {
			return nil, nil, fmt.Errorf("err: apply message: %v ", err)
		}

		if err2 := vmError(); err2 != nil {
			return nil, nil, fmt.Errorf("err: emit evmerror: %v ", err2)
		}

		if data.Err != nil {
			return nil, nil, fmt.Errorf("err: inner data: %v", data.Err)
		}

		if state.Error() != nil {
			return nil, nil, fmt.Errorf("err: state data: %v", state.Error())
		}

		result = data
	}

	return result, state.GetSelfLogs(args.Contract), nil
}
