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

func (h *HeaderGovModule) HandleVote(blockNum uint64, vote VoteData) error {
	h.cache.AddVote(calcEpochIdx(blockNum, h.epoch), blockNum, vote)

	var data StoredUint64Array = h.cache.VoteBlockNums()
	WriteGovVoteDataBlockNums(h.ChainKv, &data)

	for i, myvote := range h.myVotes {
		if reflect.DeepEqual(&myvote, vote) {
			h.PopMyVotes(i)
			break
		}
	}

	return nil
}

func (h *HeaderGovModule) HandleGov(blockNum uint64, gov *GovData) error {
	h.cache.AddGov(blockNum, *gov)

	// merge gov based on latest effective params.
	gp, err := h.EffectiveParams(blockNum)
	if err != nil {
		logger.Error("kaiax.HandleGov error fetching EffectiveParams", "blockNum", blockNum, "gov", *gov, "err", err)
		return err
	}

	err = gp.SetFromGovernanceData(gov)
	if err != nil {
		logger.Error("kaiax.HandleGov error setting paramset", "blockNum", blockNum, "gov", *gov, "err", err, "gp", gp)
		return err
	}

	var data StoredUint64Array = h.cache.GovBlockNums()
	WriteGovDataBlockNums(h.ChainKv, &data)
	return nil
}

func calcEpochIdx(blockNum uint64, epoch uint64) uint64 {
	return blockNum / epoch
}
