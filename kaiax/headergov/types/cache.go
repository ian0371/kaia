package types

import (
	"sort"

	"github.com/kaiachain/kaia/log"
)

var logger = log.NewModuleLogger(log.KaiaXGov)

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
	logger.Warn("kaiax.VoteBlockNums", "blockNums", blockNums)
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
	logger.Warn("kaiax.GovBlockNums", "blockNums", blockNums)
	return blockNums
}

func (h *GovernanceCache) AddVote(num uint64, vote VoteData) {
	logger.Warn("kaiax.AddVote", "num", num, "vote", vote)
	h.Votes[num] = vote
}

func (h *GovernanceCache) AddGovernance(num uint64, gov GovernanceData) {
	logger.Warn("kaiax.AddGovernance", "num", num, "gov", gov)
	h.Governances[num] = gov
}

func (h *GovernanceCache) GetGovernanceHistory() GovernanceHistory {
	logger.Warn("kaiax.GetGovernanceHistory", "votes", h.Votes, "governances", h.Governances)
	return GetGovernanceHistory(h.Governances)
}
