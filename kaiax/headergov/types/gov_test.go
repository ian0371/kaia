package types

import (
	"fmt"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestHeaderGovSerialization(t *testing.T) {
	tcs := []struct {
		serializedGovData string
		blockNum          uint64
		data              *GovernanceData
	}{
		{serializedGovData: "0xb901c17b22676f7665726e616e63652e676f7665726e616e63656d6f6465223a2273696e676c65222c22676f7665726e616e63652e676f7665726e696e676e6f6465223a22307835326434316361373261663631356131616333333031623061393365666132323265636337353431222c22676f7665726e616e63652e756e69747072696365223a32353030303030303030302c22697374616e62756c2e636f6d6d697474656573697a65223a32322c22697374616e62756c2e65706f6368223a3630343830302c22697374616e62756c2e706f6c696379223a322c227265776172642e64656665727265647478666565223a747275652c227265776172642e6d696e696d756d7374616b65223a2235303030303030222c227265776172642e6d696e74696e67616d6f756e74223a2239363030303030303030303030303030303030222c227265776172642e70726f706f736572757064617465696e74657276616c223a333630302c227265776172642e726174696f223a2233342f35342f3132222c227265776172642e7374616b696e67757064617465696e74657276616c223a38363430302c227265776172642e75736567696e69636f656666223a747275657d", blockNum: 0, data: &GovernanceData{BlockNum: 0, Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.GovernanceMode]:          "single",
			governance.GovernanceKeyMapReverse[params.GoverningNode]:           common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541"),
			governance.GovernanceKeyMapReverse[params.UnitPrice]:               uint64(25000000000),
			governance.GovernanceKeyMapReverse[params.CommitteeSize]:           uint64(22),
			governance.GovernanceKeyMapReverse[params.Epoch]:                   uint64(604800),
			governance.GovernanceKeyMapReverse[params.Policy]:                  uint64(2),
			governance.GovernanceKeyMapReverse[params.DeferredTxFee]:           true,
			governance.GovernanceKeyMapReverse[params.MinimumStake]:            "5000000",
			governance.GovernanceKeyMapReverse[params.MintingAmount]:           "9600000000000000000",
			governance.GovernanceKeyMapReverse[params.ProposerRefreshInterval]: uint64(3600),
			governance.GovernanceKeyMapReverse[params.Ratio]:                   "34/54/12",
			governance.GovernanceKeyMapReverse[params.StakeUpdateInterval]:     uint64(86400),
			governance.GovernanceKeyMapReverse[params.UseGiniCoeff]:            true,
		}}},
		{serializedGovData: "0xa57b22676f7665726e616e63652e756e69747072696365223a3735303030303030303030307d", blockNum: 86486400, data: &GovernanceData{BlockNum: 86486400, Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.UnitPrice]: uint64(750000000000),
		}}},
		{serializedGovData: "0xa57b22676f7665726e616e63652e756e69747072696365223a3235303030303030303030307d", blockNum: 90720000, data: &GovernanceData{BlockNum: 90720000, Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.UnitPrice]: uint64(250000000000),
		}}},
		{serializedGovData: "0x9d7b22697374616e62756c2e636f6d6d697474656573697a65223a33317d", blockNum: 95558400, data: &GovernanceData{BlockNum: 95558400, Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.CommitteeSize]: uint64(31),
		}}},
		{serializedGovData: "0xb8487b227265776172642e6d696e74696e67616d6f756e74223a2236343030303030303030303030303030303030222c227265776172642e726174696f223a2235302f34302f3130227d", blockNum: 105840000, data: &GovernanceData{BlockNum: 105840000, Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.MintingAmount]: "6400000000000000000",
			governance.GovernanceKeyMapReverse[params.Ratio]:         "50/40/10",
		}}},
		{serializedGovData: "0x9b7b227265776172642e726174696f223a2235302f32302f3330227d", blockNum: 119145600, data: &GovernanceData{BlockNum: 119145600, Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.Ratio]: "50/20/30",
		}}},
		{serializedGovData: "0xb8497b22676f7665726e616e63652e676f7665726e696e676e6f6465223a22307863306362653163373730666263653165623737383662666261316163323131356435633061343536227d", blockNum: 126403200, data: &GovernanceData{BlockNum: 126403200, Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.GoverningNode]: common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456"),
		}}},
		{serializedGovData: "0xb8477b226b697037312e676173746172676574223a31353030303030302c226b697037312e6d6178626c6f636b67617375736564666f7262617365666565223a33303030303030307d", blockNum: 140918400, data: &GovernanceData{BlockNum: 140918400, Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.GasTarget]:                 uint64(15000000),
			governance.GovernanceKeyMapReverse[params.MaxBlockGasUsedForBaseFee]: uint64(30000000),
		}}},
		{serializedGovData: "0xb8647b22697374616e62756c2e636f6d6d697474656573697a65223a35302c227265776172642e6d696e74696e67616d6f756e74223a2239363030303030303030303030303030303030222c227265776172642e726174696f223a2235302f32352f3235227d", blockNum: 162086400, data: &GovernanceData{BlockNum: 162086400, Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.CommitteeSize]: uint64(50),
			governance.GovernanceKeyMapReverse[params.MintingAmount]: "9600000000000000000",
			governance.GovernanceKeyMapReverse[params.Ratio]:         "50/25/25",
		}}},
	}

	for i, tc := range tcs {
		actual, err := DeserializeHeaderGov(hexutil.MustDecode(tc.serializedGovData), tc.blockNum)
		assert.NoError(t, err)
		assert.Equal(t, tc.data, actual, fmt.Sprintf("DeserializeHeaderGov() tcs[%d] failed", i))

		// serialized is order dependent, so compare the deserialized result
		serialized, err := actual.Serialize()
		assert.NoError(t, err)
		deserialized, err := DeserializeHeaderGov(serialized, tc.blockNum)
		assert.NoError(t, err)
		assert.Equal(t, tc.data, deserialized, fmt.Sprintf("governanceData.Serialize() tcs[%d] failed", i))
	}
}
