package headergov

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/params"
)

func (h *HeaderGovModule) VerifyHeader(header *types.Header) error {
	if header.Number.Uint64() == 0 {
		return nil
	}

	// 1. Check Vote
	if header.Vote != nil {
		vote, err := headergov_types.DeserializeHeaderVote(header.Vote, header.Number.Uint64())
		if err != nil {
			return err
		}
		err = h.VerifyVote(vote)
		if err != nil {
			return err
		}
	}

	// 2. Check Governance
	if header.Number.Uint64()%h.epoch != 0 {
		if len(header.Governance) == 0 {
			return nil
		} else {
			return errors.New("governance is not allowed in non-epoch block")
		}
	} else {
		expected := h.getExpectedGovernance(calcEpochIdx(header.Number.Uint64(), h.epoch))
		actual, err := headergov_types.DeserializeHeaderGov(header.Governance, header.Number.Uint64())
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(expected, actual) {
			return fmt.Errorf("expected governance: %v, actual: %v", expected, actual)
		}

		return nil
	}
}

func (h *HeaderGovModule) PrepareHeader(header *types.Header) (*types.Header, error) {
	// if epoch block & vote exists in the last epoch, put Governance to header.
	if len(h.MyVotes) > 0 {
		header.Vote, _ = h.MyVotes[0].Serialize()
		h.MyVotes = h.MyVotes[1:]
	}

	if header.Number.Uint64()%h.epoch == 0 {
		gov := h.getExpectedGovernance(calcEpochIdx(header.Number.Uint64(), h.epoch))
		header.Governance, _ = gov.Serialize()
	}

	return header, nil
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
		h.HandleVote(vote)
	}

	if len(b.Header().Governance) > 0 {
		gov, err := headergov_types.DeserializeHeaderGov(b.Header().Vote, b.NumberU64())
		if err != nil {
			return err
		}
		h.HandleGov(gov)
	}

	return nil
}

func (h *HeaderGovModule) HandleVote(vote *VoteData) error {
	h.cache.AddVote(vote.BlockNum, *vote)

	var data StoredVoteBlockNums = h.cache.VoteBlockNums()
	WriteVoteDataBlockNums(h.ChainKv, &data)
	return nil
}

func (h *HeaderGovModule) HandleGov(gov *GovernanceData) error {
	h.cache.AddGovernance(gov.BlockNum, *gov)

	// merge gov based on latest effective params.
	gp, err := h.EffectiveParams(gov.BlockNum)
	if err != nil {
		return err
	}

	gp.SetFromGovernanceData(gov)
	WriteGovernanceParam(h.ChainKv, gov.BlockNum, &gp)
	var data StoredGovBlockNums = h.cache.GovBlockNums()
	WriteGovDataBlockNums(h.ChainKv, &data)
	return nil
}

func (h *HeaderGovModule) getExpectedGovernance(epochIdx uint64) GovernanceData {
	votes := h.getExpectedVotes(epochIdx)
	govs := GovernanceData{
		BlockNum: epochIdx * h.epoch,
		Params:   make(map[string]interface{}),
	}

	// TODO: add tally
	for _, vote := range votes {
		govs.Params[vote.Name] = vote.Value
	}

	return govs
}

func (h *HeaderGovModule) getExpectedVotes(epochIdx uint64) []VoteData {
	ret := make([]VoteData, 0)
	for _, vote := range h.cache.Votes {
		if calcEpochIdx(vote.BlockNum, h.epoch) == epochIdx {
			ret = append(ret, vote)
		}
	}

	return ret
}

func (h *HeaderGovModule) VerifyVote(vote *VoteData) error {
	if vote.BlockNum != h.Chain.CurrentBlock().NumberU64()+1 {
		return errors.New("vote block number is not the previous block number")
	}

	gp := GovernanceParam{}
	err := gp.SetFromVoteData(vote)
	if err != nil {
		return err
	}

	/*
		if key == "kip71.lowerboundbasefee" {
			if val.(uint64) > pset.UpperBoundBaseFee() {
				return errInvalidLowerBound
			}
		}
		if key == "kip71.upperboundbasefee" {
			if val.(uint64) < pset.LowerBoundBaseFee() {
				return errInvalidUpperBound
			}
		}
	*/
	return nil
}

func (h *HeaderGovModule) VerifyGov(key string, val interface{}) error {
	_, err := params.NewGovParamSetStrMap(map[string]interface{}{
		key: val,
	})
	if err != nil {
		return err
	}

	return nil
}

func calcEpochIdx(blockNum uint64, epoch uint64) uint64 {
	return blockNum / epoch
}
