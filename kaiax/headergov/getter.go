package headergov

import (
	"github.com/kaiachain/kaia/params"
)

func (h *HeaderGovModule) EffectiveParams(num uint64) (*params.GovParamSet, error) {
	// TODO: only return when num <= head + 1
	sourceBlock := PrevEpochStart(num, h.epoch, h.isKoreHF(num))
	return h.cache.GovMap.GetItem(uint(sourceBlock)), nil
}

func PrevEpochStart(num, epoch uint64, isKore bool) uint64 {
	if num <= epoch {
		return 0
	}
	if !isKore {
		num -= 1
	}
	return num - num%epoch - epoch
}
