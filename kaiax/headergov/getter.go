package headergov

import (
	"github.com/kaiachain/kaia/params"
)

func (h *HeaderGovModule) EffectiveParams(num uint64) (*params.GovParamSet, error) {
	// TODO: only return when num <= head + 1
	govBlock := CalcGovDataBlock(num, h.epoch, h.isKoreHF(num))
	return h.cache.GovMap.GetItem(uint(govBlock)), nil
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
