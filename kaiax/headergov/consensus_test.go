package headergov

import (
	"fmt"
	"math/big"
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

func TestVerifyHeader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()
	config := &params.ChainConfig{
		KoreCompatibleBlock: big.NewInt(999999999),
		Istanbul: &params.IstanbulConfig{
			Epoch: 1000,
		},
	}
	h := &HeaderGovModule{}
	err := h.Init(&InitOpts{
		Chain:       chain,
		ChainKv:     db,
		ChainConfig: config,
	})
	require.NoError(t, err)
	h.HandleVote(500, NewVoteData(common.Address{1}, governance.GovernanceKeyMapReverse[params.UnitPrice], uint64(100)))

	gov := GovData{
		Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.UnitPrice]: uint64(100),
		},
	}
	govBytes, err := gov.Serialize()
	require.NoError(t, err)

	tcs := []struct {
		blockNum uint64
		gov      []byte
		isError  bool
	}{
		{999, nil, false},
		{999, govBytes, true},

		{1000, nil, true},
		{1000, govBytes, false},

		{1001, nil, false},
		{1001, govBytes, true},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("BlockNum=%d,HasGov=%v", tc.blockNum, tc.gov != nil), func(t *testing.T) {
			err := h.VerifyHeader(&types.Header{
				Number:     big.NewInt(int64(tc.blockNum)),
				Governance: tc.gov,
			})
			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetVotesInEpoch(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()
	config := &params.ChainConfig{
		Istanbul: &params.IstanbulConfig{
			Epoch: 1000,
		},
	}
	h := &HeaderGovModule{}
	err := h.Init(&InitOpts{
		Chain:       chain,
		ChainKv:     db,
		ChainConfig: config,
	})
	require.NoError(t, err)
	v1 := NewVoteData(common.Address{1}, governance.GovernanceKeyMapReverse[params.UnitPrice], uint64(100))
	h.HandleVote(500, v1)
	v2 := NewVoteData(common.Address{2}, governance.GovernanceKeyMapReverse[params.UnitPrice], uint64(200))
	h.HandleVote(1500, v2)

	assert.Equal(t, []VoteData{*v1}, h.getVotesInEpoch(0))
	assert.Equal(t, []VoteData{*v2}, h.getVotesInEpoch(1))
}

func TestGetExpectedGovernance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()
	config := &params.ChainConfig{
		KoreCompatibleBlock: big.NewInt(999999999),
		Istanbul: &params.IstanbulConfig{
			Epoch: 1000,
		},
	}
	h := &HeaderGovModule{}
	err := h.Init(&InitOpts{
		Chain:       chain,
		ChainKv:     db,
		ChainConfig: config,
	})
	require.NoError(t, err)

	v1 := NewVoteData(common.Address{1}, governance.GovernanceKeyMapReverse[params.UnitPrice], uint64(100))
	h.HandleVote(500, v1)
	v2 := NewVoteData(common.Address{2}, governance.GovernanceKeyMapReverse[params.UnitPrice], uint64(200))
	h.HandleVote(1500, v2)

	g1 := &GovData{
		Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.UnitPrice]: uint64(100),
		},
	}
	h.HandleGov(1000, g1)
	g2 := &GovData{
		Params: map[string]interface{}{
			governance.GovernanceKeyMapReverse[params.UnitPrice]: uint64(200),
		},
	}
	h.HandleGov(2000, g2)
	assert.Equal(t, *g1, h.getExpectedGovernance(1000))
	assert.Equal(t, *g2, h.getExpectedGovernance(2000))
}
