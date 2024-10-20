package filters

import (
	"math/big"

	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/consensus/misc"
	"github.com/scroll-tech/go-ethereum/core/types"
	"github.com/scroll-tech/go-ethereum/internal/ethapi"
	"github.com/scroll-tech/go-ethereum/params"
)

func NewRPCPendingTransaction(tx *types.Transaction, current *types.Header, config *params.ChainConfig) *ethapi.RPCTransaction {
	var baseFee *big.Int
	blockNumber := uint64(0)

	if current != nil {
		baseFee = misc.CalcBaseFee(config, current, big.NewInt(0))
		blockNumber = current.Number.Uint64()
	}
	return ethapi.NewRPCTransaction(tx, common.Hash{}, blockNumber, 0, baseFee, config)
}
