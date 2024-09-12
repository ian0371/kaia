package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamSet(t *testing.T) {
	p := ParamSet{}
	err := p.Set("governance.governancemode", "none")
	assert.NoError(t, err)
	assert.Equal(t, p.GovernanceMode, "none")
}

func TestGetDefaultGovernanceParamSet(t *testing.T) {
	p := GetDefaultGovernanceParamSet()
	assert.NotNil(t, p)
	for name, param := range Params {
		fieldName := param.ParamSetFieldName
		fieldValue := reflect.ValueOf(p).Elem().FieldByName(fieldName).Interface()
		assert.Equal(t, param.DefaultValue, fieldValue, "Mismatch for %s", name)
	}
}
