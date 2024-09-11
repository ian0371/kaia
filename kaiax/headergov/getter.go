package headergov

func (h *headerGovModule) EffectiveParamSet(blockNum uint64) (ParamSet, error) {
	// TODO: only return when num <= head + 1
	prevEpochStart := PrevEpochStart(blockNum, h.epoch, h.isKoreHF(blockNum))
	gh := h.GetGovernanceHistory()
	gp, err := gh.Search(prevEpochStart)
	if err != nil {
		logger.Error("kaiax.EffectiveParams error", "prevEpochStart", prevEpochStart, "blockNum", blockNum, "err", err,
			"govHistory", gh, "govs", h.cache.Govs())
		return ParamSet{}, err
	} else {
		return gp, nil
	}
}

func (h *headerGovModule) EffectiveParamsPartial(blockNum uint64) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	for num, gov := range h.cache.Govs() {
		if num > blockNum {
			continue
		}
		for name, value := range gov.Items() {
			ret[name] = value
		}
	}

	return ret, nil
}

func (h *headerGovModule) GetGovernanceHistory() History {
	return h.cache.History()
}

func PrevEpochStart(blockNum, epoch uint64, isKore bool) uint64 {
	if blockNum <= epoch {
		return 0
	}
	if !isKore {
		blockNum -= 1
	}
	return blockNum - blockNum%epoch - epoch
}
