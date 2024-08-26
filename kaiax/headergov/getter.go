package headergov

func (h *HeaderGovModule) EffectiveParams(num uint64) (GovernanceParam, error) {
	// TODO: only return when num <= head + 1
	sourceBlock := PrevEpochStart(num, h.epoch, h.isKoreHF(num))
	gh := h.GetGovernanceHistory()
	gp, err := gh.Search(sourceBlock)
	if err != nil {
		return GovernanceParam{}, err
	} else {
		return gp, nil
	}
}

func (h *HeaderGovModule) GetGovernanceHistory() GovernanceHistory {
	return h.cache.GetGovernanceHistory()
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
