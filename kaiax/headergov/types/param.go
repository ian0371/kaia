package types

import (
	"bytes"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/common"
)

type GovParam struct {
	GovParamSetFieldName string
	Canonicalizer        func(v interface{}) (interface{}, error)
	Validator            func(cv interface{}) (bool, error) // validation on canonical value.

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

var govParams = map[string]GovParam{
	"governance.governancemode": {
		GovParamSetFieldName: "GovernanceMode",
		Canonicalizer:        stringCanonicalizer,
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
		GovParamSetFieldName: "GoverningNode",
		Canonicalizer:        addressCanonicalizer,
		Validator:            nil,
		DefaultValue:         defaultGoverningNode,
		Forbidden:            false,
	},
	"governance.govparamcontract": {
		GovParamSetFieldName: "GovParamContract",
		Canonicalizer:        addressCanonicalizer,
		Validator:            nil,
		DefaultValue:         defaultGovParamContract,
		Forbidden:            false,
	},
	"istanbul.committeesize": {
		GovParamSetFieldName: "CommitteeSize",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil,
		DefaultValue:         defaultCommitteeSize,
		Forbidden:            false,
	},
	"istanbul.policy": {
		GovParamSetFieldName: "ProposerPolicy",
		Canonicalizer:        uint64Canonicalizer,
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
		GovParamSetFieldName: "Epoch",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil,
		DefaultValue:         defaultEpoch,
		Forbidden:            false,
	},
	"reward.ratio": {
		GovParamSetFieldName: "Ratio",
		Canonicalizer:        stringCanonicalizer,
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
		GovParamSetFieldName: "Kip82Ratio",
		Canonicalizer:        stringCanonicalizer,
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
		GovParamSetFieldName: "StakeUpdateInterval",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil,
		DefaultValue:         defaultStakeUpdateInterval,
		Forbidden:            false,
	},
	"reward.proposerupdateinterval": {
		GovParamSetFieldName: "ProposerRefreshInterval",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil,
		DefaultValue:         defaultProposerRefreshInterval,
		Forbidden:            false,
	},
	"reward.mintingamount": {
		GovParamSetFieldName: "MintingAmount",
		Canonicalizer:        bigIntCanonicalizer,
		Validator:            nil,
		DefaultValue:         defaultMintingAmount,
		Forbidden:            false,
	},
	"reward.minimumstake": {
		GovParamSetFieldName: "MinimumStake",
		Canonicalizer:        bigIntCanonicalizer,
		Validator:            nil,
		DefaultValue:         defaultMinimumStake,
		Forbidden:            false,
	},
	"reward.useginicoeff": {
		GovParamSetFieldName: "UseGiniCoeff",
		Canonicalizer:        boolCanonicalizer,
		Validator:            nil,
		DefaultValue:         defaultUseGiniCoeff,
		Forbidden:            false,
	},
	"reward.deferredtxfee": {
		GovParamSetFieldName: "DeferredTxFee",
		Canonicalizer:        boolCanonicalizer,
		Validator:            nil,
		DefaultValue:         defaultDeferredTxFee,
		Forbidden:            false,
	},
	"kip71.lowerboundbasefee": {
		GovParamSetFieldName: "LowerBoundBaseFee",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil, // TODO: validate if lower bound is less than upper bound. or validate in vote API.
		DefaultValue:         defaultLowerBoundBaseFee,
		Forbidden:            false,
	},
	"kip71.upperboundbasefee": {
		GovParamSetFieldName: "UpperBoundBaseFee",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil, // TODO: validate if upper bound is greater than lower bound. or validate in vote API.
		DefaultValue:         defaultUpperBoundBaseFee,
		Forbidden:            false,
	},
	"kip71.gastarget": {
		GovParamSetFieldName: "GasTarget",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil,
		DefaultValue:         defaultGasTarget,
		Forbidden:            false,
	},
	"kip71.maxblockgasusedforbasefee": {
		GovParamSetFieldName: "MaxBlockGasUsedForBaseFee",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil,
		DefaultValue:         defaultMaxBlockGasUsedForBaseFee,
		Forbidden:            false,
	},
	"kip71.basefeedenominator": {
		GovParamSetFieldName: "BaseFeeDenominator",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil,
		DefaultValue:         defaultBaseFeeDenominator,
		Forbidden:            false,
	},
	"governance.deriveshaimpl": {
		GovParamSetFieldName: "DeriveShaImpl",
		Canonicalizer:        uint64Canonicalizer,
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
		GovParamSetFieldName: "UnitPrice",
		Canonicalizer:        uint64Canonicalizer,
		Validator:            nil,
		DefaultValue:         defaultUnitPrice,
		Forbidden:            false,
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
