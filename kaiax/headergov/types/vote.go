package types

import (
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
)

type VoteData interface {
	Voter() common.Address
	Name() string
	Value() interface{}

	Serialize() ([]byte, error)
}

var _ VoteData = (*voteData)(nil)

type voteData struct {
	voter common.Address
	name  string
	value interface{} // canonicalized value
}

// NewVoteData returns a canonical & formatted vote data. Consistency is NOT checked.
func NewVoteData(voter common.Address, name string, value interface{}) VoteData {
	v := &voteData{voter: voter, name: name, value: value}
	param, ok := Params[v.name]
	if !ok {
		if name == "governance.addvalidator" || name == "governance.removevalidator" {
			v.value = []common.Address{} // don't care about the value
			return v
		} else {
			return nil
		}
	}

	if param.VoteForbidden {
		return nil
	}

	cv, err := param.Canonicalizer(v.value)
	if err != nil {
		return nil
	}

	if param.FormatChecker != nil && !param.FormatChecker(cv) {
		return nil
	}

	v.value = cv
	return v
}

func (vote *voteData) Voter() common.Address {
	return vote.voter
}

func (vote *voteData) Name() string {
	return vote.name
}

func (vote *voteData) Value() interface{} {
	return vote.value
}

func (vote *voteData) Serialize() ([]byte, error) {
	v := &struct {
		Validator common.Address
		Key       string
		Value     interface{}
	}{
		Validator: vote.voter,
		Key:       vote.name,
		Value:     vote.value,
	}

	if cv, ok := vote.value.(*big.Int); ok {
		v.Value = cv.String()
	}

	return rlp.EncodeToBytes(v)
}

func DeserializeHeaderVote(b []byte, blockNum uint64) (VoteData, error) {
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
