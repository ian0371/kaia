package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/common"
)

type VoteParam struct {
	GovFieldName  string
	Canonicalizer func(v interface{}) (interface{}, error)
	Validator     func(cv interface{}) (bool, error) // validation on canonical value.

	DefaultValue interface{}
	Forbidden    bool
}

func stringCanonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		return string(v), nil
	case string:
		return v, nil
	}
	return nil, errors.New("invalid type")
}

func addressCanonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		return common.HexToAddress(string(v)), nil
	case string:
		return common.HexToAddress(v), nil
	case common.Address:
		return v, nil
	}
	return nil, errors.New("invalid type")
}

func uint64Canonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		return new(big.Int).SetBytes(v).Uint64(), nil
	case float64:
		return uint64(v), nil
	case uint64:
		return v, nil
	}
	return nil, errors.New("invalid type")
}

func bigIntCanonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		return new(big.Int).SetBytes(v), nil
	case string:
		return new(big.Int).SetBytes([]byte(v)), nil
	}
	return nil, errors.New("invalid type")
}

func boolCanonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		if bytes.Compare(v, []byte{0x01}) == 0 {
			return true, nil
		} else if bytes.Compare(v, []byte{0x00}) == 0 {
			return false, nil
		} else {
			return nil, errors.New("invalid type")
		}
	case bool:
		return v, nil
	}
	return nil, errors.New("invalid type")
}

func addressArrayCanonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		// Handle single address or multiple addresses joined by comma
		addresses := strings.Split(string(v), ",")
		var result []common.Address
		for _, addr := range addresses {
			trimmedAddr := strings.TrimSpace(addr)
			if !common.IsHexAddress(trimmedAddr) {
				return nil, errors.New("invalid address format")
			}
			result = append(result, common.HexToAddress(trimmedAddr))
		}
		return result, nil
	}
	return nil, errors.New("invalid type")
}

var govParams = map[string]VoteParam{
	"governance.governancemode": {
		GovFieldName:  "GovernanceMode",
		Canonicalizer: stringCanonicalizer,
		Validator: func(cv interface{}) (bool, error) {
			switch v := cv.(type) {
			case string:
				if v == "none" || v == "single" || v == "ballot" {
					return true, nil
				}
			}
			return false, errors.New("invalid governance mode")
		},
		DefaultValue: defaultGovernanceMode,
		Forbidden:    false,
	},
	"governance.governingnode": {
		GovFieldName:  "GoverningNode",
		Canonicalizer: addressCanonicalizer,
		Validator:     nil,
		DefaultValue:  defaultGoverningNode,
		Forbidden:     false,
	},
	"governance.govparamcontract": {
		GovFieldName:  "GovParamContract",
		Canonicalizer: addressCanonicalizer,
		Validator:     nil,
		DefaultValue:  defaultGovParamContract,
		Forbidden:     false,
	},
	"istanbul.committeesize": {
		GovFieldName:  "CommitteeSize",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil,
		DefaultValue:  defaultCommitteeSize,
		Forbidden:     false,
	},
	"istanbul.policy": {
		GovFieldName:  "ProposerPolicy",
		Canonicalizer: uint64Canonicalizer,
		Validator: func(cv interface{}) (bool, error) {
			v, ok := cv.(uint64)
			if !ok {
				return false, errors.New("invalid type")
			}
			return v < ProposerPolicy_End, nil
		},
		DefaultValue: defaultProposerPolicy,
		Forbidden:    false,
	},
	"istanbul.epoch": {
		GovFieldName:  "Epoch",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil,
		DefaultValue:  defaultEpoch,
		Forbidden:     false,
	},
	"reward.ratio": {
		GovFieldName:  "Ratio",
		Canonicalizer: stringCanonicalizer,
		Validator: func(cv interface{}) (bool, error) {
			v, ok := cv.(string)
			if !ok {
				return false, errors.New("invalid type")
			}
			parts := strings.Split(v, "/")
			if len(parts) != 3 {
				return false, errors.New("invalid format: must be a/b/c")
			}
			sum := 0
			for _, part := range parts {
				num, err := strconv.Atoi(part)
				if err != nil {
					return false, errors.New("invalid number in ratio")
				}
				if num < 0 {
					return false, errors.New("negative numbers not allowed in ratio")
				}
				sum += num
			}
			if sum != 100 {
				return false, errors.New("sum of ratio parts must equal 100")
			}
			return true, nil
		},
		DefaultValue: defaultRatio,
		Forbidden:    false,
	},
	"reward.kip82ratio": {
		GovFieldName:  "Kip82Ratio",
		Canonicalizer: stringCanonicalizer,
		Validator: func(cv interface{}) (bool, error) {
			v, ok := cv.(string)
			if !ok {
				return false, errors.New("invalid type")
			}
			parts := strings.Split(v, "/")
			if len(parts) != 2 {
				return false, errors.New("invalid format: must be a/b")
			}
			sum := 0
			for _, part := range parts {
				num, err := strconv.Atoi(part)
				if err != nil {
					return false, errors.New("invalid number in ratio")
				}
				if num < 0 {
					return false, errors.New("negative numbers not allowed in ratio")
				}
				sum += num
			}
			if sum != 100 {
				return false, errors.New("sum of ratio parts must equal 100")
			}
			return true, nil
		},
		DefaultValue: defaultKip82Ratio,
		Forbidden:    false,
	},
	"reward.stakingupdateinterval": {
		GovFieldName:  "StakeUpdateInterval",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil,
		DefaultValue:  defaultStakeUpdateInterval,
		Forbidden:     false,
	},
	"reward.proposerupdateinterval": {
		GovFieldName:  "ProposerRefreshInterval",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil,
		DefaultValue:  defaultProposerRefreshInterval,
		Forbidden:     false,
	},
	"reward.mintingamount": {
		GovFieldName:  "MintingAmount",
		Canonicalizer: bigIntCanonicalizer,
		Validator:     nil,
		DefaultValue:  defaultMintingAmount,
		Forbidden:     false,
	},
	"reward.minimumstake": {
		GovFieldName:  "MinimumStake",
		Canonicalizer: bigIntCanonicalizer,
		Validator:     nil,
		DefaultValue:  defaultMinimumStake,
		Forbidden:     false,
	},
	"reward.useginicoeff": {
		GovFieldName:  "UseGiniCoeff",
		Canonicalizer: boolCanonicalizer,
		Validator:     nil,
		DefaultValue:  defaultUseGiniCoeff,
		Forbidden:     false,
	},
	"reward.deferredtxfee": {
		GovFieldName:  "DeferredTxFee",
		Canonicalizer: boolCanonicalizer,
		Validator:     nil,
		DefaultValue:  defaultDeferredTxFee,
		Forbidden:     false,
	},
	"kip71.lowerboundbasefee": {
		GovFieldName:  "LowerBoundBaseFee",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil, // TODO: validate if lower bound is less than upper bound. or validate in vote API.
		DefaultValue:  defaultLowerBoundBaseFee,
		Forbidden:     false,
	},
	"kip71.upperboundbasefee": {
		GovFieldName:  "UpperBoundBaseFee",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil, // TODO: validate if upper bound is greater than lower bound. or validate in vote API.
		DefaultValue:  defaultUpperBoundBaseFee,
		Forbidden:     false,
	},
	"kip71.gastarget": {
		GovFieldName:  "GasTarget",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil,
		DefaultValue:  defaultGasTarget,
		Forbidden:     false,
	},
	"kip71.maxblockgasusedforbasefee": {
		GovFieldName:  "MaxBlockGasUsedForBaseFee",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil,
		DefaultValue:  defaultMaxBlockGasUsedForBaseFee,
		Forbidden:     false,
	},
	"kip71.basefeedenominator": {
		GovFieldName:  "BaseFeeDenominator",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil,
		DefaultValue:  defaultBaseFeeDenominator,
		Forbidden:     false,
	},
	"governance.deriveshaimpl": {
		GovFieldName:  "DeriveShaImpl",
		Canonicalizer: uint64Canonicalizer,
		Validator: func(cv interface{}) (bool, error) {
			v, ok := cv.(uint64)
			if !ok {
				return false, errors.New("invalid type")
			}
			return v <= 2, nil
		},
		DefaultValue: defaultDeriveShaImpl,
		Forbidden:    false,
	},
	"governance.unitprice": {
		GovFieldName:  "UnitPrice",
		Canonicalizer: uint64Canonicalizer,
		Validator:     nil,
		DefaultValue:  defaultUnitPrice,
		Forbidden:     false,
	},
	"governance.addvalidator": {
		Canonicalizer: addressArrayCanonicalizer,
		Validator:     nil,
		DefaultValue:  defaultUnitPrice,
		Forbidden:     false,
	},
	"governance.removevalidator": {
		Canonicalizer: addressArrayCanonicalizer,
		Validator:     nil, // TODO: validate if addresses are validators. or validate in vote API.
		DefaultValue:  defaultUnitPrice,
		Forbidden:     false,
	},
}

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

const (
	// Proposer policy
	RoundRobin = iota
	Sticky
	WeightedRandom
	ProposerPolicy_End
)

var (
	// Default Values: Constants used for getting default values for configuration
	defaultGovernanceMode            = "none"
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

func GetDefaultGovernanceParam() *GovParamSet {
	return &GovParamSet{
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
