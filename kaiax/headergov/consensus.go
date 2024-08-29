package headergov

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
)

func (h *HeaderGovModule) VerifyHeader(header *types.Header) error {
	if header.Number.Uint64() == 0 {
		return nil
	}

	// 1. Check Vote
	if len(header.Vote) > 0 {
		vote, err := headergov_types.DeserializeHeaderVote(header.Vote, header.Number.Uint64())
		if err != nil {
			logger.Error("Failed to parse vote", "num", header.Number.Uint64(), "err", err)
			return err
		}
		err = h.VerifyVote(vote)
		if err != nil {
			logger.Error("Failed to verify vote", "num", header.Number.Uint64(), "err", err)
			return err
		}
	}

	// 2. Check Governance
	// TODO-kaiax: fail fast
	if header.Number.Uint64()%h.epoch != 0 {
		if len(header.Governance) > 0 {
			logger.Error("governance is not allowed in non-epoch block", "num", header.Number.Uint64())
			return errors.New("governance is not allowed in non-epoch block")
		} else {
			return nil
		}
	} else {
		expected := h.getExpectedGovernance(header.Number.Uint64())
		if len(expected.Params) == 0 && len(header.Governance) == 0 {
			return nil
		}

		actual, err := headergov_types.DeserializeHeaderGov(header.Governance, header.Number.Uint64())
		if err != nil {
			logger.Error("Failed to parse governance", "num", header.Number.Uint64(), "err", err)
			return err
		}
		if !reflect.DeepEqual(&expected, actual) {
			logger.Error("governance mismatch", "num", header.Number.Uint64(), "expected", &expected, "actual", actual)
			return fmt.Errorf("expected governance: %v, actual: %v", &expected, actual)
		}

		return nil
	}
}

func (h *HeaderGovModule) PrepareHeader(header *types.Header) (*types.Header, error) {
	// if epoch block & vote exists in the last epoch, put Governance to header.
	if len(h.MyVotes) > 0 {
		header.Vote, _ = h.MyVotes[0].Serialize()
	}

	if header.Number.Uint64()%h.epoch == 0 {
		gov := h.getExpectedGovernance(header.Number.Uint64())
		if len(gov.Params) > 0 {
			header.Governance, _ = gov.Serialize()
		}
	}

	return header, nil
}

func (h *HeaderGovModule) FinalizeBlock(b *types.Block) (*types.Block, error) {
	// TODO-kaiax: must be removed later. only for testing.
	h.PostInsertBlock(b)
	return b, nil
}

func (h *HeaderGovModule) VerifyVote(vote *VoteData) error {
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

func (h *HeaderGovModule) VerifyGov(gov *GovernanceData) error {
	gp := GovernanceParam{}
	err := gp.SetFromGovernanceData(gov)
	if err != nil {
		return err
	}

	return nil
}

func (h *HeaderGovModule) getExpectedGovernance(blockNum uint64) GovernanceData {
	prevEpochVotes := h.getVotesInEpoch(calcEpochIdx(blockNum, h.epoch) - 1)
	govs := GovernanceData{
		Params: make(map[string]interface{}),
	}

	// TODO: add tally
	for _, vote := range prevEpochVotes {
		govs.Params[vote.Name] = vote.Value
	}

	return govs
}

func (h *HeaderGovModule) getVotesInEpoch(epochIdx uint64) map[uint64]VoteData {
	votes := make(map[uint64]VoteData)
	for blockNum, vote := range h.cache.GroupedVotes()[epochIdx] {
		votes[blockNum] = vote
	}
	return votes
}
