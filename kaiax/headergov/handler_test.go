package headergov

import (
	"math/big"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderVerify(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()
	config := &params.ChainConfig{
		KoreCompatibleBlock: big.NewInt(999999999),
		Istanbul: &params.IstanbulConfig{
			Epoch: 604800,
		},
	}
	h := &HeaderGovModule{}
	err := h.Init(&InitOpts{
		Chain:       chain,
		ChainKv:     db,
		ChainConfig: config,
	})
	require.NoError(t, err)
	h.AddVote(&VoteData{
		BlockNum: 1,
		Name:     governance.GovernanceKeyMapReverse[params.UnitPrice],
		Value:    uint64(100),
	})

	gov := GovernanceData{
		BlockNum: 604800,
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
		{604799, nil, false},
		{604799, govBytes, true},

		{604800, nil, true},
		{604800, govBytes, false},

		{604801, nil, false},
		{604801, govBytes, true},
	}

	for _, tc := range tcs {
		err = h.VerifyHeader(&types.Header{
			Number:     big.NewInt(int64(tc.blockNum)),
			Governance: tc.gov,
		})
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
