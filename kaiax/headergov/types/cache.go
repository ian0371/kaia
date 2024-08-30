package types

import (
	"sort"
)

type VotesInEpoch map[uint64]VoteData

type GovHeaderCache struct {
	groupedGovVotes map[uint64]VotesInEpoch
	valVotes        VotesInEpoch
	governances     map[uint64]GovData
	govHistory      GovHistory
}

func NewHeaderGovCache() *GovHeaderCache {
	return &GovHeaderCache{
		groupedGovVotes: make(map[uint64]VotesInEpoch),
		governances:     make(map[uint64]GovData),
		govHistory:      GovHistory{},
	}
}

func (h *GovHeaderCache) GroupedGovVotes() map[uint64]VotesInEpoch {
	votes := make(map[uint64]VotesInEpoch)
	for epochIdx, votesInEpoch := range h.groupedGovVotes {
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

func (h *GovHeaderCache) GovHistory() GovHistory {
	return h.govHistory
}

func (h *GovHeaderCache) GovVoteBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.groupedGovVotes {
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

func (h *GovHeaderCache) AddGovVote(epochIdx, blockNum uint64, vote VoteData) {
	if _, ok := h.groupedGovVotes[epochIdx]; !ok {
		h.groupedGovVotes[epochIdx] = make(map[uint64]VoteData)
	}
	h.groupedGovVotes[epochIdx][blockNum] = vote
}

func (h *GovHeaderCache) AddGov(blockNum uint64, gov GovData) {
	h.governances[blockNum] = gov

	h.govHistory = GetGovHistory(h.governances)
}

func (h *GovHeaderCache) RemoveVotesAfter(blockNum uint64) {
	for epochIdxIter, votes := range h.groupedGovVotes {
		for blockNumIter := range votes {
			if blockNumIter > blockNum {
				delete(h.groupedGovVotes[epochIdxIter], blockNumIter)
			}
		}
		// If all votes for this epoch have been removed, delete the epoch entry
		if len(h.groupedGovVotes[epochIdxIter]) == 0 {
			delete(h.groupedGovVotes, epochIdxIter)
		}
	}
}

func (h *GovHeaderCache) RemoveGovernanceAfter(blockNum uint64) {
	for blockNumIter := range h.governances {
		if blockNumIter > blockNum {
			delete(h.governances, blockNumIter)
		}
	}

	// Regenerate the governance history after removing entries
	h.govHistory = GetGovHistory(h.governances)
}
