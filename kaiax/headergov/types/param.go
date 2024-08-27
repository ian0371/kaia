package types

import (
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
	defaultPeriod                    = uint64(1)
	defaultDeriveShaImpl             = uint64(0) // Orig
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
		if val, ok := value.(int); !ok {
			return errors.New("invalid governance mode")
		} else {
			if val < 0 || val >= GovernanceMode_End {
				return errors.New("invalid governance mode")
			}
			p.GovernanceMode = val
		}
	case "governance.governingnode":
		if val, ok := value.(common.Address); !ok {
			return errors.New("invalid governing node")
		} else {
			p.GoverningNode = val
		}
	case "governance.govparamcontract":
		if val, ok := value.(common.Address); !ok {
			return errors.New("invalid governance param contract")
		} else {
			p.GovParamContract = val
		}
	case "istanbul.epoch":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid epoch")
		} else {
			p.Epoch = val
		}
	case "istanbul.policy":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid proposer policy")
		} else {
			p.ProposerPolicy = val
		}
	case "istanbul.committeesize":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid committee size")
		} else {
			p.CommitteeSize = val
		}
	case "governance.unitprice":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid unit price")
		} else {
			p.UnitPrice = val
		}
	case "governance.deriveshaimpl":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid derive sha impl")
		} else {
			p.DeriveShaImpl = val
		}
	case "kip71.lowerboundbasefee":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid lower bound base fee")
		} else {
			p.LowerBoundBaseFee = val
		}
	case "kip71.gastarget":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid gas target")
		} else {
			p.GasTarget = val
		}
	case "kip71.maxblockgasusedforbasefee":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid max block gas used for base fee")
		} else {
			p.MaxBlockGasUsedForBaseFee = val
		}
	case "kip71.basefeedenominator":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid base fee denominator")
		} else {
			p.BaseFeeDenominator = val
		}
	case "kip71.upperboundbasefee":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid upper bound base fee")
		} else {
			p.UpperBoundBaseFee = val
		}
	case "reward.mintingamount":
		if val, ok := value.(*big.Int); !ok {
			return errors.New("invalid minting amount")
		} else {
			p.MintingAmount = val
		}
	case "reward.ratio":
		if val, ok := value.(string); !ok {
			return errors.New("invalid ratio")
		} else {
			p.Ratio = val
		}
	case "reward.kip82ratio":
		if val, ok := value.(string); !ok {
			return errors.New("invalid kip82 ratio")
		} else {
			p.Kip82Ratio = val
		}
	case "reward.useginicoeff":
		if val, ok := value.(bool); !ok {
			return errors.New("invalid use gini coeff")
		} else {
			p.UseGiniCoeff = val
		}
	case "reward.deferredtxfee":
		if val, ok := value.(bool); !ok {
			return errors.New("invalid deferred tx fee")
		} else {
			p.DeferredTxFee = val
		}
	case "reward.minimumstake":
		if val, ok := value.(*big.Int); !ok {
			return errors.New("invalid minimum stake")
		} else {
			p.MinimumStake = val
		}
	case "reward.stakingupdateinterval":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid stake update interval")
		} else {
			p.StakeUpdateInterval = val
		}
	case "reward.proposerupdateinterval":
		if val, ok := value.(uint64); !ok {
			return errors.New("invalid proposer refresh interval")
		} else {
			p.ProposerRefreshInterval = val
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
			return err
		}
	}
	return nil
}
