package ethclient

import (
	"context"

	"github.com/scroll-tech/go-ethereum"
	"github.com/scroll-tech/go-ethereum/common/hexutil"
)

func (ec *Client) EstimateGasWithLogs(ctx context.Context, msg ethereum.BatchCallMsg) (ethereum.EstimateGasResult, error) {
	var result ethereum.EstimateGasResult
	err := ec.c.CallContext(ctx, &result, "eth_estimateGasWithLogs", toBatchCallArg(msg))
	if err != nil {
		return result, err
	}
	return result, nil
}

func toBatchCallArg(call ethereum.BatchCallMsg) interface{} {
	items := make([]interface{}, 0, len(call.CallParams))
	for _, msg := range call.CallParams {
		arg := map[string]interface{}{
			"from": msg.From,
			"to":   msg.To,
		}
		if len(msg.Data) > 0 {
			arg["data"] = hexutil.Bytes(msg.Data)
		}
		if msg.Value != nil {
			arg["value"] = (*hexutil.Big)(msg.Value)
		}
		if msg.Gas != 0 {
			arg["gas"] = hexutil.Uint64(msg.Gas)
		}
		if msg.GasPrice != nil {
			arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
		}
		if msg.GasFeeCap != nil {
			arg["maxFeePerGas"] = (*hexutil.Big)(msg.GasFeeCap)
		}
		if msg.GasTipCap != nil {
			arg["maxPriorityFeePerGas"] = (*hexutil.Big)(msg.GasTipCap)
		}
		if msg.AccessList != nil {
			arg["accessList"] = msg.AccessList
		}
		if msg.BlobGasFeeCap != nil {
			arg["maxFeePerBlobGas"] = (*hexutil.Big)(msg.BlobGasFeeCap)
		}
		if msg.BlobHashes != nil {
			arg["blobVersionedHashes"] = msg.BlobHashes
		}
		items = append(items, arg)
	}
	batches := make(map[string]interface{}, len(call.CallParams))
	if call.Contract != nil {
		batches["contract"] = *call.Contract
	}
	batches["callParams"] = items
	return batches
}
