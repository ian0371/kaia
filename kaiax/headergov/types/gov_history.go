package types

import "errors"

type GovernanceHistory map[uint64]GovernanceParam

func GetGovernanceHistory(govList []GovernanceData) GovernanceHistory {
	gp := GovernanceParam{}
	gh := make(map[uint64]GovernanceParam)
	for _, gov := range govList {
		gp.SetFromGovernanceData(&gov)
		gh[gov.BlockNum] = gp
	}
	return gh
}

func (g *GovernanceHistory) Search(blockNum uint64) (GovernanceParam, error) {
	idx := uint64(0)
	for num := range *g {
		if num > idx && num <= blockNum {
			idx = num
		}
	}
	if ret, ok := (*g)[idx]; ok {
		return ret, nil
	} else {
		return GovernanceParam{}, errors.New("blockNum not found from governance history")
	}
}
