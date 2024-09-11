package gov

import "github.com/kaiachain/kaia/blockchain/types"

func (g *govModule) VerifyHeader(header *types.Header) error {
	return g.hgm.VerifyHeader(header)
}

func (g *govModule) PrepareHeader(header *types.Header) (*types.Header, error) {
	return g.hgm.PrepareHeader(header)
}

func (g *govModule) FinalizeBlock(b *types.Block) (*types.Block, error) {
	return g.hgm.FinalizeBlock(b)
}
