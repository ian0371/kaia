package types

import (
	"github.com/kaiachain/kaia/kaiax"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
)

//go:generate mockgen -destination=kaiax/contractgov/mocks/contractgov_mock.go github.com/kaiachain/kaia/kaiax/contractgov/types ContractGovModule
type ContractGovModule interface {
	kaiax.BaseModule

	EffectiveParamSet(blockNum uint64) (headergov_types.ParamSet, error)
	EffectiveParamsPartial(blockNum uint64) (map[string]interface{}, error)
}
