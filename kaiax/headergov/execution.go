package headergov

import (
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
)

func (h *HeaderGovModule) PostInsertBlock(b *types.Block) error {
	if len(b.Header().Vote) > 0 {
		vote, err := headergov_types.DeserializeHeaderVote(b.Header().Vote, b.Number().Uint64())
		if err != nil {
			logger.Error("kaiax.PostInsertBlock error", "vote", b.Header().Vote, "err", err)
			return err
		}
		h.HandleVote(b.NumberU64(), vote)
	}

	if len(b.Header().Governance) > 0 {
		gov, err := headergov_types.DeserializeHeaderGov(b.Header().Governance, b.NumberU64())
		if err != nil {
			logger.Error("kaiax.PostInsertBlock error", "governance", b.Header().Governance, "err", err)
			return err
		}
		h.HandleGov(b.NumberU64(), gov)
	}

	return nil
}

func (h *HeaderGovModule) HandleVote(blockNum uint64, vote *VoteData) error {
	h.cache.AddVote(calcEpochIdx(blockNum, h.epoch), blockNum, *vote)

	var data StoredVoteBlockNums = h.cache.VoteBlockNums()
	WriteVoteDataBlockNums(h.ChainKv, &data)

	for i, myvote := range h.MyVotes {
		if reflect.DeepEqual(myvote, vote) {
			h.MyVotes = append(h.MyVotes[:i], h.MyVotes[i+1:]...)
		}
	}

	return nil
}

func (h *HeaderGovModule) HandleGov(blockNum uint64, gov *GovernanceData) error {
	h.cache.AddGovernance(blockNum, *gov)

	// merge gov based on latest effective params.
	gp, err := h.EffectiveParams(blockNum)
	if err != nil {
		logger.Error("kaiax.HandleGov error", "blockNum", blockNum, "gov", *gov, "err", err)
		return err
	}

	gp.SetFromGovernanceData(gov)
	var data StoredGovBlockNums = h.cache.GovBlockNums()
	WriteGovDataBlockNums(h.ChainKv, &data)
	return nil
}

func calcEpochIdx(blockNum uint64, epoch uint64) uint64 {
	return blockNum / epoch
}
