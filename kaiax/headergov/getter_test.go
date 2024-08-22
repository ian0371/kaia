package headergov

import (
	"fmt"
	"math/big"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEffectiveParams(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	paramName := governance.GovernanceKeyMapReverse[params.UnitPrice]
	gov := []GovernanceData{
		{
			BlockNum: 0,
			Params: map[string]Param{
				paramName: {
					Name:  paramName,
					Value: uint64(25),
				},
			},
		},
		{
			BlockNum: 604800,
			Params: map[string]Param{
				paramName: {
					Name:  paramName,
					Value: uint64(750),
				},
			},
		},
	}

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

	for _, gov := range gov {
		require.NoError(t, h.AddGov(&gov))
	}

	pset, err := h.EffectiveParams(0)
	require.NoError(t, err)
	assert.Equal(t, uint64(25), pset.UnitPrice())

	pset, err = h.EffectiveParams(604800 * 2)
	require.NoError(t, err)
	assert.Equal(t, uint64(25), pset.UnitPrice())

	pset, err = h.EffectiveParams(604800*2 + 1)
	require.NoError(t, err)
	assert.Equal(t, uint64(750), pset.UnitPrice())

	config.KoreCompatibleBlock = big.NewInt(0)
	db = database.NewMemDB()
	h = &HeaderGovModule{}
	err = h.Init(&InitOpts{
		Chain:       chain,
		ChainKv:     db,
		ChainConfig: config,
	})
	require.NoError(t, err)
	for _, gov := range gov {
		assert.NoError(t, h.AddGov(&gov))
	}

	pset, err = h.EffectiveParams(0)
	require.NoError(t, err)
	assert.Equal(t, uint64(25), pset.UnitPrice())

	pset, err = h.EffectiveParams(604800 * 2)
	require.NoError(t, err)
	assert.Equal(t, uint64(750), pset.UnitPrice())

	pset, err = h.EffectiveParams(604800*2 + 1)
	require.NoError(t, err)
	assert.Equal(t, uint64(750), pset.UnitPrice())
}

func TestCalcGovDataBlock(t *testing.T) {
	epoch := uint64(604800)
	koreHF := epoch * 3

	testCases := []struct {
		blockNum    uint64
		expectedGov uint64
	}{
		{0, 0},
		{epoch - 1, 0},
		{epoch, 0},
		{epoch + 1, 0},
		{epoch*2 - 1, 0},
		{epoch * 2, 0},
		{epoch*2 + 1, epoch},
		{epoch*3 - 1, epoch},
		{epoch * 3, epoch * 2},
		{epoch*3 + 1, epoch * 2},
		{epoch*4 - 1, epoch * 2},
		{epoch * 4, epoch * 3},
		{epoch*4 + 1, epoch * 3},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Block %d", tc.blockNum), func(t *testing.T) {
			result := CalcGovDataBlock(tc.blockNum, epoch, tc.blockNum >= koreHF)
			assert.Equal(t, tc.expectedGov, result, "Incorrect governance data block for block %d", tc.blockNum)
		})
	}
}
