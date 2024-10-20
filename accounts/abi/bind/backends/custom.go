package backends

import (
	"github.com/scroll-tech/go-ethereum/core/types"
	"github.com/scroll-tech/go-ethereum/params"
)

func (fb *filterBackend) CurrentHeader() *types.Header {
	return fb.bc.CurrentHeader()
}
func (fb *filterBackend) ChainConfig() *params.ChainConfig {
	return fb.bc.Config()
}
