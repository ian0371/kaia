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
)

const (
	// Proposer policy
	RoundRobin = iota
	Sticky
	WeightedRandom
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
		p.GovernanceMode = value.(int)
	case "governance.governingnode":
		p.GoverningNode = value.(common.Address)
	case "governance.govparamcontract":
		p.GovParamContract = value.(common.Address)
	case "istanbul.epoch":
		p.Epoch = value.(uint64)
	case "istanbul.policy":
		p.ProposerPolicy = value.(uint64)
	case "istanbul.committeesize":
		p.CommitteeSize = value.(uint64)
	case "governance.unitprice":
		p.UnitPrice = value.(uint64)
	case "governance.deriveshaimpl":
		p.DeriveShaImpl = value.(uint64)
	case "kip71.lowerboundbasefee":
		p.LowerBoundBaseFee = value.(uint64)
	case "kip71.gastarget":
		p.GasTarget = value.(uint64)
	case "kip71.maxblockgasusedforbasefee":
		p.MaxBlockGasUsedForBaseFee = value.(uint64)
	case "kip71.basefeedenominator":
		p.BaseFeeDenominator = value.(uint64)
	case "kip71.upperboundbasefee":
		p.UpperBoundBaseFee = value.(uint64)
	case "reward.mintingamount":
		p.MintingAmount = value.(*big.Int)
	case "reward.ratio":
		p.Ratio = value.(string)
	case "reward.kip82ratio":
		p.Kip82Ratio = value.(string)
	case "reward.useginicoeff":
		p.UseGiniCoeff = value.(bool)
	case "reward.deferredtxfee":
		p.DeferredTxFee = value.(bool)
	case "reward.minimumstake":
		p.MinimumStake = value.(*big.Int)
	case "reward.stakingupdateinterval":
		p.StakeUpdateInterval = value.(uint64)
	case "reward.proposerupdateinterval":
		p.ProposerRefreshInterval = value.(uint64)
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
