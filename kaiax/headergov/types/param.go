package types

import (
	"bytes"
	"math/big"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/common"
)

type Param struct {
	ParamSetFieldName string
	Canonicalizer     func(v interface{}) (interface{}, error)
	FormatChecker     func(cv interface{}) bool // validation on canonical value.

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
	return nil, ErrCanonicalizeString
}

func addressCanonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		if len(v) != common.AddressLength {
			return nil, ErrCanonicalizeByteToAddress
		}
		return common.BytesToAddress(v), nil
	case string:
		if !common.IsHexAddress(v) {
			return nil, ErrCanonicalizeStringToAddress
		}
		return common.HexToAddress(v), nil
	case common.Address:
		return v, nil
	}
	return nil, ErrCanonicalizeToAddress
}

func addressListCanonicalizer(v interface{}) (interface{}, error) {
	stringToAddressList := func(v string) ([]common.Address, error) {
		ret := []common.Address{}
		for _, address := range strings.Split(v, ",") {
			if !common.IsHexAddress(address) {
				return nil, ErrCanonicalizeStringToAddress
			}
			ret = append(ret, common.HexToAddress(address))
		}
		return ret, nil
	}

	switch v := v.(type) {
	case []byte:
		return stringToAddressList(string(v))
	case string:
		return stringToAddressList(v)
	}
	return nil, ErrCanonicalizeToAddressList
}

func uint64Canonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		if len(v) > 8 {
			return nil, ErrCanonicalizeByteToUint64
		}
		return new(big.Int).SetBytes(v).Uint64(), nil
	case float64:
		if float64(uint64(v)) != v {
			return nil, ErrCanonicalizeFloatToUint64
		}

		return uint64(v), nil
	case uint64:
		return v, nil
	}
	return nil, ErrCanonicalizeUint64
}

func bigIntCanonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		cv, ok := new(big.Int).SetString(string(v), 10)
		if !ok {
			return nil, ErrCanonicalizeByteToBigInt
		}
		return cv, nil
	case string:
		cv, ok := new(big.Int).SetString(v, 10)
		if !ok {
			return nil, ErrCanonicalizeStringToBigInt
		}
		return cv, nil
	case *big.Int:
		return v, nil
	}
	return nil, ErrCanonicalizeBigInt
}

func boolCanonicalizer(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []byte:
		if bytes.Equal(v, []byte{0x01}) {
			return true, nil
		} else if bytes.Equal(v, []byte{0x00}) {
			return false, nil
		} else {
			return nil, ErrCanonicalizeByteToBool
		}
	case bool:
		return v, nil
	}
	return nil, ErrCanonicalizeBool
}

func noopFormatChecker(cv interface{}) bool {
	return true
}

var Params = map[string]Param{
	"governance.governancemode": {
		ParamSetFieldName: "GovernanceMode",
		Canonicalizer:     stringCanonicalizer,
		FormatChecker: func(cv interface{}) bool {
			v, ok := cv.(string)
			if !ok {
				return false
			}
			if v == "none" || v == "single" {
				return true
			}
			return false
		},
		DefaultValue:  "none",
		VoteForbidden: true,
	},
	"governance.governingnode": {
		ParamSetFieldName: "GoverningNode",
		Canonicalizer:     addressCanonicalizer,
		FormatChecker: func(cv interface{}) bool {
			_, ok := cv.(common.Address)
			return ok
		},
		DefaultValue:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		VoteForbidden: false,
	},
	"governance.govparamcontract": {
		ParamSetFieldName: "GovParamContract",
		Canonicalizer:     addressCanonicalizer,
		FormatChecker: func(cv interface{}) bool {
			_, ok := cv.(common.Address)
			return ok
		},
		DefaultValue:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		VoteForbidden: false,
	},
	"istanbul.committeesize": {
		ParamSetFieldName: "CommitteeSize",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker: func(cv interface{}) bool {
			v, ok := cv.(uint64)
			if !ok {
				return false
			}
			return v > 0
		},
		DefaultValue:  uint64(21),
		VoteForbidden: false,
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
		DefaultValue:  uint64(RoundRobin),
		VoteForbidden: true,
	},
	"istanbul.epoch": {
		ParamSetFieldName: "Epoch",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(604800),
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
		DefaultValue:  "100/0/0",
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
		DefaultValue:  "20/80",
		VoteForbidden: false,
	},
	"reward.stakingupdateinterval": {
		ParamSetFieldName: "StakingUpdateInterval",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(86400),
		VoteForbidden:     true,
	},
	"reward.proposerupdateinterval": {
		ParamSetFieldName: "ProposerUpdateInterval",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(3600),
		VoteForbidden:     true,
	},
	"reward.mintingamount": {
		ParamSetFieldName: "MintingAmount",
		Canonicalizer:     bigIntCanonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      big.NewInt(0),
		VoteForbidden:     false,
	},
	"reward.minimumstake": {
		ParamSetFieldName: "MinimumStake",
		Canonicalizer:     bigIntCanonicalizer,
		FormatChecker: func(cv interface{}) bool {
			v, ok := cv.(*big.Int)
			if !ok {
				return false
			}
			return v.Sign() >= 0
		},
		DefaultValue:  big.NewInt(2000000),
		VoteForbidden: true,
	},
	"reward.useginicoeff": {
		ParamSetFieldName: "UseGiniCoeff",
		Canonicalizer:     boolCanonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      false,
		VoteForbidden:     true,
	},
	"reward.deferredtxfee": {
		ParamSetFieldName: "DeferredTxFee",
		Canonicalizer:     boolCanonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      false,
		VoteForbidden:     true,
	},
	"kip71.lowerboundbasefee": {
		ParamSetFieldName: "LowerBoundBaseFee",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(25000000000),
		VoteForbidden:     false,
	},
	"kip71.upperboundbasefee": {
		ParamSetFieldName: "UpperBoundBaseFee",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(750000000000),
		VoteForbidden:     false,
	},
	"kip71.gastarget": {
		ParamSetFieldName: "GasTarget",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(30000000),
		VoteForbidden:     false,
	},
	"kip71.maxblockgasusedforbasefee": {
		ParamSetFieldName: "MaxBlockGasUsedForBaseFee",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(60000000),
		VoteForbidden:     false,
	},
	"kip71.basefeedenominator": {
		ParamSetFieldName: "BaseFeeDenominator",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker: func(cv interface{}) bool {
			v, ok := cv.(uint64)
			return ok && v != 0
		},
		DefaultValue:  uint64(20),
		VoteForbidden: false,
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
		DefaultValue:  uint64(0),
		VoteForbidden: false,
	},
	"governance.unitprice": {
		ParamSetFieldName: "UnitPrice",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(250e9),
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
