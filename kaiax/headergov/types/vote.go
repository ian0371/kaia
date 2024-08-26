package types

import (
	"errors"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

type VoteData struct {
	BlockNum uint64
	Voter    common.Address
	Name     string
	Value    interface{}
}

func (vote *VoteData) ToParamSet() (*params.GovParamSet, error) {
	tmp := map[string]interface{}{
		vote.Name: vote.Value,
	}

	return params.NewGovParamSetStrMap(tmp)
}

func (vote *VoteData) Serialize() ([]byte, error) {
	v := &struct {
		Validator common.Address
		Key       string
		Value     interface{}
	}{
		Validator: vote.Voter,
		Key:       vote.Name,
		Value:     vote.Value,
	}

	return rlp.EncodeToBytes(v)
}

func DeserializeHeaderVote(b []byte, blockNum uint64) (*VoteData, error) {
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
		Name:     v.Key,
		Value:    value,
	}, nil
}
