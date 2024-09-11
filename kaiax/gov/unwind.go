package gov

func (m *govModule) Unwind(blockNum uint64) error {
	return m.hgm.Unwind(blockNum)
}
