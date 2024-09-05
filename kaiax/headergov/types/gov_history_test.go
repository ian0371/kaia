package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHistory(t *testing.T) {
	govs := map[uint64]GovData{
		0: {
			Params: map[string]interface{}{
				"governance.unitprice": uint64(100),
			},
		},
		4: {
			Params: map[string]interface{}{
				"governance.unitprice": uint64(200),
			},
		},
	}

	history := GetHistory(govs)
	assert.Equal(t, uint64(100), history[0].UnitPrice)
	assert.Equal(t, uint64(200), history[4].UnitPrice)
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		blockNumber       uint64
		expectedUnitPrice uint64
	}{
		{0, 100},
		{3, 100},
		{4, 200},
		{5, 200},
	}

	govs := map[uint64]GovData{
		0: {
			Params: map[string]interface{}{
				"governance.unitprice": uint64(100),
			},
		},
		4: {
			Params: map[string]interface{}{
				"governance.unitprice": uint64(200),
			},
		},
	}

	gh := GetHistory(govs)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Block %d", tc.blockNumber), func(t *testing.T) {
			gp, err := gh.Search(tc.blockNumber)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedUnitPrice, gp.UnitPrice)
		})
	}
}
