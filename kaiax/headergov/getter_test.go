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
	gasPrice := governance.GovernanceKeyMapReverse[params.UnitPrice]
	gov := []GovernanceData{
		{
			BlockNum: 0,
			Params: map[string]interface{}{
				gasPrice: uint64(25),
			},
		},
		{
			BlockNum: 604800,
			Params: map[string]interface{}{
				gasPrice: uint64(750),
			},
		},
	}

	testCases := []struct {
		desc          string
		koreBlock     uint64
		blockNum      uint64
		expectedPrice uint64
	}{
		{"Pre-Kore, Block 0", 999999999, 0, 25},
		{"Pre-Kore, Block 1209600", 999999999, 1209600, 25},
		{"Pre-Kore, Block 1209601", 999999999, 1209601, 750},
		{"Post-Kore, Block 0", 0, 0, 25},
		{"Post-Kore, Block 1209600", 0, 1209600, 750},
		{"Post-Kore, Block 1209601", 0, 1209601, 750},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			chain := mocks.NewMockBlockChain(mockCtrl)
			db := database.NewMemDB()
			config := &params.ChainConfig{
				KoreCompatibleBlock: big.NewInt(int64(tc.koreBlock)),
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

			for _, g := range gov {
				h.HandleGov(&g)
			}

			gp, err := h.EffectiveParams(tc.blockNum)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedPrice, gp.UnitPrice)
		})
	}
}

func TestSourceBlockNum(t *testing.T) {
	epoch := uint64(1000)
	testCases := []struct {
		blockNum    uint64
		isKore      bool
		expectedGov uint64
	}{
		{0, false, 0},
		{999, false, 0},
		{1000, false, 0},
		{1001, false, 0},
		{1999, false, 0},
		{2000, false, 0},
		{2001, false, 1000},
		{2999, false, 1000},
		{3000, false, 1000},
		{3001, false, 2000},

		{0, true, 0},
		{999, true, 0},
		{1000, true, 0},
		{1001, true, 0},
		{1999, true, 0},
		{2000, true, 1000},
		{2001, true, 1000},
		{2999, true, 1000},
		{3000, true, 2000},
		{3001, true, 2000},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Block %d", tc.blockNum), func(t *testing.T) {
			result := PrevEpochStart(tc.blockNum, epoch, tc.isKore)
			assert.Equal(t, tc.expectedGov, result, "Incorrect governance data block for block %d", tc.blockNum)
		})
	}
}
