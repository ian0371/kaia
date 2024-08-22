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

	testCases := []struct {
		desc          string
		koreBlock     *big.Int
		blockNum      uint64
		expectedPrice uint64
	}{
		{"Pre-Kore, Block 0", big.NewInt(999999999), 0, 25},
		{"Pre-Kore, Block 1209600", big.NewInt(999999999), 604800 * 2, 25},
		{"Pre-Kore, Block 1209601", big.NewInt(999999999), 604800*2 + 1, 750},
		{"Post-Kore, Block 0", big.NewInt(0), 0, 25},
		{"Post-Kore, Block 1209600", big.NewInt(0), 604800 * 2, 750},
		{"Post-Kore, Block 1209601", big.NewInt(0), 604800*2 + 1, 750},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			chain := mocks.NewMockBlockChain(mockCtrl)
			db := database.NewMemDB()
			config := &params.ChainConfig{
				KoreCompatibleBlock: tc.koreBlock,
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
				require.NoError(t, h.AddGov(&g))
			}

			pset, err := h.EffectiveParams(tc.blockNum)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedPrice, pset.UnitPrice())
		})
	}
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
