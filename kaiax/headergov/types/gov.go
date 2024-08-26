package types

import (
	"encoding/json"
	"errors"

	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

type GovernanceData struct {
	BlockNum uint64
	Params   map[string]interface{} // canonicalized value
}

func (g *GovernanceData) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for name, value := range g.Params {
		tmp[name] = value
	}

	return json.Marshal(tmp)
}

func (g *GovernanceData) Serialize() ([]byte, error) {
	j, err := g.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(j)
}

func DeserializeHeaderGov(b []byte, blockNum uint64) (*GovernanceData, error) {
	rlpDecoded := []byte("")
	err := rlp.DecodeBytes(b, &rlpDecoded)
	if err != nil {
		return nil, err
	}

	j := make(map[string]interface{})
	json.Unmarshal(rlpDecoded, &j)

	ps, err := params.NewGovParamSetStrMap(j)
	if err != nil {
		return nil, err
	}

	params := make(map[string]interface{}, len(j))
	for k := range j {
		value, ok := ps.Get(governance.GovernanceKeyMap[k])
		if !ok {
			return nil, errors.New("key not found")
		}
		params[k] = value
	}

	return &GovernanceData{
		BlockNum: blockNum,
		Params:   params,
	}, nil
}
