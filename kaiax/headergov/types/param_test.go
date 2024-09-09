package types

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestParam(t *testing.T) {
	for name, param := range Params {
		assert.NotEmpty(t, param.ParamSetFieldName, name)
		assert.NotEmpty(t, param.Canonicalizer, name)
		// assert.NotEmpty(t, param.FormatChecker, name)
	}
}

func TestUint64Canonicalizer(t *testing.T) {
	tcs := []struct {
		name          string
		input         interface{}
		expected      uint64
		expectedError error
	}{
		{name: "Valid uint64", input: uint64(12345), expected: 12345, expectedError: nil},
		{name: "Valid float64", input: float64(67890), expected: 67890, expectedError: nil},
		{name: "Valid byte slice", input: []byte{0, 0, 0, 0, 0, 0x10, 0, 0}, expected: 0x100000, expectedError: nil},
		{name: "Invalid float64 (not an integer)", input: float64(123.45), expected: 0, expectedError: ErrCanonicalizeFloatToUint64},
		{name: "Invalid type (string)", input: "12345", expected: 0, expectedError: ErrCanonicalizeUint64},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := uint64Canonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.(uint64))
			}
		})
	}
}

func TestBigIntCanonicalizer(t *testing.T) {
	tcs := []struct {
		name          string
		input         interface{}
		expected      *big.Int
		expectedError error
	}{
		{name: "Valid big.Int", input: big.NewInt(12345), expected: big.NewInt(12345), expectedError: nil},
		{name: "Valid string", input: "67890", expected: big.NewInt(67890), expectedError: nil},
		{name: "Valid byte slice", input: []byte("100000"), expected: big.NewInt(100000), expectedError: nil},
		{name: "Invalid string", input: "invalid", expected: nil, expectedError: ErrCanonicalizeStringToBigInt},
		{name: "Invalid byte slice", input: []byte("invalid"), expected: nil, expectedError: ErrCanonicalizeByteToBigInt},
		{name: "Invalid type (int)", input: 12345, expected: nil, expectedError: ErrCanonicalizeBigInt},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := bigIntCanonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.(*big.Int))
			}
		})
	}
}

func TestAddressCanonicalizer(t *testing.T) {
	tcs := []struct {
		name          string
		input         interface{}
		expected      common.Address
		expectedError error
	}{
		{name: "Valid byte slice", input: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, expected: common.HexToAddress("0x0102030405060708090a0b0c0d0e0f1011121314")},
		{name: "Valid hex string", input: "0x1234567890123456789012345678901234567890", expected: common.HexToAddress("0x1234567890123456789012345678901234567890")},
		{name: "Valid common.Address", input: common.HexToAddress("0x1234567890123456789012345678901234567890"), expected: common.HexToAddress("0x1234567890123456789012345678901234567890")},
		{name: "Invalid byte slice length", input: []byte{1, 2, 3}, expectedError: ErrCanonicalizeByteToAddress},
		{name: "Invalid hex string", input: "0xinvalid", expectedError: ErrCanonicalizeStringToAddress},
		{name: "Invalid type", input: 123, expectedError: ErrCanonicalizeToAddress},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := addressCanonicalizer(tc.input)
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result.(common.Address))
			}
		})
	}
}

func TestBoolCanonicalizer(t *testing.T) {
	tcs := []struct {
		name          string
		input         interface{}
		expected      bool
		expectedError error
	}{
		{name: "Valid bool true", input: true, expected: true, expectedError: nil},
		{name: "Valid bool false", input: false, expected: false, expectedError: nil},
		{name: "Valid byte slice true", input: []byte{0x01}, expected: true, expectedError: nil},
		{name: "Valid byte slice false", input: []byte{0x00}, expected: false, expectedError: nil},
		{name: "Invalid byte slice", input: []byte{0x02}, expected: false, expectedError: ErrCanonicalizeByteToBool},
		{name: "Invalid type string", input: "true", expected: false, expectedError: ErrCanonicalizeBool},
		{name: "Invalid type int", input: 1, expected: false, expectedError: ErrCanonicalizeBool},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := boolCanonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.(bool))
			}
		})
	}
}

func TestAddressListCanonicalizer(t *testing.T) {
	tcs := []struct {
		name          string
		input         interface{}
		expected      []common.Address
		expectedError error
	}{
		{name: "Valid single address string", input: "0x1234567890123456789012345678901234567890", expected: []common.Address{common.HexToAddress("0x1234567890123456789012345678901234567890")}},
		{name: "Valid multiple address string", input: "0x1234567890123456789012345678901234567890,0x0987654321098765432109876543210987654321", expected: []common.Address{common.HexToAddress("0x1234567890123456789012345678901234567890"), common.HexToAddress("0x0987654321098765432109876543210987654321")}},
		{name: "Invalid address string", input: "0xinvalid", expectedError: ErrCanonicalizeStringToAddress},
		{name: "Valid byte slice", input: []byte("0x1234567890123456789012345678901234567890"), expected: []common.Address{common.HexToAddress("0x1234567890123456789012345678901234567890")}},
		{name: "Invalid type", input: 123, expectedError: ErrCanonicalizeToAddressList},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := addressListCanonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.([]common.Address))
			}
		})
	}
}
