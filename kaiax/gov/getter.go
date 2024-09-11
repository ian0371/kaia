package gov

func (m *govModule) EffectiveParamSet(blockNum uint64) (ParamSet, error) {
	p1 := m.hgm.EffectiveParamsPartial(blockNum)
	ret := ParamSet{}
	for k, v := range p1 {
		ret.Set(k, v)
	}

	p2 := m.cgm.EffectiveParamsPartial(blockNum)
	for k, v := range p2 {
		ret.Set(k, v)
	}
	return ret, nil
}
