package types

import (
	"github.com/kaiachain/kaia/common"
)

type Param struct {
	Name  string
	Value interface{} // canonical value
}

type VoteData struct {
	BlockNum uint64
	Voter    common.Address
	Param    Param
}

type GovernanceData struct {
	BlockNum uint64
	Params   map[string]Param
}

type AllParamsHistory map[string]*PartitionList[Param] // p1 -> [(activation1, value1), ...], p2 -> ...

type GovernanceCache struct {
	Votes []VoteData
	Govs  []GovernanceData
}

func (h *GovernanceCache) VoteBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for _, vote := range h.Votes {
		blockNums = append(blockNums, vote.BlockNum)
	}
	return blockNums
}

func (h *GovernanceCache) GovBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for _, vote := range h.Govs {
		blockNums = append(blockNums, vote.BlockNum)
	}
	return blockNums
}

func (h *GovernanceCache) AddVote(num uint64, vote VoteData) {
	h.Votes = append(h.Votes, vote)
}

func (h *GovernanceCache) AddGov(num uint64, g GovernanceData) {
	h.Govs = append(h.Govs, g)
}

func GetAllParamsHistory(govs []GovernanceData, epoch uint64, koreHf uint64) AllParamsHistory {
	ret := make(AllParamsHistory)

	for _, g := range govs {
		for paramName, p := range g.Params {
			activation := uint64(0)
			if g.BlockNum < koreHf {
				activation = headerGovActivationBlockPreKore(g.BlockNum, epoch)
			} else {
				activation = headerGovActivationBlockPostKore(g.BlockNum, epoch)
			}

			if _, exists := ret[paramName]; !exists {
				ret[paramName] = &PartitionList[Param]{}
			}
			ret[paramName].AddRecord(uint(activation), p)
		}
	}

	return ret
}

func headerGovActivationBlockPreKore(govDataBlockNum, epoch uint64) uint64 {
	if govDataBlockNum == 0 {
		return 0
	}
	return govDataBlockNum + epoch + 1
}

func headerGovActivationBlockPostKore(govDataBlockNum, epoch uint64) uint64 {
	if govDataBlockNum == 0 {
		return 0
	}
	return govDataBlockNum + epoch
}
