package headergov

import (
	gov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/params"
)

func (h *HeaderGovModule) EffectiveParams(num uint64) (*params.GovParamSet, error) {
	// TODO: only return when num <= head + 1
	allParamsHistory := gov_types.GetAllParamsHistory(h.cache.Govs, h.epoch, h.ChainConfig.KoreCompatibleBlock.Uint64())
	ret := make(map[string]interface{})
	logger.Debug("EffectiveParams", "num", num, "allParamHistory", allParamsHistory)
	for _, paramHistory := range allParamsHistory {
		param := paramHistory.GetItem(uint(num))
		ret[param.Name] = param.Value
	}
	return params.NewGovParamSetStrMap(ret)
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
