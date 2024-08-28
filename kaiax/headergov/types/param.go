package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/common"
)

type GovernanceParam struct {
	// governance
	GovernanceMode                  int
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

const (
	GovernanceMode_None = iota
	GovernanceMode_Single
	GovernanceMode_Ballot
	GovernanceMode_End
)

const (
	// Proposer policy
	RoundRobin = iota
	Sticky
	WeightedRandom
	ProposerPolicy_End
)

var (
	// Default Values: Constants used for getting default values for configuration
	defaultGovernanceMode            = GovernanceMode_None
	defaultGoverningNode             = common.HexToAddress("0x0000000000000000000000000000000000000000")
	defaultGovParamContract          = common.HexToAddress("0x0000000000000000000000000000000000000000")
	defaultEpoch                     = uint64(604800)
	defaultProposerPolicy            = uint64(RoundRobin)
	defaultCommitteeSize             = uint64(21)
	defaultUnitPrice                 = uint64(250000000000)
	defaultLowerBoundBaseFee         = uint64(25000000000)
	defaultUpperBoundBaseFee         = uint64(750000000000)
	defaultGasTarget                 = uint64(30000000)
	defaultMaxBlockGasUsedForBaseFee = uint64(60000000)
	defaultBaseFeeDenominator        = uint64(20)
	defaultMintingAmount             = big.NewInt(0)
	defaultRatio                     = "100/0/0"
	defaultKip82Ratio                = "20/80"
	defaultUseGiniCoeff              = false
	defaultDeferredTxFee             = false
	defaultMinimumStake              = big.NewInt(2000000)
	defaultStakeUpdateInterval       = uint64(86400) // 1 day
	defaultProposerRefreshInterval   = uint64(3600)  // 1 hour
	defaultDeriveShaImpl             = uint64(0)     // Orig
)

func GetDefaultGovernanceParam() *GovernanceParam {
	return &GovernanceParam{
		GovernanceMode:            defaultGovernanceMode,
		GoverningNode:             defaultGoverningNode,
		GovParamContract:          defaultGovParamContract,
		CommitteeSize:             defaultCommitteeSize,
		ProposerPolicy:            defaultProposerPolicy,
		Epoch:                     defaultEpoch,
		Ratio:                     defaultRatio,
		Kip82Ratio:                defaultKip82Ratio,
		StakeUpdateInterval:       defaultStakeUpdateInterval,
		ProposerRefreshInterval:   defaultProposerRefreshInterval,
		MintingAmount:             defaultMintingAmount,
		MinimumStake:              defaultMinimumStake,
		UseGiniCoeff:              defaultUseGiniCoeff,
		DeferredTxFee:             defaultDeferredTxFee,
		LowerBoundBaseFee:         defaultLowerBoundBaseFee,
		UpperBoundBaseFee:         defaultUpperBoundBaseFee,
		GasTarget:                 defaultGasTarget,
		MaxBlockGasUsedForBaseFee: defaultMaxBlockGasUsedForBaseFee,
		BaseFeeDenominator:        defaultBaseFeeDenominator,
		DeriveShaImpl:             defaultDeriveShaImpl,
		UnitPrice:                 defaultUnitPrice,
	}
}

func (p *GovernanceParam) Set(key string, value interface{}) error {
	switch key {
	case "governance.governancemode":
		switch v := value.(type) {
		case int:
			if v < 0 || v >= GovernanceMode_End {
				return errors.New("invalid governance mode")
			}
			p.GovernanceMode = v
		case string:
			switch v {
			case "none":
				p.GovernanceMode = GovernanceMode_None
			case "single":
				p.GovernanceMode = GovernanceMode_Single
			case "ballot":
				p.GovernanceMode = GovernanceMode_Ballot
			default:
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

func (p *GovernanceParam) SetFromVoteData(v *VoteData) error {
	return p.Set(v.Name, v.Value)
}

func (p *GovernanceParam) SetFromGovernanceData(g *GovernanceData) error {
	for name, value := range g.Params {
		err := p.Set(name, value)
		if err != nil {
			logger.Warn("kaiax.SetFromGovernanceData error", "name", name, "value", value, "err", err)
			continue
		}
	}
	return nil
}

func (p *GovernanceParam) ToJSON() (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (p *GovernanceParam) ToStrMap() (map[string]interface{}, error) {
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
