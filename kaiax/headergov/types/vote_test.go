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

func TestHeaderVoteSerialization(t *testing.T) {
	v1 := common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541")
	v2 := common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")

	tcs := []struct {
		serializedVoteData string
		blockNum           uint64
		voteData           *VoteData
	}{
		{serializedVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e6974707269636585ae9f7bcc00", blockNum: 86119166, voteData: &VoteData{BlockNum: 86119166, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(750000000000)}}},
		{serializedVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e69747072696365853a35294400", blockNum: 90355962, voteData: &VoteData{BlockNum: 90355962, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(250000000000)}}},
		{serializedVoteData: "0xed9452d41ca72af615a1ac3301b0a93efa222ecc754196697374616e62756c2e636f6d6d697474656573697a651f", blockNum: 95352567, voteData: &VoteData{BlockNum: 95352567, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.CommitteeSize], Value: uint64(31)}}},
		{serializedVoteData: "0xf83e9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e749336343030303030303030303030303030303030", blockNum: 105629058, voteData: &VoteData{BlockNum: 105629058, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.MintingAmount], Value: "6400000000000000000"}}},
		{serializedVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f34302f3130", blockNum: 105629111, voteData: &VoteData{BlockNum: 105629111, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "50/40/10"}}},
		{serializedVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f32302f3330", blockNum: 118753908, voteData: &VoteData{BlockNum: 118753908, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "50/20/30"}}},
		{serializedVoteData: "0xf8439452d41ca72af615a1ac3301b0a93efa222ecc754198676f7665726e616e63652e676f7665726e696e676e6f646594c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456", blockNum: 126061533, voteData: &VoteData{BlockNum: 126061533, Voter: v1, Param: Param{Name: governance.GovernanceKeyMapReverse[params.GoverningNode], Value: v2}}},
		{serializedVoteData: "0xef94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45698676f7665726e616e63652e646572697665736861696d706c80", blockNum: 127692621, voteData: &VoteData{BlockNum: 127692621, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.DeriveShaImpl], Value: uint64(0)}}},
		{serializedVoteData: "0xe994c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568f6b697037312e67617374617267657483e4e1c0", blockNum: 140916059, voteData: &VoteData{BlockNum: 140916059, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.GasTarget], Value: uint64(15000000)}}},
		{serializedVoteData: "0xf83a94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4569f6b697037312e6d6178626c6f636b67617375736564666f72626173656665658401c9c380", blockNum: 140916152, voteData: &VoteData{BlockNum: 140916152, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.MaxBlockGasUsedForBaseFee], Value: uint64(30000000)}}},
		{serializedVoteData: "0xed94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45696697374616e62756c2e636f6d6d697474656573697a6532", blockNum: 161809335, voteData: &VoteData{BlockNum: 161809335, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.CommitteeSize], Value: uint64(50)}}},
		{serializedVoteData: "0xf83e94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456947265776172642e6d696e74696e67616d6f756e749339363030303030303030303030303030303030", blockNum: 161809370, voteData: &VoteData{BlockNum: 161809370, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.MintingAmount], Value: "9600000000000000000"}}},
		{serializedVoteData: "0xeb94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568c7265776172642e726174696f8835302f32352f3235", blockNum: 161809416, voteData: &VoteData{BlockNum: 161809416, Voter: v2, Param: Param{Name: governance.GovernanceKeyMapReverse[params.Ratio], Value: "50/25/25"}}},
	}

	for i, tc := range tcs {
		actual, err := DeserializeHeaderVote(hexutil.MustDecode(tc.serializedVoteData), tc.blockNum)
		assert.NoError(t, err)
		assert.Equal(t, tc.voteData, actual, fmt.Sprintf("DeserializeHeaderVote() tcs[%d] failed", i))

		serialized, err := tc.voteData.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, tc.serializedVoteData, hexutil.Encode(serialized), fmt.Sprintf("voteData.Serialize() tcs[%d] failed", i))
	}
}
