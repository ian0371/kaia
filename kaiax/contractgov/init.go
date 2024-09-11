package contractgov

import (
	"errors"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	logger = log.NewModuleLogger(log.KaiaXGov)

	errNoChainConfig          = errors.New("ChainConfig or Istanbul is not set")
	errContractEngineNotReady = errors.New("ContractEngine is not ready")
	errParamsAtFail           = errors.New("headerGov EffectiveParams() failed")
	errGovParamNotExist       = errors.New("GovParam does not exist")
	errInvalidGovParam        = errors.New("GovParam conversion failed")
)

type ParamSet = headergov_types.ParamSet
type HeaderGovModule = headergov_types.HeaderGovModule

var Params = headergov_types.Params

type chain interface {
	blockchain.ChainContext

	GetHeaderByNumber(number uint64) *types.Header
	CurrentBlock() *types.Block
	State() (*state.StateDB, error)
	StateAt(root common.Hash) (*state.StateDB, error)
	Config() *params.ChainConfig
	GetBlock(hash common.Hash, number uint64) *types.Block
}

type InitOpts struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	hgm         HeaderGovModule
}

type contractGovModule struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	hgm         HeaderGovModule
}

func (c *contractGovModule) Init(opts *InitOpts) error {
	c.ChainKv = opts.ChainKv
	c.ChainConfig = opts.ChainConfig
	c.Chain = opts.Chain
	c.hgm = opts.hgm
	if c.ChainConfig == nil || c.ChainConfig.Istanbul == nil {
		return errNoChainConfig
	}

	return nil
}
