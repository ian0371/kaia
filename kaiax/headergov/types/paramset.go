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

func GetDefaultGovernanceParamSet() *ParamSet {
	p := &ParamSet{}
	for name, param := range Params {
		p.Set(name, param.DefaultValue)
	}

	return p
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

func (p *ParamSet) SetFromGovernanceData(g GovData) error {
	for name, value := range g.Items() {
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

// ToStrMap is used for test and getParams API.
func (p *ParamSet) ToStrMap() (map[string]interface{}, error) {
	ret := make(map[string]interface{})

	// Iterate through all params in Params and ensure they're in the result
	for paramName, param := range Params {
		field := reflect.ValueOf(p).Elem().FieldByName(param.ParamSetFieldName)
		if field.IsValid() {
			// Convert big.Int to string for JSON compatibility
			if bigIntValue, ok := field.Interface().(*big.Int); ok {
				ret[paramName] = bigIntValue.String()
			} else {
				ret[paramName] = field.Interface()
			}
		}
	}

	return ret, nil
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
