package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParam(t *testing.T) {
	for name, param := range Params {
		assert.NotEmpty(t, param.ParamSetFieldName, name)
		assert.NotEmpty(t, param.Canonicalizer, name)
		// assert.NotEmpty(t, param.FormatChecker, name)
	}
}
