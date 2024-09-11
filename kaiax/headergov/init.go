package headergov

import (
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ HeaderGovModule = (*headerGovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaXGov)

	errZeroEpoch     = errors.New("epoch cannot be zero")
	errNoChainConfig = errors.New("ChainConfig or Istanbul is not set")
)

type VoteData = headergov_types.VoteData
type VotesInEpoch = headergov_types.VotesInEpoch
type GovData = headergov_types.GovData
type History = headergov_types.History
type ParamSet = headergov_types.ParamSet
type HeaderCache = headergov_types.HeaderCache
type HeaderGovModule = headergov_types.HeaderGovModule

var NewVoteData = headergov_types.NewVoteData
var NewGovData = headergov_types.NewGovData
var Params = headergov_types.Params
var DeserializeHeaderVote = headergov_types.DeserializeHeaderVote
var DeserializeHeaderGov = headergov_types.DeserializeHeaderGov

type chain interface {
	GetHeaderByNumber(number uint64) *types.Header
	CurrentBlock() *types.Block
	State() (*state.StateDB, error)
}

type InitOpts struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	NodeAddress common.Address
}

//go:generate mockgen -destination=kaiax/headergov/mocks/headergov_mock.go github.com/kaiachain/kaia/kaiax/headergov HeaderGovModule
type headerGovModule struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	NodeAddress common.Address
	myVotes     []VoteData // queue

	epoch uint64
	cache HeaderCache
}

func NewHeaderGovModule() *headerGovModule {
	return &headerGovModule{}
}

func (h *headerGovModule) Init(opts *InitOpts) error {
	h.ChainKv = opts.ChainKv
	h.ChainConfig = opts.ChainConfig
	h.Chain = opts.Chain
	h.NodeAddress = opts.NodeAddress
	h.myVotes = make([]VoteData, 0)
	if h.ChainConfig == nil || h.ChainConfig.Istanbul == nil {
		return errNoChainConfig
	}

	h.epoch = h.ChainConfig.Istanbul.Epoch
	if h.epoch == 0 {
		return errZeroEpoch
	}

	votes := readVoteDataFromDB(h.Chain, h.ChainKv)
	govs := readGovDataFromDB(h.Chain, h.ChainKv)

	h.cache = *headergov_types.NewHeaderGovCache()
	for blockNum, vote := range votes {
		h.cache.AddVote(calcEpochIdx(blockNum, h.epoch), blockNum, vote)
	}
	for blockNum, gov := range govs {
		h.cache.AddGov(blockNum, gov)
	}

	return nil
}

func (s *headerGovModule) Start() error {
	logger.Info("HeaderGovModule started")
	return nil
}

func (s *headerGovModule) Stop() {
	logger.Info("HeaderGovModule stopped")
}

func (s *headerGovModule) isKoreHF(num uint64) bool {
	return s.ChainConfig.IsKoreForkEnabled(new(big.Int).SetUint64(num))
}

func (s *headerGovModule) PushMyVotes(vote VoteData) {
	s.myVotes = append(s.myVotes, vote)
}

func (s *headerGovModule) PopMyVotes(idx int) {
	s.myVotes = append(s.myVotes[:idx], s.myVotes[idx+1:]...)
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

			votes[blockNum] = parsedVote
		}
	}

	return votes
}

func readGovDataFromDB(chain chain, db database.Database) map[uint64]GovData {
	govBlocks := ReadGovDataBlockNums(db)
	govs := make(map[uint64]GovData)

	if govBlocks == nil {
		panic("govBlocks does not exist")
	}

	for _, blockNum := range *govBlocks {
		header := chain.GetHeaderByNumber(blockNum)
		parsedGov, err := headergov_types.DeserializeHeaderGov(header.Governance, blockNum)
		if err != nil {
			logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			panic(err)
		}

		govs[blockNum] = parsedGov
	}

	return govs
}
