package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/stretchr/testify/assert"
)

var _ VoteData = (*voteData)(nil)

func TestNewVote(t *testing.T) {
	goodVotes := []struct {
		name  string
		value interface{}
	}{
		{name: "governance.addvalidator", value: common.Address{}},
		{name: "governance.deriveshaimpl", value: float64(0.0)},
		{name: "governance.deriveshaimpl", value: uint64(2)},
		{name: "governance.governingnode", value: "000000000000000000000000000abcd000000000"},
		{name: "governance.governingnode", value: "0x0000000000000000000000000000000000000000"},
		{name: "governance.governingnode", value: "0x000000000000000000000000000abcd000000000"},
		{name: "governance.governingnode", value: "0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456"},
		{name: "governance.governingnode", value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{name: "governance.governingnode", value: common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")},
		{name: "governance.govparamcontract", value: "000000000000000000000000000abcd000000000"},
		{name: "governance.govparamcontract", value: "0x0000000000000000000000000000000000000000"},
		{name: "governance.govparamcontract", value: "0x000000000000000000000000000abcd000000000"},
		{name: "governance.govparamcontract", value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{name: "governance.removevalidator", value: common.Address{}},
		{name: "governance.unitprice", value: float64(0.0)},
		{name: "governance.unitprice", value: float64(25e9)},
		{name: "governance.unitprice", value: uint64(25000000000)},
		{name: "governance.unitprice", value: uint64(25e9)},
		{name: "istanbul.committeesize", value: float64(7.0)},
		{name: "istanbul.committeesize", value: uint64(7)},
		{name: "kip71.basefeedenominator", value: uint64(64)},
		{name: "kip71.gastarget", value: uint64(15000000)},
		{name: "kip71.gastarget", value: uint64(30000000)},
		{name: "kip71.lowerboundbasefee", value: uint64(25000000000)},
		{name: "kip71.maxblockgasusedforbasefee", value: uint64(84000000)},
		{name: "kip71.upperboundbasefee", value: uint64(750000000000)},
		{name: "reward.kip82ratio", value: "10/90"},
		{name: "reward.kip82ratio", value: "20/80"},
		{name: "reward.mintingamount", value: "0"},
		{name: "reward.mintingamount", value: "9600000000000000000"},
		{name: "reward.mintingamount", value: new(big.Int).SetUint64(9.6e18)},
		{name: "reward.ratio", value: "0/0/100"},
		{name: "reward.ratio", value: "0/100/0"},
		{name: "reward.ratio", value: "10/10/80"},
		{name: "reward.ratio", value: "100/0/0"},
		{name: "reward.ratio", value: "30/40/30"},
		{name: "reward.ratio", value: "50/25/25"},
	}

	badVotes := []struct {
		name  string
		value interface{}
	}{
		// {name: "governance.addvalidator", value: common.Address{}},
		{name: "governance.deriveshaimpl", value: "2"},
		{name: "governance.deriveshaimpl", value: false},
		{name: "governance.deriveshaimpl", value: float64(-1)},
		{name: "governance.deriveshaimpl", value: float64(0.1)},
		{name: "governance.governancemode", value: "ballot"},
		{name: "governance.governancemode", value: "none"},
		{name: "governance.governancemode", value: "single"},
		{name: "governance.governancemode", value: "unexpected"},
		{name: "governance.governancemode", value: 0},
		{name: "governance.governancemode", value: 1},
		{name: "governance.governancemode", value: 2},
		{name: "governance.governingnode", value: "0x00000000000000000000"},
		{name: "governance.governingnode", value: "0x000000000000000000000000000xxxx000000000"},
		{name: "governance.governingnode", value: "address"},
		{name: "governance.governingnode", value: 0},
		{name: "governance.governingnode", value: []byte{0}},
		{name: "governance.governingnode", value: []byte{}},
		{name: "governance.governingnode", value: false},
		{name: "governance.govparamcontract", value: "0x00000000000000000000"},
		{name: "governance.govparamcontract", value: "0x000000000000000000000000000xxxx000000000"},
		{name: "governance.govparamcontract", value: "address"},
		{name: "governance.govparamcontract", value: 0},
		{name: "governance.govparamcontract", value: []byte{0}},
		{name: "governance.govparamcontract", value: []byte{}},
		{name: "governance.govparamcontract", value: false},
		// {name: "governance.removevalidator", value: common.Address{}},
		{name: "governance.unitprice", value: "25000000000"},
		{name: "governance.unitprice", value: false},
		{name: "governance.unitprice", value: float64(-10)},
		{name: "governance.unitprice", value: float64(0.1)},
		{name: "istanbul.Epoch", value: float64(30000.10)},
		{name: "istanbul.committeesize", value: "7"},
		{name: "istanbul.committeesize", value: false},
		{name: "istanbul.committeesize", value: float64(-7)},
		{name: "istanbul.committeesize", value: float64(7.1)},
		{name: "istanbul.committeesize", value: uint64(0)},
		{name: "istanbul.epoch", value: "bad"},
		{name: "istanbul.epoch", value: false},
		{name: "istanbul.epoch", value: float64(30000.00)},
		{name: "istanbul.epoch", value: uint64(30000)},
		{name: "istanbul.policy", value: "RoundRobin"},
		{name: "istanbul.policy", value: "WeightedRandom"},
		{name: "istanbul.policy", value: "roundrobin"},
		{name: "istanbul.policy", value: "sticky"},
		{name: "istanbul.policy", value: "weightedrandom"},
		{name: "istanbul.policy", value: false},
		{name: "istanbul.policy", value: float64(1.0)},
		{name: "istanbul.policy", value: float64(1.2)},
		{name: "istanbul.policy", value: uint64(0)},
		{name: "istanbul.policy", value: uint64(1)},
		{name: "istanbul.policy", value: uint64(2)},
		{name: "kip71.basefeedenominator", value: "64"},
		{name: "kip71.basefeedenominator", value: "sixtyfour"},
		{name: "kip71.basefeedenominator", value: 64},
		{name: "kip71.basefeedenominator", value: false},
		{name: "kip71.gastarget", value: "30000"},
		{name: "kip71.gastarget", value: 3000},
		{name: "kip71.gastarget", value: false},
		{name: "kip71.gastarget", value: true},
		{name: "kip71.lowerboundbasefee", value: "250000000"},
		{name: "kip71.lowerboundbasefee", value: "test"},
		{name: "kip71.lowerboundbasefee", value: 25000000},
		{name: "kip71.lowerboundbasefee", value: false},
		{name: "kip71.maxblockgasusedforbasefee", value: "84000"},
		{name: "kip71.maxblockgasusedforbasefee", value: 0},
		{name: "kip71.maxblockgasusedforbasefee", value: 840000},
		{name: "kip71.maxblockgasusedforbasefee", value: false},
		{name: "kip71.upperboundbasefee", value: "750000"},
		{name: "kip71.upperboundbasefee", value: 7500000},
		{name: "kip71.upperboundbasefee", value: false},
		{name: "kip71.upperboundbasefee", value: true},
		{name: "reward.deferredtxfee", value: "false"},
		{name: "reward.deferredtxfee", value: 0},
		{name: "reward.deferredtxfee", value: 1},
		{name: "reward.deferredtxfee", value: false},
		{name: "reward.deferredtxfee", value: true},
		{name: "reward.kip82ratio", value: "30/30/40"},
		{name: "reward.kip82ratio", value: "30/80"},
		{name: "reward.kip82ratio", value: "49.5/50.5"},
		{name: "reward.kip82ratio", value: "50.5/50.5"},
		{name: "reward.minimumstake", value: "-1"},
		{name: "reward.minimumstake", value: "0"},
		{name: "reward.minimumstake", value: "2000000000000000000000000"},
		{name: "reward.minimumstake", value: 0},
		{name: "reward.minimumstake", value: 1.1},
		{name: "reward.minimumstake", value: 200000000000000},
		{name: "reward.mintingamount", value: "many"},
		{name: "reward.mintingamount", value: 96000},
		{name: "reward.mintingamount", value: false},
		{name: "reward.proposerupdateinterval", value: "20"},
		{name: "reward.proposerupdateinterval", value: float64(20.0)},
		{name: "reward.proposerupdateinterval", value: float64(20.2)},
		{name: "reward.proposerupdateinterval", value: uint64(20)},
		{name: "reward.ratio", value: "0/0/0"},
		{name: "reward.ratio", value: "30.5/40/29.5"},
		{name: "reward.ratio", value: "30.5/40/30.5"},
		{name: "reward.ratio", value: "30/40/29"},
		{name: "reward.ratio", value: "30/40/31"},
		{name: "reward.ratio", value: "30/70"},
		{name: "reward.ratio", value: 30 / 40 / 30},
		{name: "reward.stakingupdateinterval", value: "20"},
		{name: "reward.stakingupdateinterval", value: float64(20.0)},
		{name: "reward.stakingupdateinterval", value: float64(20.2)},
		{name: "reward.stakingupdateinterval", value: uint64(20)},
		{name: "reward.useginicoeff", value: "false"},
		{name: "reward.useginicoeff", value: 0},
		{name: "reward.useginicoeff", value: 1},
		{name: "reward.useginicoeff", value: false},
		{name: "reward.useginicoeff", value: true},
	}

	for _, tc := range goodVotes {
		t.Run("goodVote/"+tc.name, func(t *testing.T) {
			assert.NotNil(t, NewVoteData(common.Address{}, tc.name, tc.value))
		})
	}

	for _, tc := range badVotes {
		t.Run("badVote/"+tc.name, func(t *testing.T) {
			assert.Nil(t, NewVoteData(common.Address{}, tc.name, tc.value))
		})
	}
}

func TestHeaderVoteSerialization(t *testing.T) {
	v1 := common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541")
	v2 := common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")

	tcs := []struct {
		serializedVoteData string
		blockNum           uint64
		voteData           VoteData
	}{
		///// all vote datas.
		{serializedVoteData: "0xf8439452d41ca72af615a1ac3301b0a93efa222ecc754198676f7665726e616e63652e676f7665726e696e676e6f64659452d41ca72af615a1ac3301b0a93efa222ecc7541", blockNum: 1, voteData: NewVoteData(v1, "governance.governingnode", v1)},
		{serializedVoteData: "0xed9452d41ca72af615a1ac3301b0a93efa222ecc7541917265776172642e6b69703832726174696f8533332f3637", blockNum: 2, voteData: NewVoteData(v1, "reward.kip82ratio", "33/67")},
		{serializedVoteData: "0xf39452d41ca72af615a1ac3301b0a93efa222ecc7541976b697037312e6c6f776572626f756e64626173656665658505d21dba00", blockNum: 3, voteData: NewVoteData(v1, "kip71.lowerboundbasefee", uint64(25e9))},
		{serializedVoteData: "0xf39452d41ca72af615a1ac3301b0a93efa222ecc7541976b697037312e7570706572626f756e646261736566656585ae9f7bcc00", blockNum: 4, voteData: NewVoteData(v1, "kip71.upperboundbasefee", uint64(750e9))},
		{serializedVoteData: "0xef9452d41ca72af615a1ac3301b0a93efa222ecc7541986b697037312e6261736566656564656e6f6d696e61746f7264", blockNum: 5, voteData: NewVoteData(v1, "kip71.basefeedenominator", uint64(100))},
		{serializedVoteData: "0xf83e9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e749331303030303030303030303030303030303030", blockNum: 6, voteData: NewVoteData(v1, "reward.mintingamount", big.NewInt(1000000000000000000))},
		// TODO: add govparamcontract from baobab

		///// Real mainnet vote data.
		{serializedVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e6974707269636585ae9f7bcc00", blockNum: 86119166, voteData: NewVoteData(v1, "governance.unitprice", uint64(750000000000))},
		{serializedVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e69747072696365853a35294400", blockNum: 90355962, voteData: NewVoteData(v1, "governance.unitprice", uint64(250000000000))},
		{serializedVoteData: "0xed9452d41ca72af615a1ac3301b0a93efa222ecc754196697374616e62756c2e636f6d6d697474656573697a651f", blockNum: 95352567, voteData: NewVoteData(v1, "istanbul.committeesize", uint64(31))},
		{serializedVoteData: "0xf83e9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e749336343030303030303030303030303030303030", blockNum: 105629058, voteData: NewVoteData(v1, "reward.mintingamount", big.NewInt(6400000000000000000))},
		{serializedVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f34302f3130", blockNum: 105629111, voteData: NewVoteData(v1, "reward.ratio", "50/40/10")},
		{serializedVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f32302f3330", blockNum: 118753908, voteData: NewVoteData(v1, "reward.ratio", "50/20/30")},
		{serializedVoteData: "0xf8439452d41ca72af615a1ac3301b0a93efa222ecc754198676f7665726e616e63652e676f7665726e696e676e6f646594c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456", blockNum: 126061533, voteData: NewVoteData(v1, "governance.governingnode", v2)},
		{serializedVoteData: "0xef94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45698676f7665726e616e63652e646572697665736861696d706c80", blockNum: 127692621, voteData: NewVoteData(v2, "governance.deriveshaimpl", uint64(0))},
		{serializedVoteData: "0xe994c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568f6b697037312e67617374617267657483e4e1c0", blockNum: 140916059, voteData: NewVoteData(v2, "kip71.gastarget", uint64(15000000))},
		{serializedVoteData: "0xf83a94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4569f6b697037312e6d6178626c6f636b67617375736564666f72626173656665658401c9c380", blockNum: 140916152, voteData: NewVoteData(v2, "kip71.maxblockgasusedforbasefee", uint64(30000000))},
		{serializedVoteData: "0xed94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45696697374616e62756c2e636f6d6d697474656573697a6532", blockNum: 161809335, voteData: NewVoteData(v2, "istanbul.committeesize", uint64(50))},
		{serializedVoteData: "0xf83e94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456947265776172642e6d696e74696e67616d6f756e749339363030303030303030303030303030303030", blockNum: 161809370, voteData: NewVoteData(v2, "reward.mintingamount", new(big.Int).SetUint64(9.6e18))},
		{serializedVoteData: "0xeb94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568c7265776172642e726174696f8835302f32352f3235", blockNum: 161809416, voteData: NewVoteData(v2, "reward.ratio", "50/25/25")},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("TestCase_block_%d", tc.blockNum), func(t *testing.T) {
			// Test deserialization
			actual, err := DeserializeHeaderVote(hexutil.MustDecode(tc.serializedVoteData), tc.blockNum)
			assert.NoError(t, err)
			assert.Equal(t, tc.voteData, actual, "DeserializeHeaderVote() failed")

			// Test serialization
			serialized, err := tc.voteData.Serialize()
			assert.NoError(t, err)
			assert.Equal(t, tc.serializedVoteData, hexutil.Encode(serialized), "voteData.Serialize() failed")
		})
	}
}
