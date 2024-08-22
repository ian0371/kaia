package headergov

import (
	"encoding/json"
	"errors"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

func deserializeHeaderVote(b []byte, blockNum uint64) (*VoteData, error) {
	var v struct {
		Validator common.Address
		Key       string
		Value     []byte
	}

	err := rlp.DecodeBytes(b, &v)
	if err != nil {
		return nil, err
	}

	// canonicalize. e.g., [0x1, 0xc9, 0xc3, 0x80] -> 0x1c9c380
	ps, err := params.NewGovParamSetBytesMap(map[string][]byte{
		v.Key: v.Value,
	})
	if err != nil {
		return nil, err
	}

	value, ok := ps.Get(governance.GovernanceKeyMap[v.Key])
	if !ok {
		return nil, errors.New("key not found")
	}

	return &VoteData{
		BlockNum: blockNum,
		Voter:    v.Validator,
		Param:    Param{Name: v.Key, Value: value},
	}, nil
}

func serializeVoteData(vote *VoteData) ([]byte, error) {
	v := &struct {
		Validator common.Address
		Key       string
		Value     interface{}
	}{
		Validator: vote.Voter,
		Key:       vote.Param.Name,
		Value:     vote.Param.Value,
	}

	return rlp.EncodeToBytes(v)
}

func deserializeHeaderGov(b []byte, blockNum uint64) (*GovernanceData, error) {
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

	params := make(map[string]Param, len(j))
	for k := range j {
		value, ok := ps.Get(governance.GovernanceKeyMap[k])
		if !ok {
			return nil, errors.New("key not found")
		}
		params[k] = Param{Name: k, Value: value}
	}

	return &GovernanceData{
		BlockNum: blockNum,
		Params:   params,
	}, nil
}
