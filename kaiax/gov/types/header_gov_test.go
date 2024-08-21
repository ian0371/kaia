package types

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllParamsHistoryMap(t *testing.T) {
	epoch := uint64(4)
	koreHf := uint64(10000)
	gov := []governanceData{
		0: {
			blockNum: 0,
			params: []param{
				{
					name:  "param1",
					value: 0,
				},
			},
		},
		4: {
			blockNum: 4,
			params: []param{
				{
					name:  "param1",
					value: 100,
				},
			},
		},
	}

	param1 := getAllParamsHistory(gov, epoch, koreHf)["param1"]
	assert.Equal(t, 0, param1.GetItem(0).value)
	assert.Equal(t, 0, param1.GetItem(8).value)
	assert.Equal(t, 100, param1.GetItem(9).value)
}

func TestEffectiveParams(t *testing.T) {
	koreHf := uint64(999999999)
	gov := map[uint64]governanceData{
		0: {
			blockNum: 0,
			params: []param{
				{
					name:  "param1",
					value: 0,
				},
			},
		},
		604800: {
			blockNum: 604800,
			params: []param{
				{
					name:  "param1",
					value: 100,
				},
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()
	h := NewHeaderGovernanceReader(chain, db, koreHf)
	for num, gov := range gov {
		require.NoError(t, h.AddGov(num, gov))
	}

	assert.Equal(t, 0, h.EffectiveParams(0)["param1"].value)
	assert.Equal(t, 0, h.EffectiveParams(604800 * 2)["param1"].value)
	assert.Equal(t, 100, h.EffectiveParams(604800*2 + 1)["param1"].value)

	koreHf = 0
	db = database.NewMemDB()
	h = NewHeaderGovernanceReader(chain, db, koreHf)
	for num, gov := range gov {
		assert.NoError(t, h.AddGov(num, gov))
	}

	assert.Equal(t, 0, h.EffectiveParams(0)["param1"].value)
	assert.Equal(t, 0, h.EffectiveParams(604800*2 - 1)["param1"].value)
	assert.Equal(t, 100, h.EffectiveParams(604800 * 2)["param1"].value)
}

func TestHeaderVote(t *testing.T) {
	v1 := common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541")
	v2 := common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")

	tcs := []struct {
		blockVoteData string
		blockNum      uint64
		data          *voteData
	}{
		{blockVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e6974707269636585ae9f7bcc00", blockNum: 86119166, data: &voteData{blockNum: 86119166, voter: v1, param: param{name: governance.GovernanceKeyMapReverse[params.UnitPrice], value: uint64(750000000000)}}},
		{blockVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e69747072696365853a35294400", blockNum: 90355962, data: &voteData{blockNum: 90355962, voter: v1, param: param{name: governance.GovernanceKeyMapReverse[params.UnitPrice], value: uint64(250000000000)}}},
		{blockVoteData: "0xed9452d41ca72af615a1ac3301b0a93efa222ecc754196697374616e62756c2e636f6d6d697474656573697a651f", blockNum: 95352567, data: &voteData{blockNum: 95352567, voter: v1, param: param{name: governance.GovernanceKeyMapReverse[params.CommitteeSize], value: uint64(31)}}},
		{blockVoteData: "0xf83e9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e749336343030303030303030303030303030303030", blockNum: 105629058, data: &voteData{blockNum: 105629058, voter: v1, param: param{name: governance.GovernanceKeyMapReverse[params.MintingAmount], value: "6400000000000000000"}}},
		{blockVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f34302f3130", blockNum: 105629111, data: &voteData{blockNum: 105629111, voter: v1, param: param{name: governance.GovernanceKeyMapReverse[params.Ratio], value: "50/40/10"}}},
		{blockVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f32302f3330", blockNum: 118753908, data: &voteData{blockNum: 118753908, voter: v1, param: param{name: governance.GovernanceKeyMapReverse[params.Ratio], value: "50/20/30"}}},
		{blockVoteData: "0xf8439452d41ca72af615a1ac3301b0a93efa222ecc754198676f7665726e616e63652e676f7665726e696e676e6f646594c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456", blockNum: 126061533, data: &voteData{blockNum: 126061533, voter: v1, param: param{name: governance.GovernanceKeyMapReverse[params.GoverningNode], value: v2}}},
		{blockVoteData: "0xef94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45698676f7665726e616e63652e646572697665736861696d706c80", blockNum: 127692621, data: &voteData{blockNum: 127692621, voter: v2, param: param{name: governance.GovernanceKeyMapReverse[params.DeriveShaImpl], value: uint64(0)}}},
		{blockVoteData: "0xe994c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568f6b697037312e67617374617267657483e4e1c0", blockNum: 140916059, data: &voteData{blockNum: 140916059, voter: v2, param: param{name: governance.GovernanceKeyMapReverse[params.GasTarget], value: uint64(15000000)}}},
		{blockVoteData: "0xf83a94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4569f6b697037312e6d6178626c6f636b67617375736564666f72626173656665658401c9c380", blockNum: 140916152, data: &voteData{blockNum: 140916152, voter: v2, param: param{name: governance.GovernanceKeyMapReverse[params.MaxBlockGasUsedForBaseFee], value: uint64(30000000)}}},
		{blockVoteData: "0xed94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45696697374616e62756c2e636f6d6d697474656573697a6532", blockNum: 161809335, data: &voteData{blockNum: 161809335, voter: v2, param: param{name: governance.GovernanceKeyMapReverse[params.CommitteeSize], value: uint64(50)}}},
		{blockVoteData: "0xf83e94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456947265776172642e6d696e74696e67616d6f756e749339363030303030303030303030303030303030", blockNum: 161809370, data: &voteData{blockNum: 161809370, voter: v2, param: param{name: governance.GovernanceKeyMapReverse[params.MintingAmount], value: "9600000000000000000"}}},
		{blockVoteData: "0xeb94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568c7265776172642e726174696f8835302f32352f3235", blockNum: 161809416, data: &voteData{blockNum: 161809416, voter: v2, param: param{name: governance.GovernanceKeyMapReverse[params.Ratio], value: "50/25/25"}}},
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

func TestReadVoteDataFromDB(t *testing.T) {
	voteDatas := []voteData{
		{blockNum: 1, voter: common.Address{1}, param: param{name: governance.GovernanceKeyMapReverse[params.UnitPrice], value: uint64(100)}},
		{blockNum: 50, voter: common.Address{2}, param: param{name: governance.GovernanceKeyMapReverse[params.UnitPrice], value: uint64(200)}},
		{blockNum: 100, voter: common.Address{3}, param: param{name: governance.GovernanceKeyMapReverse[params.UnitPrice], value: uint64(300)}},
	}

	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)

	db := database.NewMemDB()
	voteDataBlockNums := make(gov.StoredVoteData, 0, len(voteDatas))
	for _, voteData := range voteDatas {
		headerVoteData, err := serializeVoteData(&voteData)
		chain.EXPECT().GetHeaderByNumber(uint64(voteData.blockNum)).Return(&types.Header{Vote: headerVoteData})
		require.NoError(t, err)
		voteDataBlockNums = append(voteDataBlockNums, voteData.blockNum)
	}
	gov.WriteVoteDataBlocks(db, &voteDataBlockNums)

	assert.Equal(t, voteDatas, readVoteDataFromDB(chain, db))
}
