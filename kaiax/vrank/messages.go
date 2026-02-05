package vrank

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

type VRankPreprepare struct {
	Block *types.Block
}

type VRankCandidate struct {
	BlockNumber uint64
	Round       uint8
	BlockHash   common.Hash
	Sig         []byte
}
