package types

import (
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/params"
)

type HeaderGovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.ExecutionModule
	kaiax.UnwindableModule

	EffectiveParams(num uint64) (*params.GovParamSet, error)
}
