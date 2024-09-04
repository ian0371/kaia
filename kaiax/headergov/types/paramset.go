package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/common"
)

type GovParamSet struct {
	// governance
	GovernanceMode                  string
	GoverningNode, GovParamContract common.Address

	// istanbul
	CommitteeSize, ProposerPolicy, Epoch uint64

	// reward
	Ratio, Kip82Ratio                            string
	StakeUpdateInterval, ProposerRefreshInterval uint64
	MintingAmount, MinimumStake                  *big.Int
	UseGiniCoeff, DeferredTxFee                  bool

	// KIP-71
	LowerBoundBaseFee, UpperBoundBaseFee, GasTarget, MaxBlockGasUsedForBaseFee, BaseFeeDenominator uint64

	// etc.
	DeriveShaImpl uint64
	UnitPrice     uint64
}

// TODO: add tests, compare from gov/default
func (p *GovParamSet) Set(key string, value interface{}) error {
	switch key {
	case "governance.governancemode":
		switch v := value.(type) {
		case string:
			if v == "none" || v == "single" || v == "ballot" {
				p.GovernanceMode = v
			} else {
				return errors.New("invalid governance mode")
			}
		default:
			return errors.New("invalid governance mode")
		}
	case "governance.governingnode":
		switch v := value.(type) {
		case common.Address:
			p.GoverningNode = v
		default:
			return errors.New("invalid governing node")
		}
	case "governance.govparamcontract":
		switch v := value.(type) {
		case common.Address:
			p.GovParamContract = v
		default:
			return errors.New("invalid governance param contract")
		}
	case "istanbul.epoch":
		switch v := value.(type) {
		case uint64:
			p.Epoch = v
		default:
			return errors.New("invalid epoch")
		}
	case "istanbul.policy":
		switch v := value.(type) {
		case uint64:
			p.ProposerPolicy = v
		default:
			return errors.New("invalid proposer policy")
		}
	case "istanbul.committeesize":
		switch v := value.(type) {
		case uint64:
			p.CommitteeSize = v
		default:
			return errors.New("invalid committee size")
		}
	case "governance.unitprice":
		switch v := value.(type) {
		case uint64:
			p.UnitPrice = v
		default:
			return errors.New("invalid unit price")
		}
	case "governance.deriveshaimpl":
		switch v := value.(type) {
		case uint64:
			p.DeriveShaImpl = v
		default:
			return errors.New("invalid derive sha impl")
		}
	case "kip71.lowerboundbasefee":
		switch v := value.(type) {
		case uint64:
			p.LowerBoundBaseFee = v
		default:
			return errors.New("invalid lower bound base fee")
		}
	case "kip71.gastarget":
		switch v := value.(type) {
		case uint64:
			p.GasTarget = v
		default:
			return errors.New("invalid gas target")
		}
	case "kip71.maxblockgasusedforbasefee":
		switch v := value.(type) {
		case uint64:
			p.MaxBlockGasUsedForBaseFee = v
		default:
			return errors.New("invalid max block gas used for base fee")
		}
	case "kip71.basefeedenominator":
		switch v := value.(type) {
		case uint64:
			p.BaseFeeDenominator = v
		default:
			return errors.New("invalid base fee denominator")
		}
	case "kip71.upperboundbasefee":
		switch v := value.(type) {
		case uint64:
			p.UpperBoundBaseFee = v
		default:
			return errors.New("invalid upper bound base fee")
		}
	case "reward.mintingamount":
		switch v := value.(type) {
		case *big.Int:
			p.MintingAmount = v
		case uint64:
			p.MintingAmount = big.NewInt(int64(v))
		case float64:
			p.MintingAmount = big.NewInt(int64(v))
		case string:
			var ok bool
			p.MintingAmount, ok = new(big.Int).SetString(v, 10)
			if !ok {
				return errors.New("invalid minting amount")
			}
		default:
			return errors.New("invalid minting amount")
		}
	case "reward.ratio":
		switch v := value.(type) {
		case string:
			p.Ratio = v
		default:
			return errors.New("invalid ratio")
		}
	case "reward.kip82ratio":
		switch v := value.(type) {
		case string:
			p.Kip82Ratio = v
		default:
			return errors.New("invalid kip82 ratio")
		}
	case "reward.useginicoeff":
		switch v := value.(type) {
		case bool:
			p.UseGiniCoeff = v
		default:
			return errors.New("invalid use gini coeff")
		}
	case "reward.deferredtxfee":
		switch v := value.(type) {
		case bool:
			p.DeferredTxFee = v
		default:
			return errors.New("invalid deferred tx fee")
		}
	case "reward.minimumstake":
		switch v := value.(type) {
		case *big.Int:
			p.MinimumStake = v
		case uint64:
			p.MinimumStake = big.NewInt(int64(v))
		case float64:
			p.MinimumStake = big.NewInt(int64(v))
		case string:
			var ok bool
			p.MinimumStake, ok = new(big.Int).SetString(v, 10)
			if !ok {
				return errors.New("invalid minimum stake")
			}
		default:
			return errors.New("invalid minimum stake")
		}
	case "reward.stakingupdateinterval":
		switch v := value.(type) {
		case uint64:
			p.StakeUpdateInterval = v
		default:
			return errors.New("invalid stake update interval")
		}
	case "reward.proposerupdateinterval":
		switch v := value.(type) {
		case uint64:
			p.ProposerRefreshInterval = v
		default:
			return errors.New("invalid proposer refresh interval")
		}
	default:
		return errors.New("unknown parameter name")
	}

	return nil
}

func (p *GovParamSet) SetFromVoteData(v *VoteData) error {
	return p.Set(v.Name, v.Value)
}

func (p *GovParamSet) SetFromGovernanceData(g *GovData) error {
	for name, value := range g.Params {
		err := p.Set(name, value)
		if err != nil {
			continue
		}
	}
	return nil
}

func (p *GovParamSet) ToJSON() (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (p *GovParamSet) ToStrMap() (map[string]interface{}, error) {
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
