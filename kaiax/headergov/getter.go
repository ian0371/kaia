package headergov

import (
	gov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/params"
)

func (h *HeaderGovModule) EffectiveParams(num uint64) (*params.GovParamSet, error) {
	// TODO: only return when num <= head + 1
	allParamsHistory, err := gov_types.GetAllParamsHistory(h.cache.Govs)
	if err != nil {
		return nil, err
	}

	govBlock := CalcGovDataBlock(num, h.epoch, h.isKoreHF(num))
	return allParamsHistory.GetItem(uint(govBlock)), nil
}

func CalcGovDataBlock(num uint64, epoch uint64, isKore bool) uint64 {
	if num <= epoch {
		return 0
	}
	if isKore {
		if num%epoch == 0 {
			return num - epoch
		} else {
			return num - num%epoch - epoch
		}
	} else {
		if num%epoch == 0 {
			return num - 2*epoch
		} else {
			return num - num%epoch - epoch
		}
	}
}
