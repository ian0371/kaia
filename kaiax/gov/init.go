package gov

import (
	contractgov_types "github.com/kaiachain/kaia/kaiax/contractgov/types"
	gov_types "github.com/kaiachain/kaia/kaiax/gov/types"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/log"
)

var (
	_ GovModule = (*govModule)(nil)

	logger = log.NewModuleLogger(log.KaiaXGov)
)

type HeaderGovModule = headergov_types.HeaderGovModule
type ContractGovModule = contractgov_types.ContractGovModule
type ParamSet = headergov_types.ParamSet
type GovModule = gov_types.GovModule

type govModule struct {
	hgm HeaderGovModule
	cgm ContractGovModule
}

type InitOpts struct {
	hgm HeaderGovModule
	cgm ContractGovModule
}

func (m *govModule) Init(opts *InitOpts) error {
	m.hgm = opts.hgm
	m.cgm = opts.cgm
	return nil
}

func (m *govModule) Start() error {
	logger.Info("GovModule started")
	return nil
}

func (m *govModule) Stop() {
	logger.Info("GovModule stopped")
}
