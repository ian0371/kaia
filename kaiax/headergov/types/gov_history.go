package types

import (
	"errors"
	"sort"
)

type GovHistory map[uint64]GovParamSet

func GetGovHistory(govs map[uint64]GovData) GovHistory {
	gh := make(map[uint64]GovParamSet)

	var sortedNums []uint64
	for num := range govs {
		sortedNums = append(sortedNums, num)
	}
	sort.Slice(sortedNums, func(i, j int) bool {
		return sortedNums[i] < sortedNums[j]
	})

	gp := GovParamSet{}
	for _, num := range sortedNums {
		govData := govs[num]
		gp.SetFromGovernanceData(&govData)
		gh[num] = gp
	}
	return gh
}

func (g *GovHistory) Search(blockNum uint64) (GovParamSet, error) {
	idx := uint64(0)
	for num := range *g {
		if num > idx && num <= blockNum {
			idx = num
		}
	}
	if ret, ok := (*g)[idx]; ok {
		return ret, nil
	} else {
		return GovParamSet{}, errors.New("blockNum not found from governance history")
	}
}
