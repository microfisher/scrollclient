package state

import (
	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/core/types"
)

func (s *StateDB) GetSelfLogs(contract *common.Address) []*types.Log {

	if contract == nil {
		return nil
	}

	logs := s.logs[s.thash]
	items := make([]*types.Log, 0, len(logs))
	for _, item := range logs {
		if item.Address == *contract {
			items = append(items, item)
		}
	}
	return items

}
