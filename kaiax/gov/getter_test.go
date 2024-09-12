package gov

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	contractgov_mock "github.com/kaiachain/kaia/kaiax/contractgov/mocks"
	headergov_mock "github.com/kaiachain/kaia/kaiax/headergov/mocks"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func newHeaderGovModuleMock(t *testing.T, config *params.ChainConfig) *headergov_mock.MockHeaderGovModule {
	mock := headergov_mock.NewMockHeaderGovModule(gomock.NewController(t))
	return mock
}

func newContractGovModuleMock(t *testing.T, config *params.ChainConfig) *contractgov_mock.MockContractGovModule {
	mock := contractgov_mock.NewMockContractGovModule(gomock.NewController(t))
	return mock
}

func TestEffectiveParamSet(t *testing.T) {
	config := params.TestChainConfig
	hgm := newHeaderGovModuleMock(t, config)
	cgm := newContractGovModuleMock(t, config)
	m := &govModule{
		hgm: hgm,
		cgm: cgm,
	}

	// default value returned
	{
		hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil, nil)
		cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil, nil)
		ps, _ := m.EffectiveParamSet(1)
		assert.Equal(t, uint64(250e9), ps.UnitPrice)
	}

	// headergov value returned
	{
		val := uint64(123)
		hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{"governance.unitprice": val}, nil)
		cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{}, nil)
		ps, _ := m.EffectiveParamSet(1)
		assert.Equal(t, val, ps.UnitPrice)
	}

	// contractgov value returned
	{
		val := uint64(456)
		hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{"governance.unitprice": 0}, nil)
		cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{"governance.unitprice": val}, nil)
		ps, _ := m.EffectiveParamSet(1)
		assert.Equal(t, val, ps.UnitPrice)
	}
}
