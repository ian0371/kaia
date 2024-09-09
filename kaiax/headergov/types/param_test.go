package types

import (
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
	testCases := []struct {
		name          string
		input         interface{}
		expected      uint64
		expectedError error
	}{
		{
			name:          "Valid uint64",
			input:         uint64(12345),
			expected:      12345,
			expectedError: nil,
		},
		{
			name:          "Valid float64",
			input:         float64(67890),
			expected:      67890,
			expectedError: nil,
		},
		{
			name:          "Valid byte slice",
			input:         []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x86, 0xA0},
			expected:      100000,
			expectedError: nil,
		},
		{
			name:          "Invalid float64 (not an integer)",
			input:         float64(123.45),
			expected:      0,
			expectedError: ErrCanonicalizeFloatToUint64,
		},
		{
			name:          "Invalid type (string)",
			input:         "12345",
			expected:      0,
			expectedError: ErrCanonicalizeUint64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := uint64Canonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.(uint64))
			}
		})
	}
}

func TestAddressListCanonicalizer(t *testing.T) {
	testCases := []struct {
		name          string
		input         interface{}
		expected      []common.Address
		expectedError error
	}{
		{
			name:     "Valid single address string",
			input:    "0x1234567890123456789012345678901234567890",
			expected: []common.Address{common.HexToAddress("0x1234567890123456789012345678901234567890")},
		},
		{
			name:  "Valid multiple address string",
			input: "0x1234567890123456789012345678901234567890,0x0987654321098765432109876543210987654321",
			expected: []common.Address{
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
				common.HexToAddress("0x0987654321098765432109876543210987654321"),
			},
		},
		{
			name:          "Invalid address string",
			input:         "0xinvalid",
			expectedError: ErrCanonicalizeStringToAddress,
		},
		{
			name:     "Valid byte slice",
			input:    []byte("0x1234567890123456789012345678901234567890"),
			expected: []common.Address{common.HexToAddress("0x1234567890123456789012345678901234567890")},
		},
		{
			name:          "Invalid type",
			input:         123,
			expectedError: ErrCanonicalizeToAddressList,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := addressListCanonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.([]common.Address))
			}
		})
	}
}
