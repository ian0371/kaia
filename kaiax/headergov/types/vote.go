package types

import (
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

type VoteData struct {
	Voter common.Address
	Name  string
	Value interface{} // canonicalized value
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

	if cv, ok := vote.Value.(*big.Int); ok {
		v.Value = cv.String()
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

	param, ok := Params[v.Key]
	if !ok {
		return nil, errors.New("invalid param")
	}

	cv, err := param.Canonicalizer(v.Value)
	if err != nil {
		return nil, err
	}

	return &VoteData{
		Voter: v.Validator,
		Name:  v.Key,
		Value: cv,
	}, nil
}
