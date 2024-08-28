package types

import (
	"sort"
)

type GovernanceCache struct {
	Votes       map[uint64]VoteData
	Governances map[uint64]GovernanceData
}

func (h *GovernanceCache) VoteBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.Votes {
		blockNums = append(blockNums, num)
	}
	sort.Slice(blockNums, func(i, j int) bool {
		return blockNums[i] < blockNums[j]
	})
	return blockNums
}

func (h *GovernanceCache) GovBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.Governances {
		blockNums = append(blockNums, num)
	}
	sort.Slice(blockNums, func(i, j int) bool {
		return blockNums[i] < blockNums[j]
	})
	return blockNums
}

func (h *GovernanceCache) AddVote(num uint64, vote VoteData) {
	h.Votes[num] = vote
}

func (h *GovernanceCache) AddGovernance(num uint64, gov GovernanceData) {
	h.Governances[num] = gov
}

func (h *GovernanceCache) GetGovernanceHistory() GovernanceHistory {
	return GetGovernanceHistory(h.Governances)
}
