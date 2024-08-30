package headergov

import (
	"encoding/json"
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
type VotesInEpoch = headergov_types.VotesInEpoch
type GovData = headergov_types.GovData
type GovHistory = headergov_types.GovHistory
type GovParam = headergov_types.GovParamSet
type GovHeaderCache = headergov_types.GovHeaderCache

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
	myVotes     []VoteData // queue

	epoch uint64
	cache GovHeaderCache
}

func NewHeaderGovModule() *HeaderGovModule {
	return &HeaderGovModule{}
}

func (h *HeaderGovModule) Init(opts *InitOpts) error {
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
		h.cache.AddGovVote(calcEpochIdx(blockNum, h.epoch), blockNum, vote)
	}
	for blockNum, gov := range govs {
		h.cache.AddGov(blockNum, gov)
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

func (s *HeaderGovModule) PushMyVotes(vote VoteData) {
	s.myVotes = append(s.myVotes, vote)
}

func (s *HeaderGovModule) PopMyVotes(idx int) {
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

			votes[blockNum] = *parsedVote
		}
	}

	return votes
}

func groupVotesByEpoch(votes map[uint64]VoteData, epoch uint64) map[uint64]VotesInEpoch {
	groupedVotes := make(map[uint64]VotesInEpoch)
	for blockNum, vote := range votes {
		epochIdx := calcEpochIdx(blockNum, epoch)
		if _, ok := groupedVotes[epochIdx]; !ok {
			groupedVotes[epochIdx] = make(VotesInEpoch)
		}
		groupedVotes[epochIdx][blockNum] = vote
	}
	return groupedVotes
}

func readGovDataFromDB(chain chain, db database.Database) map[uint64]GovData {
	govBlocks := ReadGovDataBlockNums(db)
	govs := make(map[uint64]GovData)

	// TODO: remove this.
	if govBlocks == nil {
		governanceHistoryKey := []byte("governanceIdxHistory")
		history, err := db.Get(governanceHistoryKey)
		if err != nil {
			logger.Error("Failed to read governance history", "err", err)
			return govs
		}
		idxHistory := make([]uint64, 0)
		if err := json.Unmarshal(history, &idxHistory); err != nil {
			logger.Error("Failed to unmarshal governance history", "err", err)
			return govs
		}
		govBlocks = (*StoredGovBlockNums)(&idxHistory)
	}

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
