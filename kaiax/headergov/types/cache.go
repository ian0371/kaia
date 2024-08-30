package types

import (
	"sort"
)

type VotesInEpoch map[uint64]VoteData

type GovHeaderCache struct {
	groupedVotes map[uint64]VotesInEpoch
	governances  map[uint64]GovData
	govHistory   GovHistory
}

func NewGovernanceCache() *GovHeaderCache {
	return &GovHeaderCache{
		groupedVotes: make(map[uint64]VotesInEpoch),
		governances:  make(map[uint64]GovData),
		govHistory:   GovHistory{},
	}
}

func (h *GovHeaderCache) GroupedVotes() map[uint64]VotesInEpoch {
	votes := make(map[uint64]VotesInEpoch)
	for epochIdx, votesInEpoch := range h.groupedVotes {
		votes[epochIdx] = make(VotesInEpoch)
		for blockNum, vote := range votesInEpoch {
			votes[epochIdx][blockNum] = vote
		}
	}
	return votes
}

func (h *GovHeaderCache) Govs() map[uint64]GovData {
	govs := make(map[uint64]GovData)
	for blockNum, gov := range h.governances {
		govs[blockNum] = gov
	}
	return govs
}

func (h *GovHeaderCache) VoteBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.groupedVotes {
		blockNums = append(blockNums, num)
	}
	sort.Slice(blockNums, func(i, j int) bool {
		return blockNums[i] < blockNums[j]
	})
	return blockNums
}

func (h *GovHeaderCache) GovBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.governances {
		blockNums = append(blockNums, num)
	}
	sort.Slice(blockNums, func(i, j int) bool {
		return blockNums[i] < blockNums[j]
	})
	return blockNums
}

func (h *GovHeaderCache) AddVote(epochIdx, blockNum uint64, vote VoteData) {
	if _, ok := h.groupedVotes[epochIdx]; !ok {
		h.groupedVotes[epochIdx] = make(map[uint64]VoteData)
	}
	h.groupedVotes[epochIdx][blockNum] = vote
}

func (h *GovHeaderCache) AddGovernance(blockNum uint64, gov GovData) {
	h.governances[blockNum] = gov

	h.govHistory = GetGovHistory(h.governances)
}

// TODO-kaiax: rename
func (h *GovHeaderCache) GetGovernanceHistory() GovHistory {
	return h.govHistory
}
