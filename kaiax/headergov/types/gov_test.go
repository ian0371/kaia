package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/stretchr/testify/assert"
)

var _ GovData = (*govData)(nil)

func TestNewGov(t *testing.T) {
	tcs := []struct {
		name    string
		value   interface{}
		invalid bool
	}{
		{name: "istanbul.epoch", value: uint64(30000), invalid: false},
		{name: "istanbul.epoch", value: float64(30000.00), invalid: false},
		{name: "istanbul.policy", value: uint64(0), invalid: false},
		{name: "istanbul.policy", value: uint64(1), invalid: false},
		{name: "istanbul.policy", value: uint64(2), invalid: false},
		{name: "governance.governancemode", value: "none", invalid: false},
		{name: "governance.governancemode", value: "single", invalid: false},
		{name: "governance.governancemode", value: 0, invalid: true},
		{name: "governance.governancemode", value: 1, invalid: true},
		{name: "governance.governancemode", value: "unexpected", invalid: true},
		{name: "reward.useginicoeff", value: true, invalid: false},
		{name: "reward.useginicoeff", value: false, invalid: false},
		{name: "reward.useginicoeff", value: []byte{0}, invalid: false},
		{name: "reward.useginicoeff", value: []byte{1}, invalid: false},
		{name: "reward.useginicoeff", value: "false", invalid: true},
		{name: "reward.useginicoeff", value: 0, invalid: true},
		{name: "reward.useginicoeff", value: 1, invalid: true},
		{name: "reward.minimumstake", value: "2000000000000000000000000", invalid: false},
		{name: "reward.minimumstake", value: 200000000000000, invalid: true},
		{name: "reward.minimumstake", value: "-1", invalid: true},
		{name: "reward.minimumstake", value: "0", invalid: false},
		{name: "reward.minimumstake", value: 0, invalid: true},
		{name: "reward.minimumstake", value: 1.1, invalid: true},
		{name: "reward.stakingupdateinterval", value: uint64(20), invalid: false},
		{name: "reward.stakingupdateinterval", value: float64(20.0), invalid: false},
		{name: "reward.stakingupdateinterval", value: float64(20.2), invalid: true},
		{name: "reward.stakingupdateinterval", value: "20", invalid: true},
		{name: "reward.proposerupdateinterval", value: uint64(20), invalid: false},
		{name: "reward.proposerupdateinterval", value: float64(20.0), invalid: false},
		{name: "reward.proposerupdateinterval", value: float64(20.2), invalid: true},
		{name: "reward.proposerupdateinterval", value: "20", invalid: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			vote := NewGovData(map[string]interface{}{tc.name: tc.value})
			if tc.invalid {
				assert.Nil(t, vote)
			} else {
				assert.NotNil(t, vote)
			}
		})
	}
}
func TestHeaderGovSerialization(t *testing.T) {
	tcs := []struct {
		serializedGovData string
		blockNum          uint64
		data              GovData
	}{
		{serializedGovData: "0xb8b67b227265776172642e6b69703832726174696f223a2233332f3637222c226b697037312e6c6f776572626f756e6462617365666565223a32353030303030303030302c226b697037312e7570706572626f756e6462617365666565223a3735303030303030303030302c226b697037312e6261736566656564656e6f6d696e61746f72223a3130302c227265776172642e6d696e74696e67616d6f756e74223a2239363030303030303030303030303030303030227d", blockNum: 604800, data: NewGovData(map[string]interface{}{
			"reward.kip82ratio":        "33/67",
			"kip71.lowerboundbasefee":  uint64(25000000000),
			"kip71.upperboundbasefee":  uint64(750000000000),
			"kip71.basefeedenominator": uint64(100),
			"reward.mintingamount":     new(big.Int).SetUint64(9.6e18),
		})},

		///// Real mainnet governance data.
		{serializedGovData: "0xb901c17b22676f7665726e616e63652e676f7665726e616e63656d6f6465223a2273696e676c65222c22676f7665726e616e63652e676f7665726e696e676e6f6465223a22307835326434316361373261663631356131616333333031623061393365666132323265636337353431222c22676f7665726e616e63652e756e69747072696365223a32353030303030303030302c22697374616e62756c2e636f6d6d697474656573697a65223a32322c22697374616e62756c2e65706f6368223a3630343830302c22697374616e62756c2e706f6c696379223a322c227265776172642e64656665727265647478666565223a747275652c227265776172642e6d696e696d756d7374616b65223a2235303030303030222c227265776172642e6d696e74696e67616d6f756e74223a2239363030303030303030303030303030303030222c227265776172642e70726f706f736572757064617465696e74657276616c223a333630302c227265776172642e726174696f223a2233342f35342f3132222c227265776172642e7374616b696e67757064617465696e74657276616c223a38363430302c227265776172642e75736567696e69636f656666223a747275657d", blockNum: 0, data: NewGovData(map[string]interface{}{
			"governance.governancemode":     "single",
			"governance.governingnode":      common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541"),
			"governance.unitprice":          uint64(25000000000),
			"istanbul.committeesize":        uint64(22),
			"istanbul.epoch":                uint64(604800),
			"istanbul.policy":               uint64(2),
			"reward.deferredtxfee":          true,
			"reward.minimumstake":           big.NewInt(5e6),
			"reward.mintingamount":          new(big.Int).SetUint64(9.6e18),
			"reward.proposerupdateinterval": uint64(3600),
			"reward.ratio":                  "34/54/12",
			"reward.stakingupdateinterval":  uint64(86400),
			"reward.useginicoeff":           true,
		})},
		{serializedGovData: "0xa57b22676f7665726e616e63652e756e69747072696365223a3735303030303030303030307d", blockNum: 86486400, data: NewGovData(map[string]interface{}{
			"governance.unitprice": uint64(750000000000),
		})},
		{serializedGovData: "0xa57b22676f7665726e616e63652e756e69747072696365223a3235303030303030303030307d", blockNum: 90720000, data: NewGovData(map[string]interface{}{
			"governance.unitprice": uint64(250000000000),
		})},
		{serializedGovData: "0x9d7b22697374616e62756c2e636f6d6d697474656573697a65223a33317d", blockNum: 95558400, data: NewGovData(map[string]interface{}{
			"istanbul.committeesize": uint64(31),
		})},
		{serializedGovData: "0xb8487b227265776172642e6d696e74696e67616d6f756e74223a2236343030303030303030303030303030303030222c227265776172642e726174696f223a2235302f34302f3130227d", blockNum: 105840000, data: NewGovData(map[string]interface{}{
			"reward.mintingamount": big.NewInt(6.4e18),
			"reward.ratio":         "50/40/10",
		})},
		{serializedGovData: "0x9b7b227265776172642e726174696f223a2235302f32302f3330227d", blockNum: 119145600, data: NewGovData(map[string]interface{}{
			"reward.ratio": "50/20/30",
		})},
		{serializedGovData: "0xb8497b22676f7665726e616e63652e676f7665726e696e676e6f6465223a22307863306362653163373730666263653165623737383662666261316163323131356435633061343536227d", blockNum: 126403200, data: NewGovData(map[string]interface{}{
			"governance.governingnode": common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456"),
		})},
		{serializedGovData: "0xb8477b226b697037312e676173746172676574223a31353030303030302c226b697037312e6d6178626c6f636b67617375736564666f7262617365666565223a33303030303030307d", blockNum: 140918400, data: NewGovData(map[string]interface{}{
			"kip71.gastarget":                 uint64(15000000),
			"kip71.maxblockgasusedforbasefee": uint64(30000000),
		})},
		{serializedGovData: "0xb8647b22697374616e62756c2e636f6d6d697474656573697a65223a35302c227265776172642e6d696e74696e67616d6f756e74223a2239363030303030303030303030303030303030222c227265776172642e726174696f223a2235302f32352f3235227d", blockNum: 162086400, data: NewGovData(map[string]interface{}{
			"istanbul.committeesize": uint64(50),
			"reward.mintingamount":   new(big.Int).SetUint64(9.6e18),
			"reward.ratio":           "50/25/25",
		})},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("TestCase_%d", tc.blockNum), func(t *testing.T) {
			// Test deserialization
			actual, err := DeserializeHeaderGov(hexutil.MustDecode(tc.serializedGovData), tc.blockNum)
			assert.NoError(t, err)
			assert.Equal(t, tc.data, actual, "DeserializeHeaderGov() failed")

			// Test serialization
			serialized, err := actual.Serialize()
			assert.NoError(t, err)
			deserialized, err := DeserializeHeaderGov(serialized, tc.blockNum)
			assert.NoError(t, err)
			assert.Equal(t, tc.data, deserialized, "governanceData.Serialize() failed")
		})
	}
}
