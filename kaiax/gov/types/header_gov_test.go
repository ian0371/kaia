package types

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAllParamsHistoryMap(t *testing.T) {
	epoch := uint64(4)
	koreHf := uint64(10000)
	gov := []governanceData{
		0: {
			blockNum: 0,
			params: []param{
				{
					name:  "param1",
					value: 0,
				},
			},
		},
		4: {
			blockNum: 4,
			params: []param{
				{
					name:  "param1",
					value: 100,
				},
			},
		},
	}

	param1 := getAllParamsHistory(gov, epoch, koreHf)["param1"]
	assert.Equal(t, 0, param1.GetItem(0).value)
	assert.Equal(t, 0, param1.GetItem(8).value)
	assert.Equal(t, 100, param1.GetItem(9).value)
}

func TestEffectiveParams(t *testing.T) {
	koreHf := uint64(999999999)
	gov := map[uint64]governanceData{
		0: {
			blockNum: 0,
			params: []param{
				{
					name:  "param1",
					value: 0,
				},
			},
		},
		604800: {
			blockNum: 604800,
			params: []param{
				{
					name:  "param1",
					value: 100,
				},
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()
	h := NewHeaderGovernanceReader(chain, db, koreHf)
	for num, gov := range gov {
		assert.NoError(t, h.AddGov(num, gov))
	}

	assert.Equal(t, 0, h.EffectiveParams(0)["param1"].value)
	assert.Equal(t, 0, h.EffectiveParams(604800 * 2)["param1"].value)
	assert.Equal(t, 100, h.EffectiveParams(604800*2 + 1)["param1"].value)

	koreHf = 0
	db = database.NewMemDB()
	h = NewHeaderGovernanceReader(chain, db, koreHf)
	for num, gov := range gov {
		assert.NoError(t, h.AddGov(num, gov))
	}

	assert.Equal(t, 0, h.EffectiveParams(0)["param1"].value)
	assert.Equal(t, 0, h.EffectiveParams(604800*2 - 1)["param1"].value)
	assert.Equal(t, 100, h.EffectiveParams(604800 * 2)["param1"].value)
}
