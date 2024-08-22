package headergov

import "github.com/kaiachain/kaia/blockchain/types"

func (h *HeaderGovModule) VerifyHeader(*types.Header) error {
	return nil // TODO: implement
}

func (h *HeaderGovModule) PrepareHeader(*types.Header) (*types.Header, error) {
	return nil, nil // TODO: implement
}

func (h *HeaderGovModule) FinalizeBlock() (*types.Block, error) {
	return nil, nil
}

func (h *HeaderGovModule) PostInsertBlock(b *types.Block) error {
	vote, err := deserializeHeaderVote(b.Header().Vote, b.NumberU64())
	if err != nil {
		return err
	}
	h.AddVote(vote)

	gov, err := deserializeHeaderGov(b.Header().Vote, b.NumberU64())
	if err != nil {
		return err
	}
	h.AddGov(gov)

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
