package headergov

import (
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	logger = log.NewModuleLogger(log.KaiaXGov)

	errZeroEpoch     = errors.New("epoch cannot be zero")
	errNoChainConfig = errors.New("ChainConfig or Istanbul is not set")
)

type VoteData = headergov_types.VoteData
type GovernanceData = headergov_types.GovernanceData
type GovernanceHistory = headergov_types.GovernanceHistory
type GovernanceParam = headergov_types.GovernanceParam
type GovernanceCache = headergov_types.GovernanceCache

type chain interface {
	GetHeaderByNumber(number uint64) *types.Header
	CurrentBlock() *types.Block
}

type InitOpts struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	NodeAddress common.Address
}

type HeaderGovModule struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	NodeAddress common.Address
	MyVotes     []VoteData // queue

	epoch uint64
	cache GovernanceCache
}

func NewHeaderGovModule() *HeaderGovModule {
	return &HeaderGovModule{}
}

func (h *HeaderGovModule) Init(opts *InitOpts) error {
	h.ChainKv = opts.ChainKv
	h.ChainConfig = opts.ChainConfig
	h.Chain = opts.Chain
	h.NodeAddress = opts.NodeAddress
	h.MyVotes = make([]VoteData, 0)
	if h.ChainConfig == nil || h.ChainConfig.Istanbul == nil {
		return errNoChainConfig
	}

	h.epoch = h.ChainConfig.Istanbul.Epoch
	if h.epoch == 0 {
		return errZeroEpoch
	}

	h.cache = GovernanceCache{
		Votes:       readVoteDataFromDB(h.Chain, h.ChainKv),
		Governances: readGovDataFromDB(h.Chain, h.ChainKv),
	}

	return nil
}

func (s *HeaderGovModule) Start() error {
	logger.Info("HeaderGovModule started")
	return nil
}

func (s *HeaderGovModule) Stop() {
	logger.Info("HeaderGovModule stopped")
}

func (s *HeaderGovModule) isKoreHF(num uint64) bool {
	return s.ChainConfig.IsKoreForkEnabled(new(big.Int).SetUint64(num))
}

func readVoteDataFromDB(chain chain, db database.Database) map[uint64]VoteData {
	voteBlocks := ReadVoteDataBlockNums(db)
	votes := make(map[uint64]VoteData)
	if voteBlocks != nil {
		for _, blockNum := range *voteBlocks {
			header := chain.GetHeaderByNumber(blockNum)
			parsedVote, err := headergov_types.DeserializeHeaderVote(header.Vote, blockNum)
			if err != nil {
				logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			}

			votes[blockNum] = *parsedVote
		}
	}

	return votes
}

func readGovDataFromDB(chain chain, db database.Database) map[uint64]GovernanceData {
	govBlocks := ReadGovDataBlockNums(db)
	govs := make(map[uint64]GovernanceData)
	if govBlocks != nil {
		for _, blockNum := range *govBlocks {
			header := chain.GetHeaderByNumber(blockNum)
			parsedGov, err := headergov_types.DeserializeHeaderGov(header.Governance, blockNum)
			if err != nil {
				logger.Error("Failed to parse vote", "num", blockNum, "err", err)
				panic(err)
			}

			govs[blockNum] = *parsedGov
		}
	}

	return govs
}
