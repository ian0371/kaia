package types

import (
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/params"
)

type ContractGovModule interface {
	kaiax.BaseModule

	EffectiveParams(num uint64) (*params.GovParamSet, error)
}
