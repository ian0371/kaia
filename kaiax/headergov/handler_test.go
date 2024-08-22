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
		Param: Param{
			Name:  governance.GovernanceKeyMapReverse[params.UnitPrice],
			Value: uint64(100),
		},
	})

	gov := GovernanceData{
		BlockNum: 604800,
		Params: map[string]Param{
			governance.GovernanceKeyMapReverse[params.UnitPrice]: {
				Name:  governance.GovernanceKeyMapReverse[params.UnitPrice],
				Value: uint64(100),
			},
		},
	}
	govBytes, err := gov.Serialize()
	require.NoError(t, err)

	err = h.VerifyHeader(&types.Header{
		Number:     big.NewInt(604799),
		Governance: govBytes,
	})
	assert.Error(t, err)

	err = h.VerifyHeader(&types.Header{
		Number:     big.NewInt(604800),
		Governance: govBytes,
	})
	assert.NoError(t, err)

	err = h.VerifyHeader(&types.Header{
		Number:     big.NewInt(604801),
		Governance: govBytes,
	})
	assert.Error(t, err)
}
