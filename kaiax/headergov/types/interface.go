package types

import (
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=kaiax/headergov/mocks/headergov_mock.go github.com/kaiachain/kaia/kaiax/headergov/types HeaderGovModule
type HeaderGovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.ExecutionModule
	kaiax.UnwindableModule

	EffectiveParamSet(blockNum uint64) (ParamSet, error)
	EffectiveParamsPartial(blockNum uint64) (map[string]interface{}, error)
}
