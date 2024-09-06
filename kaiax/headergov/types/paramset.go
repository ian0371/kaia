package types

import (
	"encoding/json"
	"errors"
	"math/big"
	"reflect"

	"github.com/kaiachain/kaia/common"
)

type ParamSet struct {
	// governance
	GovernanceMode                  string
	GoverningNode, GovParamContract common.Address

	// istanbul
	CommitteeSize, ProposerPolicy, Epoch uint64

	// reward
	Ratio, Kip82Ratio                             string
	StakingUpdateInterval, ProposerUpdateInterval uint64
	MintingAmount, MinimumStake                   *big.Int
	UseGiniCoeff, DeferredTxFee                   bool

	// KIP-71
	LowerBoundBaseFee, UpperBoundBaseFee, GasTarget, MaxBlockGasUsedForBaseFee, BaseFeeDenominator uint64

	// etc.
	DeriveShaImpl uint64
	UnitPrice     uint64
}

// TODO: add tests, compare from gov/default
func (p *ParamSet) Set(key string, cv interface{}) error {
	param, ok := Params[key]
	if !ok {
		return errors.New("invalid param key")
	}

	field := reflect.ValueOf(p).Elem().FieldByName(param.ParamSetFieldName)
	if !field.IsValid() || !field.CanSet() {
		return errors.New("invalid field or cannot set value")
	}

	fieldValue := reflect.ValueOf(cv)
	if !fieldValue.Type().AssignableTo(field.Type()) {
		return errors.New("type mismatch")
	}

	field.Set(fieldValue)
	return nil
}

func (p *ParamSet) SetFromVoteData(v *voteData) error {
	return p.Set(v.name, v.value)
}

func (p *ParamSet) SetFromGovernanceData(g *GovData) error {
	for name, value := range g.Params {
		err := p.Set(name, value)
		if err != nil {
			continue
		}
	}
	return nil
}

func (p *ParamSet) ToJSON() (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (p *ParamSet) ToStrMap() (map[string]interface{}, error) {
	jsonStr, err := p.ToJSON()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *ParamSet) Copy() *ParamSet {
	return &ParamSet{
		GovernanceMode:            p.GovernanceMode,
		GoverningNode:             p.GoverningNode,
		GovParamContract:          p.GovParamContract,
		CommitteeSize:             p.CommitteeSize,
		ProposerPolicy:            p.ProposerPolicy,
		Epoch:                     p.Epoch,
		Ratio:                     p.Ratio,
		Kip82Ratio:                p.Kip82Ratio,
		StakingUpdateInterval:     p.StakingUpdateInterval,
		ProposerUpdateInterval:    p.ProposerUpdateInterval,
		MintingAmount:             new(big.Int).Set(p.MintingAmount),
		MinimumStake:              new(big.Int).Set(p.MinimumStake),
		UseGiniCoeff:              p.UseGiniCoeff,
		DeferredTxFee:             p.DeferredTxFee,
		LowerBoundBaseFee:         p.LowerBoundBaseFee,
		UpperBoundBaseFee:         p.UpperBoundBaseFee,
		GasTarget:                 p.GasTarget,
		MaxBlockGasUsedForBaseFee: p.MaxBlockGasUsedForBaseFee,
		BaseFeeDenominator:        p.BaseFeeDenominator,
		DeriveShaImpl:             p.DeriveShaImpl,
		UnitPrice:                 p.UnitPrice,
	}
}
