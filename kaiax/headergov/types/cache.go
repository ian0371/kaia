package types

type GovernanceCache struct {
	Votes       map[uint64]VoteData
	Governances []GovernanceData
}

func (h *GovernanceCache) VoteBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.Votes {
		blockNums = append(blockNums, num)
	}
	return blockNums
}

func (h *GovernanceCache) GovBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for _, gov := range h.Governances {
		blockNums = append(blockNums, gov.BlockNum)
	}
	return blockNums
}

func (h *GovernanceCache) AddVote(num uint64, vote VoteData) {
	h.Votes[num] = vote
}

func (h *GovernanceCache) AddGovernance(num uint64, gov GovernanceData) {
	h.Governances = append(h.Governances, gov)
}

func (h *GovernanceCache) GetGovernanceHistory() GovernanceHistory {
	return GetGovernanceHistory(h.Governances)
}
