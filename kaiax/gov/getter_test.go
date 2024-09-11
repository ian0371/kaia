package gov

import (
	"testing"

	contractgov_mock "github.com/kaiachain/kaia/kaiax/contractgov/mocks"
	headergov_mock "github.com/kaiachain/kaia/kaiax/headergov/mocks"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
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

	hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{"governance.unitprice": uint64(250)}, nil).AnyTimes()
	cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{}, nil)

	ps, _ := m.EffectiveParamSet(1)
	assert.Equal(t, uint64(250), ps.UnitPrice)

	cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[string]interface{}{"governance.unitprice": uint64(500)}, nil)

	ps, _ = m.EffectiveParamSet(1)
	assert.Equal(t, uint64(500), ps.UnitPrice)
}
