package types

import (
	"bytes"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/common"
)

type Param struct {
	ParamSetFieldName  string
	Canonicalizer      func(v interface{}) (interface{}, error)
	FormatChecker      func(cv interface{}) bool // validation on canonical value.
	ConsistencyChecker func(cv interface{}) bool // validation on canonical value.

	DefaultValue  interface{}
	VoteForbidden bool
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
		return common.BytesToAddress(v), nil
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
		if float64(uint64(v)) != v {
			return nil, errors.New("value is not an integer")
		}

		return uint64(v), nil
	case uint64:
		return v, nil
	}
	return nil, errors.New("invalid type")
}

func bigIntCanonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		cv, ok := new(big.Int).SetString(string(v), 10)
		if !ok {
			return nil, errors.New("could not canonicalize []byte to big.Int")
		}
		return cv, nil
	case string:
		cv, ok := new(big.Int).SetString(v, 10)
		if !ok {
			return nil, errors.New("could not canonicalize string to big.Int")
		}
		return cv, nil
	}
	return nil, errors.New("could not canonicalize value to big.Int")
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

var Params = map[string]Param{
	"governance.governancemode": {
		ParamSetFieldName: "GovernanceMode",
		Canonicalizer:     stringCanonicalizer,
		FormatChecker: func(cv interface{}) bool {
			switch v := cv.(type) {
			case string:
				if v == "none" || v == "single" || v == "ballot" {
					return true
				}
			}
			return false
		},
		DefaultValue:  defaultGovernanceMode,
		VoteForbidden: true,
	},
	"governance.governingnode": {
		ParamSetFieldName: "GoverningNode",
		Canonicalizer:     addressCanonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultGoverningNode,
		VoteForbidden:     false,
	},
	"governance.govparamcontract": {
		ParamSetFieldName: "GovParamContract",
		Canonicalizer:     addressCanonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultGovParamContract,
		VoteForbidden:     false,
	},
	"istanbul.committeesize": {
		ParamSetFieldName: "CommitteeSize",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultCommitteeSize,
		VoteForbidden:     false,
	},
	"istanbul.policy": {
		ParamSetFieldName: "ProposerPolicy",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker: func(cv interface{}) bool {
			v, ok := cv.(uint64)
			if !ok {
				return false
			}
			return v <= 2
		},
		DefaultValue:  defaultProposerPolicy,
		VoteForbidden: true,
	},
	"istanbul.epoch": {
		ParamSetFieldName: "Epoch",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultEpoch,
		VoteForbidden:     true,
	},
	"reward.ratio": {
		ParamSetFieldName: "Ratio",
		Canonicalizer:     stringCanonicalizer,
		FormatChecker: func(cv interface{}) bool {
			v, ok := cv.(string)
			if !ok {
				return false
			}
			parts := strings.Split(v, "/")
			if len(parts) != 3 {
				return false
			}
			sum := 0
			for _, part := range parts {
				num, err := strconv.Atoi(part)
				if err != nil {
					return false
				}
				if num < 0 {
					return false
				}
				sum += num
			}

			return sum == 100
		},
		DefaultValue:  defaultRatio,
		VoteForbidden: false,
	},
	"reward.kip82ratio": {
		ParamSetFieldName: "Kip82Ratio",
		Canonicalizer:     stringCanonicalizer,
		FormatChecker: func(cv interface{}) bool {
			v, ok := cv.(string)
			if !ok {
				return false
			}
			parts := strings.Split(v, "/")
			if len(parts) != 2 {
				return false
			}
			sum := 0
			for _, part := range parts {
				num, err := strconv.Atoi(part)
				if err != nil {
					return false
				}
				if num < 0 {
					return false
				}
				sum += num
			}

			return sum == 100
		},
		DefaultValue:  defaultKip82Ratio,
		VoteForbidden: false,
	},
	"reward.stakingupdateinterval": {
		ParamSetFieldName: "StakeUpdateInterval",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultStakeUpdateInterval,
		VoteForbidden:     true,
	},
	"reward.proposerupdateinterval": {
		ParamSetFieldName: "ProposerUpdateInterval",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultProposerUpdateInterval,
		VoteForbidden:     true,
	},
	"reward.mintingamount": {
		ParamSetFieldName: "MintingAmount",
		Canonicalizer:     bigIntCanonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultMintingAmount,
		VoteForbidden:     false,
	},
	"reward.minimumstake": {
		ParamSetFieldName: "MinimumStake",
		Canonicalizer:     bigIntCanonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultMinimumStake,
		VoteForbidden:     true,
	},
	"reward.useginicoeff": {
		ParamSetFieldName: "UseGiniCoeff",
		Canonicalizer:     boolCanonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultUseGiniCoeff,
		VoteForbidden:     true,
	},
	"reward.deferredtxfee": {
		ParamSetFieldName: "DeferredTxFee",
		Canonicalizer:     boolCanonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultDeferredTxFee,
		VoteForbidden:     true,
	},
	"kip71.lowerboundbasefee": {
		ParamSetFieldName: "LowerBoundBaseFee",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil, // TODO: validate if lower bound is less than upper bound. or validate in vote API.
		DefaultValue:      defaultLowerBoundBaseFee,
		VoteForbidden:     false,
	},
	"kip71.upperboundbasefee": {
		ParamSetFieldName: "UpperBoundBaseFee",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil, // TODO: validate if upper bound is greater than lower bound. or validate in vote API.
		DefaultValue:      defaultUpperBoundBaseFee,
		VoteForbidden:     false,
	},
	"kip71.gastarget": {
		ParamSetFieldName: "GasTarget",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultGasTarget,
		VoteForbidden:     false,
	},
	"kip71.maxblockgasusedforbasefee": {
		ParamSetFieldName: "MaxBlockGasUsedForBaseFee",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultMaxBlockGasUsedForBaseFee,
		VoteForbidden:     false,
	},
	"kip71.basefeedenominator": {
		ParamSetFieldName: "BaseFeeDenominator",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultBaseFeeDenominator,
		VoteForbidden:     false,
	},
	"governance.deriveshaimpl": {
		ParamSetFieldName: "DeriveShaImpl",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker: func(cv interface{}) bool {
			v, ok := cv.(uint64)
			if !ok {
				return false
			}
			return v <= 2
		},
		DefaultValue:  defaultDeriveShaImpl,
		VoteForbidden: false,
	},
	"governance.unitprice": {
		ParamSetFieldName: "UnitPrice",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     nil,
		DefaultValue:      defaultUnitPrice,
		VoteForbidden:     false,
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
	defaultProposerUpdateInterval    = uint64(3600)  // 1 hour
	defaultDeriveShaImpl             = uint64(0)     // Orig
)

func GetDefaultGovernanceParam() *ParamSet {
	return &ParamSet{
		GovernanceMode:            defaultGovernanceMode,
		GoverningNode:             defaultGoverningNode,
		GovParamContract:          defaultGovParamContract,
		CommitteeSize:             defaultCommitteeSize,
		ProposerPolicy:            defaultProposerPolicy,
		Epoch:                     defaultEpoch,
		Ratio:                     defaultRatio,
		Kip82Ratio:                defaultKip82Ratio,
		StakeUpdateInterval:       defaultStakeUpdateInterval,
		ProposerUpdateInterval:    defaultProposerUpdateInterval,
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
