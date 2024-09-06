package types

import (
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

// Always use NewVoteData (constructor) to create a new VoteData.
type VoteData struct {
	Voter common.Address
	Name  string
	Value interface{} // canonicalized value
}

func NewVoteData(voter common.Address, name string, value interface{}) *VoteData {
	v := &VoteData{Voter: voter, Name: name, Value: value}
	err := v.Canonicalize()
	if err != nil {
		return nil
	}
	return v
}

func (vote *VoteData) ToParamSet() (*params.GovParamSet, error) {
	tmp := map[string]interface{}{
		vote.Name: vote.Value,
	}

	return params.NewGovParamSetStrMap(tmp)
}

func (vote *VoteData) Canonicalize() error {
	param, ok := Params[vote.Name]
	if !ok {
		return errors.New("invalid param key")
	}

	cv, err := param.Canonicalizer(vote.Value)
	if err != nil {
		return err
	}

	vote.Value = cv
	return nil
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

	vote := NewVoteData(v.Validator, v.Key, v.Value)
	if vote == nil {
		return nil, errors.New("failed to create vote data")
	}

	return vote, nil
}
