package les

import (
	"context"

	"github.com/scroll-tech/go-ethereum/core"
	"github.com/scroll-tech/go-ethereum/core/state"
	"github.com/scroll-tech/go-ethereum/core/types"
	"github.com/scroll-tech/go-ethereum/core/vm"
)

func (b *LesApiBackend) GetBatchEVM(ctx context.Context, msg core.Message, index int, state *state.StateDB, header *types.Header, vmConfig *vm.Config) (*vm.EVM, func() error, error) {
	if vmConfig == nil {
		vmConfig = new(vm.Config)
	}
	parameters := &types.LegacyTx{
		Nonce:    msg.Nonce(),
		GasPrice: msg.GasPrice(),
		Gas:      msg.Gas(),
		To:       msg.To(),
		Value:    msg.Value(),
		Data:     msg.Data(),
	}
	transaction := types.NewTx(parameters)

	state.SetTxContext(transaction.Hash(), index)

	txContext := core.NewEVMTxContext(msg)
	context := core.NewEVMBlockContext(header, b.eth.BlockChain(), b.eth.blockchain.Config(), nil)
	return vm.NewEVM(context, txContext, state, b.eth.blockchain.Config(), *vmConfig), state.Error, nil
}
