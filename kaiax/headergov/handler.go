package headergov

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/params"
)

// TODO: verify Vote.
// TODO: filter valid votes (i.e., tally)
func (h *HeaderGovModule) VerifyHeader(header *types.Header) error {

	// verify Governance.
	if header.Number.Uint64() == 0 {
		return nil
	}

	// 1. epoch check
	if header.Number.Uint64()%h.epoch != 0 {
		if len(header.Governance) > 0 {
			return errors.New("governance is not allowed in non-epoch block")
		} else {
			return nil
		}
	}

	// 2. vote pass check
	votes := h.getVotesInEpoch(calcEpochIdx(header.Number.Uint64()-1, h.epoch))
	expectedParamSet := new(params.GovParamSet)
	for _, vote := range votes {
		voteParamSet, err := vote.ToParamSet()
		if err != nil {
			return err
		}
		expectedParamSet = params.NewGovParamSetMerged(expectedParamSet, voteParamSet)
	}

	deserializedGov, err := headergov_types.DeserializeHeaderGov(header.Governance, header.Number.Uint64())
	if err != nil {
		return err
	}
	headerParamSet, err := deserializedGov.ToParamSet()
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expectedParamSet, headerParamSet) {
		return fmt.Errorf("expected ParamSet: %v, header ParamSet: %v", expectedParamSet, headerParamSet)
	}

	return nil
}

// TODO: Add Gov to header.
// TODO: if myVote exists, put Vote to header.
func (h *HeaderGovModule) PrepareHeader(header *types.Header) (*types.Header, error) {
	// if epoch block & vote exists in the last epoch, put Governance to header.
	if header.Number.Uint64()%h.epoch == 0 {
	}

	return header, nil // TODO: implement
}

func (h *HeaderGovModule) FinalizeBlock(b *types.Block) (*types.Block, error) {
	return b, nil
}

func (h *HeaderGovModule) PostInsertBlock(b *types.Block) error {
	if len(b.Header().Vote) > 0 {
		vote, err := headergov_types.DeserializeHeaderVote(b.Header().Vote, b.NumberU64())
		if err != nil {
			return err
		}
		h.AddVote(vote)
	}

	if len(b.Header().Governance) > 0 {
		gov, err := headergov_types.DeserializeHeaderGov(b.Header().Vote, b.NumberU64())
		if err != nil {
			return err
		}
		h.AddGov(gov)
	}

	return nil
}

func (h *HeaderGovModule) AddVote(vote *VoteData) error {
	h.cache.AddVote(vote.BlockNum, *vote)

	var data StoredVoteBlockNums = h.cache.VoteBlockNums()
	WriteVoteDataBlocks(h.ChainKv, &data)
	return nil
}

func (h *HeaderGovModule) AddGov(gov *GovernanceData) error {
	h.cache.AddGov(gov.BlockNum, *gov)

	var data StoredGovBlockNums = h.cache.GovBlockNums()
	WriteGovDataBlocks(h.ChainKv, &data)
	return nil
}

func (h *HeaderGovModule) getVotesInEpoch(epochIdx uint64) []VoteData {
	ret := make([]VoteData, 0)
	for _, vote := range h.cache.Votes {
		if calcEpochIdx(vote.BlockNum, h.epoch) == epochIdx {
			ret = append(ret, vote)
		}
	}

	return ret
}

func calcEpochIdx(blockNum uint64, epoch uint64) uint64 {
	return blockNum / epoch
}
