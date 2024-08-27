package types

import (
	"errors"
	"sort"
)

type GovernanceHistory map[uint64]GovernanceParam

func GetGovernanceHistory(govs map[uint64]GovernanceData) GovernanceHistory {
	gh := make(map[uint64]GovernanceParam)

	var sortedNums []uint64
	for num := range govs {
		sortedNums = append(sortedNums, num)
	}
	sort.Slice(sortedNums, func(i, j int) bool {
		return sortedNums[i] < sortedNums[j]
	})

	gp := GovernanceParam{}
	for _, num := range sortedNums {
		govData := govs[num]
		gp.SetFromGovernanceData(&govData)
		gh[num] = gp
	}
	return gh
}

func (g *GovernanceHistory) Search(blockNum uint64) (GovernanceParam, error) {
	idx := uint64(0)
	for num := range *g {
		if num > idx && num <= blockNum {
			idx = num
		}
	}
	if ret, ok := (*g)[idx]; ok {
		return ret, nil
	} else {
		return GovernanceParam{}, errors.New("blockNum not found from governance history")
	}
}
