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

func TestReadGovVoteBlockNumsFromDB(t *testing.T) {
	votes := map[uint64]VoteData{
		1:   {Voter: common.Address{1}, Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(100)},
		50:  {Voter: common.Address{2}, Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(200)},
		100: {Voter: common.Address{3}, Name: governance.GovernanceKeyMapReverse[params.UnitPrice], Value: uint64(300)},
	}

	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)

	db := database.NewMemDB()
	voteDataBlockNums := make(StoredUint64Array, 0, len(votes))
	for num, voteData := range votes {
		headerVoteData, err := voteData.Serialize()
		require.NoError(t, err)
		chain.EXPECT().GetHeaderByNumber(num).Return(&types.Header{Vote: headerVoteData})
		voteDataBlockNums = append(voteDataBlockNums, num)
	}
	WriteGovVoteDataBlockNums(db, &voteDataBlockNums)

	assert.Equal(t, votes, readGovVoteDataFromDB(chain, db))
}

func TestReadGovDataFromDB(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()

	ps1 := &ParamSet{
		UnitPrice: uint64(100),
	}
	ps2 := &ParamSet{
		UnitPrice: uint64(200),
	}

	WriteGovDataBlockNums(db, &StoredUint64Array{1, 2})

	govs := map[uint64]GovData{
		1: {Params: map[string]interface{}{governance.GovernanceKeyMapReverse[params.UnitPrice]: ps1.UnitPrice}},
		2: {Params: map[string]interface{}{governance.GovernanceKeyMapReverse[params.UnitPrice]: ps2.UnitPrice}},
	}
	for num, govData := range govs {
		headerGovData, err := govData.Serialize()
		require.NoError(t, err)
		chain.EXPECT().GetHeaderByNumber(num).Return(&types.Header{Governance: headerGovData})
	}

	assert.Equal(t, govs, readGovDataFromDB(chain, db))
}
