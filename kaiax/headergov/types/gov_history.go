package types

import (
	"errors"
	"sort"
)

type History map[uint64]ParamSet

func GetHistory(govs map[uint64]GovData) History {
	gh := make(map[uint64]ParamSet)

	var sortedNums []uint64
	for num := range govs {
		sortedNums = append(sortedNums, num)
	}
	sort.Slice(sortedNums, func(i, j int) bool {
		return sortedNums[i] < sortedNums[j]
	})

	gp := ParamSet{}
	for _, num := range sortedNums {
		govData := govs[num]
		gp.SetFromGovernanceData(govData)
		gh[num] = gp
	}
	return gh
}

func (g *History) Search(blockNum uint64) (ParamSet, error) {
	idx := uint64(0)
	for num := range *g {
		if num > idx && num <= blockNum {
			idx = num
		}
	}
	if ret, ok := (*g)[idx]; ok {
		return ret, nil
	} else {
		return ParamSet{}, errors.New("blockNum not found from governance history")
	}
}
