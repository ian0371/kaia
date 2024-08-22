package headergov

import (
	gov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/params"
)

func (h *HeaderGovModule) EffectiveParams(num uint64) (*params.GovParamSet, error) {
	allParamsHistory := gov_types.GetAllParamsHistory(h.cache.Govs, h.epoch, h.ChainConfig.KoreCompatibleBlock.Uint64())
	ret := make(map[string]interface{})
	for _, paramHistory := range allParamsHistory {
		param := paramHistory.GetItem(uint(num))
		logger.Debug("EffectiveParams",
			"num", num,
			"name", param.Name,
			"value", param.Value,
			"paramHistory", paramHistory,
		)
		ret[param.Name] = param.Value
	}
	return params.NewGovParamSetStrMap(ret)
}
