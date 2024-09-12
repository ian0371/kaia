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
		defaultVal := uint64(250e9)
		hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil, nil)
		cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil, nil)
		ps, _ := m.EffectiveParamSet(1)
		assert.Equal(t, defaultVal, ps.UnitPrice)
	}

	// headergov value returned
	{
		headerGovVal := uint64(123)
		hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{"governance.unitprice": headerGovVal}, nil)
		cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{}, nil)
		ps, _ := m.EffectiveParamSet(1)
		assert.Equal(t, headerGovVal, ps.UnitPrice)
	}

	// contractgov value returned
	{
		contractGovVal := uint64(456)
		hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{"governance.unitprice": 0}, nil)
		cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{"governance.unitprice": contractGovVal}, nil)
		ps, _ := m.EffectiveParamSet(1)
		assert.Equal(t, contractGovVal, ps.UnitPrice)
	}
}
