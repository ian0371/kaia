package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGovParamSet(t *testing.T) {
	p := GovParamSet{}
	err := p.Set("governance.governancemode", "none")
	assert.NoError(t, err)
	assert.Equal(t, p.GovernanceMode, "none")
}
