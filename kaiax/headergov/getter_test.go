package headergov

import (
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
