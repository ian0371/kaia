package types

import (
	"sort"
)

type VotesInEpoch map[uint64]VoteData

type GovernanceCache struct {
	GroupedVotes map[uint64]VotesInEpoch
	Governances  map[uint64]GovernanceData
	GovHistory   GovernanceHistory
}

func NewGovernanceCache() *GovernanceCache {
	return &GovernanceCache{
		GroupedVotes: make(map[uint64]VotesInEpoch),
		Governances:  make(map[uint64]GovernanceData),
		GovHistory:   GovernanceHistory{},
	}
}

func (h *GovernanceCache) VoteBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.GroupedVotes {
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

func (h *GovernanceCache) AddVote(epochIdx, blockNum uint64, vote VoteData) {
	if _, ok := h.GroupedVotes[epochIdx]; !ok {
		h.GroupedVotes[epochIdx] = make(map[uint64]VoteData)
	}
	h.GroupedVotes[epochIdx][blockNum] = vote
}

func (h *GovernanceCache) AddGovernance(blockNum uint64, gov GovernanceData) {
	h.Governances[blockNum] = gov

	h.GovHistory = GetGovernanceHistory(h.Governances)
}

func (h *GovernanceCache) GetGovernanceHistory() GovernanceHistory {
	return h.GovHistory
}
