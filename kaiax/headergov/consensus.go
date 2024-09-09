package headergov

import (
	"errors"
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types"
)

func (h *HeaderGovModule) VerifyHeader(header *types.Header) error {
	if header.Number.Uint64() == 0 {
		return nil
	}

	// 1. Check Vote
	if len(header.Vote) > 0 {
		vote, err := DeserializeHeaderVote(header.Vote, header.Number.Uint64())
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
	if header.Number.Uint64()%h.epoch != 0 {
		if len(header.Governance) > 0 {
			logger.Error("governance is not allowed in non-epoch block", "num", header.Number.Uint64())
			return errors.New("governance is not allowed in non-epoch block")
		} else {
			return nil
		}
	}

	gov, err := DeserializeHeaderGov(header.Governance, header.Number.Uint64())
	if err != nil {
		logger.Error("Failed to parse governance", "num", header.Number.Uint64(), "err", err)
		return err
	}
	return h.VerifyGov(header.Number.Uint64(), gov)
}

func (h *HeaderGovModule) PrepareHeader(header *types.Header) (*types.Header, error) {
	// if epoch block & vote exists in the last epoch, put Governance to header.
	if len(h.myVotes) > 0 {
		header.Vote, _ = h.myVotes[0].Serialize()
	}

	if header.Number.Uint64()%h.epoch == 0 {
		gov := h.getExpectedGovernance(header.Number.Uint64())
		if len(gov.Items()) > 0 {
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

// VerifyVote takes canonical VoteData.
func (h *HeaderGovModule) VerifyVote(vote VoteData) error {
	if vote == nil {
		return errors.New("vote is nil")
	}

	// handled by valset module.
	if vote.Name() == "governance.addvalidator" || vote.Name() == "governance.removevalidator" {
		return nil
	}

	param, ok := Params[vote.Name()]
	if !ok {
		return errors.New("invalid param key")
	}

	if param.VoteForbidden {
		return errors.New("parameter is forbidden to be changed")
	}

	if param.FormatChecker != nil {
		if valid := param.FormatChecker(vote.Value); !valid {
			return errors.New("invalid format")
		}
	}

	return nil
}

func (h *HeaderGovModule) VerifyGov(blockNum uint64, gov GovData) error {
	expected := h.getExpectedGovernance(blockNum)
	if !reflect.DeepEqual(expected, gov) {
		return errors.New("governance is not matched")
	}

	return nil
}

// blockNum must be greater than epoch.
func (h *HeaderGovModule) getExpectedGovernance(blockNum uint64) GovData {
	prevEpochIdx := calcEpochIdx(blockNum, h.epoch) - 1
	prevEpochVotes := h.getVotesInEpoch(prevEpochIdx)
	govs := make(map[string]interface{})

	// TODO: add tally
	for _, vote := range prevEpochVotes {
		govs[vote.Name()] = vote.Value()
	}

	return NewGovData(govs)
}

func (h *HeaderGovModule) getVotesInEpoch(epochIdx uint64) map[uint64]VoteData {
	votes := make(map[uint64]VoteData)
	for blockNum, vote := range h.cache.GroupedVotes()[epochIdx] {
		votes[blockNum] = vote
	}
	return votes
}
