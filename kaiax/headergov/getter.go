package headergov

func (h *HeaderGovModule) EffectiveParams(blockNum uint64) (GovernanceParam, error) {
	// TODO: only return when num <= head + 1
	prevEpochStart := PrevEpochStart(blockNum, h.epoch, h.isKoreHF(blockNum))
	gh := h.GetGovernanceHistory()
	gp, err := gh.Search(prevEpochStart)
	if err != nil {
		logger.Error("kaiax.EffectiveParams error", "prevEpochStart", prevEpochStart, "blockNum", blockNum, "err", err,
			"govHistory", gh, "govs", h.cache.Govs())
		return GovernanceParam{}, err
	} else {
		logger.Warn("kaiax.EffectiveParams", "prevEpochStart", prevEpochStart, "blockNum", blockNum, "gp", gp)
		return gp, nil
	}
}

func (h *HeaderGovModule) GetGovernanceHistory() GovernanceHistory {
	return h.cache.GetGovernanceHistory()
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
