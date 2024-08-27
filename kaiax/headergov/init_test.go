package headergov

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadVoteBlockNumsFromDB(t *testing.T) {
	voteDatas := []VoteData{
		{BlockNum: 1, Voter: common.Address{1}, Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(100)},
		{BlockNum: 50, Voter: common.Address{2}, Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(200)},
		{BlockNum: 100, Voter: common.Address{3}, Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(300)},
	}

	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)

	db := database.NewMemDB()
	voteDataBlockNums := make(StoredVoteBlockNums, 0, len(voteDatas))
	for _, voteData := range voteDatas {
		headerVoteData, err := voteData.Serialize()
		require.NoError(t, err)
		chain.EXPECT().GetHeaderByNumber(uint64(voteData.BlockNum)).Return(&types.Header{Vote: headerVoteData})
		voteDataBlockNums = append(voteDataBlockNums, voteData.BlockNum)
	}
	WriteVoteDataBlockNums(db, &voteDataBlockNums)

	assert.Equal(t, voteDatas, readVoteDataFromDB(chain, db))
}

func TestReadGovDataFromDB(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()

	ps1 := &GovernanceParam{
		UnitPrice: uint64(100),
	}
	ps2 := &GovernanceParam{
		UnitPrice: uint64(200),
	}

	WriteGovernanceParam(db, 1, ps1)
	WriteGovernanceParam(db, 2, ps2)
	WriteGovDataBlockNums(db, &StoredGovBlockNums{1, 2})

	govs := []GovernanceData{
		{BlockNum: 1, Params: map[string]interface{}{governance.GovernanceKeyMapReverse[params.UnitPrice]: ps1.UnitPrice}},
		{BlockNum: 2, Params: map[string]interface{}{governance.GovernanceKeyMapReverse[params.UnitPrice]: ps2.UnitPrice}},
	}
	for _, govData := range govs {
		headerGovData, err := govData.Serialize()
		require.NoError(t, err)
		chain.EXPECT().GetHeaderByNumber(uint64(govData.BlockNum)).Return(&types.Header{Governance: headerGovData})
	}

	assert.Equal(t, ps1, ReadGovernanceParam(db, 1))
	assert.Equal(t, ps2, ReadGovernanceParam(db, 2))
	assert.Equal(t, govs, readGovDataFromDB(chain, db))
}
