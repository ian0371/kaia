package types

import (
	"github.com/kaiachain/kaia/kaiax"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
)

type GovModule interface {
	kaiax.BaseModule
	// kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.ExecutionModule
	kaiax.UnwindableModule

	EffectiveParamSet(blockNum uint64) (headergov_types.ParamSet, error)
}
