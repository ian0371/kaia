package headergov

import (
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
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

type Param = headergov_types.Param
type VoteData = headergov_types.VoteData
type GovernanceData = headergov_types.GovernanceData
type GovBlockNumToGovParamSetMap = headergov_types.GovBlockNumToGovParamSetMap
type GovernanceCache = headergov_types.GovernanceCache

type chain interface {
	GetHeaderByNumber(number uint64) *types.Header
}

type InitOpts struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
}
type HeaderGovModule struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain

	epoch uint64
	cache GovernanceCache
}

func (h *HeaderGovModule) Init(opts *InitOpts) error {
	h.ChainKv = opts.ChainKv
	h.ChainConfig = opts.ChainConfig
	h.Chain = opts.Chain
	if h.ChainConfig == nil || h.ChainConfig.Istanbul == nil {
		return errNoChainConfig
	}

	h.epoch = h.ChainConfig.Istanbul.Epoch
	if h.epoch == 0 {
		return errZeroEpoch
	}

	h.cache = GovernanceCache{
		Votes:  readVoteBlockNumsFromDB(h.Chain, h.ChainKv),
		Govs:   readGovBlockNumsFromDB(h.Chain, h.ChainKv),
		GovMap: readGovMapFromDB(h.Chain, h.ChainKv),
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

func readVoteBlockNumsFromDB(chain chain, db database.Database) []VoteData {
	voteBlocks := ReadVoteDataBlockNums(db)
	votes := make([]VoteData, 0)
	if voteBlocks != nil {
		for _, blockNum := range *voteBlocks {
			header := chain.GetHeaderByNumber(blockNum)
			parsedVote, err := headergov_types.DeserializeHeaderVote(header.Vote, blockNum)
			if err != nil {
				logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			}

			votes = append(votes, *parsedVote)
		}
	}

	return votes
}

func readGovBlockNumsFromDB(chain chain, db database.Database) []GovernanceData {
	govBlocks := ReadGovDataBlockNums(db)
	govs := make([]GovernanceData, 0)
	if govBlocks != nil {
		for _, blockNum := range *govBlocks {
			header := chain.GetHeaderByNumber(blockNum)
			parsedGov, err := headergov_types.DeserializeHeaderGov(header.Governance, blockNum)
			if err != nil {
				logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			}

			govs = append(govs, *parsedGov)
		}
	}
	return govs
}

func readGovMapFromDB(chain chain, db database.Database) GovBlockNumToGovParamSetMap {
	govBlockNums := ReadGovDataBlockNums(db)
	if govBlockNums == nil || len(*govBlockNums) == 0 {
		return GovBlockNumToGovParamSetMap{}
	}

	govParams := GovBlockNumToGovParamSetMap{}
	for _, govBlockNum := range *govBlockNums {
		govParams.AddRecord(uint(govBlockNum), ReadGovParamSet(db, govBlockNum))
	}

	return govParams
}
