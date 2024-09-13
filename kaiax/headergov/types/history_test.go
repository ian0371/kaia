package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHistory(t *testing.T) {
	govs := map[uint64]GovData{
		0: NewGovData(map[string]interface{}{
			Params[GovernanceUnitPrice].Name: uint64(100),
		}),
		4: NewGovData(map[string]interface{}{
			Params[GovernanceUnitPrice].Name: uint64(200),
		}),
	}

	history := GetHistory(govs)
	assert.Equal(t, uint64(100), history[0].UnitPrice)
	assert.Equal(t, uint64(200), history[4].UnitPrice)
}

func TestSearch(t *testing.T) {
	govs := map[uint64]GovData{
		0: NewGovData(map[string]interface{}{
			Params[GovernanceUnitPrice].Name: uint64(100),
		}),
		4: NewGovData(map[string]interface{}{
			Params[GovernanceUnitPrice].Name: uint64(200),
		}),
	}

	gh := GetHistory(govs)
	for i := 0; i < 4; i++ {
		t.Run(fmt.Sprintf("Block %d", i), func(t *testing.T) {
			gp, err := gh.Search(uint64(i))
			assert.Nil(t, err)
			assert.Equal(t, uint64(100), gp.UnitPrice)
		})
	}

	for i := 4; i < 100; i++ {
		t.Run(fmt.Sprintf("Block %d", i), func(t *testing.T) {
			gp, err := gh.Search(uint64(i))
			assert.Nil(t, err)
			assert.Equal(t, uint64(200), gp.UnitPrice)
		})
	}
}
