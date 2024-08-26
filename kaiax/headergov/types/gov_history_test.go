package types

import (
	"testing"

	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestGetGovHistory(t *testing.T) {
	govs := []GovernanceData{
		0: {
			BlockNum: 0,
			Params: map[string]interface{}{
				governance.GovernanceKeyMapReverse[params.UnitPrice]: uint64(100),
			},
		},
		4: {
			BlockNum: 4,
			Params: map[string]interface{}{
				governance.GovernanceKeyMapReverse[params.UnitPrice]: uint64(200),
			},
		},
	}

	gh := GetGovernanceHistory(govs)

	assert.Equal(t, GovernanceParam{UnitPrice: uint64(100)}, gh[0])
	assert.Equal(t, GovernanceParam{UnitPrice: uint64(200)}, gh[4])
}

func TestSearch(t *testing.T) {
	govs := []GovernanceData{
		0: {
			BlockNum: 0,
			Params: map[string]interface{}{
				governance.GovernanceKeyMapReverse[params.UnitPrice]: uint64(100),
			},
		},
		4: {
			BlockNum: 4,
			Params: map[string]interface{}{
				governance.GovernanceKeyMapReverse[params.UnitPrice]: uint64(200),
			},
		},
	}

	gh := GetGovernanceHistory(govs)
	gp, err := gh.Search(0)
	assert.Nil(t, err)
	assert.Equal(t, GovernanceParam{UnitPrice: uint64(100)}, gp)

	gp, err = gh.Search(3)
	assert.Nil(t, err)
	assert.Equal(t, GovernanceParam{UnitPrice: uint64(100)}, gp)

	gp, err = gh.Search(4)
	assert.Nil(t, err)
	assert.Equal(t, GovernanceParam{UnitPrice: uint64(200)}, gp)
}
