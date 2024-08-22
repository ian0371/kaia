package headergov

import (
	"fmt"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestHeaderVoteSerialization(t *testing.T) {
	v1 := common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541")
	v2 := common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")

	tcs := []struct {
		blockVoteData string
		blockNum      uint64
		data          *VoteData
	}{
		{blockVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e6974707269636585ae9f7bcc00", blockNum: 86119166, data: &VoteData{BlockNum: 86119166, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(750000000000)}}},
		{blockVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e69747072696365853a35294400", blockNum: 90355962, data: &VoteData{BlockNum: 90355962, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(250000000000)}}},
		{blockVoteData: "0xed9452d41ca72af615a1ac3301b0a93efa222ecc754196697374616e62756c2e636f6d6d697474656573697a651f", blockNum: 95352567, data: &VoteData{BlockNum: 95352567, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.CommitteeSize], Value: uint64(31)}}},
		{blockVoteData: "0xf83e9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e749336343030303030303030303030303030303030", blockNum: 105629058, data: &VoteData{BlockNum: 105629058, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.MintingAmount], Value: "6400000000000000000"}}},
		{blockVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f34302f3130", blockNum: 105629111, data: &VoteData{BlockNum: 105629111, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "50/40/10"}}},
		{blockVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f32302f3330", blockNum: 118753908, data: &VoteData{BlockNum: 118753908, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "50/20/30"}}},
		{blockVoteData: "0xf8439452d41ca72af615a1ac3301b0a93efa222ecc754198676f7665726e616e63652e676f7665726e696e676e6f646594c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456", blockNum: 126061533, data: &VoteData{BlockNum: 126061533, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.GoverningNode], Value: v2}}},
		{blockVoteData: "0xef94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45698676f7665726e616e63652e646572697665736861696d706c80", blockNum: 127692621, data: &VoteData{BlockNum: 127692621, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.DeriveShaImpl], Value: uint64(0)}}},
		{blockVoteData: "0xe994c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568f6b697037312e67617374617267657483e4e1c0", blockNum: 140916059, data: &VoteData{BlockNum: 140916059, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.GasTarget], Value: uint64(15000000)}}},
		{blockVoteData: "0xf83a94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4569f6b697037312e6d6178626c6f636b67617375736564666f72626173656665658401c9c380", blockNum: 140916152, data: &VoteData{BlockNum: 140916152, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.MaxBlockGasUsedForBaseFee], Value: uint64(30000000)}}},
		{blockVoteData: "0xed94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45696697374616e62756c2e636f6d6d697474656573697a6532", blockNum: 161809335, data: &VoteData{BlockNum: 161809335, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.CommitteeSize], Value: uint64(50)}}},
		{blockVoteData: "0xf83e94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456947265776172642e6d696e74696e67616d6f756e749339363030303030303030303030303030303030", blockNum: 161809370, data: &VoteData{BlockNum: 161809370, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.MintingAmount], Value: "9600000000000000000"}}},
		{blockVoteData: "0xeb94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568c7265776172642e726174696f8835302f32352f3235", blockNum: 161809416, data: &VoteData{BlockNum: 161809416, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "50/25/25"}}},
	}

	for i, tc := range tcs {
		actual, err := deserializeHeaderVote(hexutil.MustDecode(tc.blockVoteData), tc.blockNum)
		assert.NoError(t, err)
		assert.Equal(t, tc.data, actual, fmt.Sprintf("deserializeHeaderVote tcs[%d] failed", i))

		serialized, err := serializeVoteData(tc.data)
		assert.NoError(t, err)
		assert.Equal(t, tc.blockVoteData, hexutil.Encode(serialized), fmt.Sprintf("serializeVoteData tcs[%d] failed", i))
	}
}

func TestHeaderGovSerialization(t *testing.T) {
	tcs := []struct {
		blockGovData string
		blockNum     uint64
		data         *GovernanceData
	}{
		{blockGovData: "0xb901c17b22676f7665726e616e63652e676f7665726e616e63656d6f6465223a2273696e676c65222c22676f7665726e616e63652e676f7665726e696e676e6f6465223a22307835326434316361373261663631356131616333333031623061393365666132323265636337353431222c22676f7665726e616e63652e756e69747072696365223a32353030303030303030302c22697374616e62756c2e636f6d6d697474656573697a65223a32322c22697374616e62756c2e65706f6368223a3630343830302c22697374616e62756c2e706f6c696379223a322c227265776172642e64656665727265647478666565223a747275652c227265776172642e6d696e696d756d7374616b65223a2235303030303030222c227265776172642e6d696e74696e67616d6f756e74223a2239363030303030303030303030303030303030222c227265776172642e70726f706f736572757064617465696e74657276616c223a333630302c227265776172642e726174696f223a2233342f35342f3132222c227265776172642e7374616b696e67757064617465696e74657276616c223a38363430302c227265776172642e75736567696e69636f656666223a747275657d", blockNum: 0, data: &GovernanceData{BlockNum: 0, Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.GovernanceMode]:          Param{Name: governance.GovernanceKeyMapReverse[params.GovernanceMode], Value: "single"},
			governance.GovernanceKeyMapReverse[params.GoverningNode]:           Param{Name: governance.GovernanceKeyMapReverse[params.GoverningNode], Value: common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541")},
			governance.GovernanceKeyMapReverse[params.UnitPrice]:               Param{Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(25000000000)},
			governance.GovernanceKeyMapReverse[params.CommitteeSize]:           Param{Name: governance.GovernanceKeyMapReverse[params.CommitteeSize], Value: uint64(22)},
			governance.GovernanceKeyMapReverse[params.Epoch]:                   Param{Name: governance.GovernanceKeyMapReverse[params.Epoch], Value: uint64(604800)},
			governance.GovernanceKeyMapReverse[params.Policy]:                  Param{Name: governance.GovernanceKeyMapReverse[params.Policy], Value: uint64(2)},
			governance.GovernanceKeyMapReverse[params.DeferredTxFee]:           Param{Name: governance.GovernanceKeyMapReverse[params.DeferredTxFee], Value: true},
			governance.GovernanceKeyMapReverse[params.MinimumStake]:            Param{Name: governance.GovernanceKeyMapReverse[params.MinimumStake], Value: "5000000"},
			governance.GovernanceKeyMapReverse[params.MintingAmount]:           Param{Name: governance.GovernanceKeyMapReverse[params.MintingAmount], Value: "9600000000000000000"},
			governance.GovernanceKeyMapReverse[params.ProposerRefreshInterval]: Param{Name: governance.GovernanceKeyMapReverse[params.ProposerRefreshInterval], Value: uint64(3600)},
			governance.GovernanceKeyMapReverse[params.Ratio]:                   Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "34/54/12"},
			governance.GovernanceKeyMapReverse[params.StakeUpdateInterval]:     Param{Name: governance.GovernanceKeyMapReverse[params.StakeUpdateInterval], Value: uint64(86400)},
			governance.GovernanceKeyMapReverse[params.UseGiniCoeff]:            Param{Name: governance.GovernanceKeyMapReverse[params.UseGiniCoeff], Value: true},
		}}},
		{blockGovData: "0xa57b22676f7665726e616e63652e756e69747072696365223a3735303030303030303030307d", blockNum: 86486400, data: &GovernanceData{BlockNum: 86486400, Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.UnitPrice]: Param{Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(750000000000)},
		}}},
		{blockGovData: "0xa57b22676f7665726e616e63652e756e69747072696365223a3235303030303030303030307d", blockNum: 90720000, data: &GovernanceData{BlockNum: 90720000, Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.UnitPrice]: Param{Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(250000000000)},
		}}},
		{blockGovData: "0x9d7b22697374616e62756c2e636f6d6d697474656573697a65223a33317d", blockNum: 95558400, data: &GovernanceData{BlockNum: 95558400, Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.CommitteeSize]: Param{Name: governance.GovernanceKeyMapReverse[params.CommitteeSize], Value: uint64(31)},
		}}},
		{blockGovData: "0xb8487b227265776172642e6d696e74696e67616d6f756e74223a2236343030303030303030303030303030303030222c227265776172642e726174696f223a2235302f34302f3130227d", blockNum: 105840000, data: &GovernanceData{BlockNum: 105840000, Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.MintingAmount]: Param{Name: governance.GovernanceKeyMapReverse[params.MintingAmount], Value: "6400000000000000000"},
			governance.GovernanceKeyMapReverse[params.Ratio]:         Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "50/40/10"},
		}}},
		{blockGovData: "0x9b7b227265776172642e726174696f223a2235302f32302f3330227d", blockNum: 119145600, data: &GovernanceData{BlockNum: 119145600, Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.Ratio]: Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "50/20/30"},
		}}},
		{blockGovData: "0xb8497b22676f7665726e616e63652e676f7665726e696e676e6f6465223a22307863306362653163373730666263653165623737383662666261316163323131356435633061343536227d", blockNum: 126403200, data: &GovernanceData{BlockNum: 126403200, Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.GoverningNode]: Param{Name: governance.GovernanceKeyMapReverse[params.GoverningNode], Value: common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")},
		}}},
		{blockGovData: "0xb8477b226b697037312e676173746172676574223a31353030303030302c226b697037312e6d6178626c6f636b67617375736564666f7262617365666565223a33303030303030307d", blockNum: 140918400, data: &GovernanceData{BlockNum: 140918400, Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.GasTarget]:                 Param{Name: governance.GovernanceKeyMapReverse[params.GasTarget], Value: uint64(15000000)},
			governance.GovernanceKeyMapReverse[params.MaxBlockGasUsedForBaseFee]: Param{Name: governance.GovernanceKeyMapReverse[params.MaxBlockGasUsedForBaseFee], Value: uint64(30000000)},
		}}},
		{blockGovData: "0xb8647b22697374616e62756c2e636f6d6d697474656573697a65223a35302c227265776172642e6d696e74696e67616d6f756e74223a2239363030303030303030303030303030303030222c227265776172642e726174696f223a2235302f32352f3235227d", blockNum: 162086400, data: &GovernanceData{BlockNum: 162086400, Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.CommitteeSize]: Param{Name: governance.GovernanceKeyMapReverse[params.CommitteeSize], Value: uint64(50)},
			governance.GovernanceKeyMapReverse[params.MintingAmount]: Param{Name: governance.GovernanceKeyMapReverse[params.MintingAmount], Value: "9600000000000000000"},
			governance.GovernanceKeyMapReverse[params.Ratio]:         Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "50/25/25"},
		}}},
	}

	for i, tc := range tcs {
		actual, err := deserializeHeaderGov(hexutil.MustDecode(tc.blockGovData), tc.blockNum)
		assert.NoError(t, err)
		assert.Equal(t, tc.data, actual, fmt.Sprintf("deserializeHeaderGov tcs[%d] failed", i))
	}
}
