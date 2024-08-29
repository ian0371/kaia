package types

import (
	"sort"
)

type VotesInEpoch map[uint64]VoteData

type GovernanceCache struct {
	groupedVotes map[uint64]VotesInEpoch
	governances  map[uint64]GovernanceData
	govHistory   GovernanceHistory
}

func NewGovernanceCache() *GovernanceCache {
	return &GovernanceCache{
		groupedVotes: make(map[uint64]VotesInEpoch),
		governances:  make(map[uint64]GovernanceData),
		govHistory:   GovernanceHistory{},
	}
}

func (h *GovernanceCache) GroupedVotes() map[uint64]VotesInEpoch {
	votes := make(map[uint64]VotesInEpoch)
	for epochIdx, votesInEpoch := range h.groupedVotes {
		votes[epochIdx] = make(VotesInEpoch)
		for blockNum, vote := range votesInEpoch {
			votes[epochIdx][blockNum] = vote
		}
	}
	return votes
}

func (h *GovernanceCache) Govs() map[uint64]GovernanceData {
	govs := make(map[uint64]GovernanceData)
	for blockNum, gov := range h.governances {
		govs[blockNum] = gov
	}
	return govs
}

func (h *GovernanceCache) VoteBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.groupedVotes {
		blockNums = append(blockNums, num)
	}
	sort.Slice(blockNums, func(i, j int) bool {
		return blockNums[i] < blockNums[j]
	})
	return blockNums
}

func (h *GovernanceCache) GovBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.governances {
		blockNums = append(blockNums, num)
	}
	sort.Slice(blockNums, func(i, j int) bool {
		return blockNums[i] < blockNums[j]
	})
	return blockNums
}

func (h *GovernanceCache) AddVote(epochIdx, blockNum uint64, vote VoteData) {
	if _, ok := h.groupedVotes[epochIdx]; !ok {
		h.groupedVotes[epochIdx] = make(map[uint64]VoteData)
	}
	h.groupedVotes[epochIdx][blockNum] = vote
}

func (h *GovernanceCache) AddGovernance(blockNum uint64, gov GovernanceData) {
	h.governances[blockNum] = gov

	h.govHistory = GetGovernanceHistory(h.governances)
}

// TODO-kaiax: rename
func (h *GovernanceCache) GetGovernanceHistory() GovernanceHistory {
	return h.govHistory
}
