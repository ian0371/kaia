package headergov

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/governance"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	logger = log.NewModuleLogger(log.KaiaXGov) // TODO: rename as "logger" (small case)

	errZeroEpoch     = errors.New("epoch cannot be zero")
	errNoChainConfig = errors.New("ChainConfig or Istanbul is not set")
)

type Param = headergov_types.Param
type VoteData = headergov_types.VoteData
type GovernanceData = headergov_types.GovernanceData
type AllParamsHistory = headergov_types.AllParamsHistory
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
		Votes: readVoteBlockNumsFromDB(h.Chain, h.ChainKv),
		Govs:  readGovBlockNumsFromDB(h.Chain, h.ChainKv),
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
	voteBlocks := ReadVoteDataBlocks(db)
	votes := make([]VoteData, 0)
	if voteBlocks != nil {
		for _, blockNum := range *voteBlocks {
			header := chain.GetHeaderByNumber(blockNum)
			parsedVote, err := deserializeHeaderVote(header.Vote, blockNum)
			if err != nil {
				logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			}

			votes = append(votes, *parsedVote)
		}
	}

	return votes
}

func readGovBlockNumsFromDB(chain chain, db database.Database) []GovernanceData {
	govBlocks := ReadGovDataBlocks(db)
	govs := make([]GovernanceData, 0)
	if govBlocks != nil {
		for _, blockNum := range *govBlocks {
			header := chain.GetHeaderByNumber(blockNum)
			parsedGov, err := deserializeHeaderGov(header.Governance, blockNum)
			if err != nil {
				logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			}

			govs = append(govs, *parsedGov)
		}
	}
	return govs
}

func deserializeHeaderVote(b []byte, blockNum uint64) (*VoteData, error) {
	var v struct {
		Validator common.Address
		Key       string
		Value     []byte
	}

	err := rlp.DecodeBytes(b, &v)
	if err != nil {
		return nil, err
	}

	// canonicalize. e.g., [0x1, 0xc9, 0xc3, 0x80] -> 0x1c9c380
	ps, err := params.NewGovParamSetBytesMap(map[string][]byte{
		v.Key: v.Value,
	})
	if err != nil {
		return nil, err
	}

	value, ok := ps.Get(governance.GovernanceKeyMap[v.Key])
	if !ok {
		return nil, errors.New("key not found")
	}

	return &VoteData{
		BlockNum: blockNum,
		Voter:    v.Validator,
		Param:    Param{Name: v.Key, Value: value},
	}, nil
}

func serializeVoteData(vote *VoteData) ([]byte, error) {
	v := &struct {
		Validator common.Address
		Key       string
		Value     interface{}
	}{
		Validator: vote.Voter,
		Key:       vote.Param.Name,
		Value:     vote.Param.Value,
	}

	return rlp.EncodeToBytes(v)
}

func deserializeHeaderGov(b []byte, blockNum uint64) (*GovernanceData, error) {
	rlpDecoded := []byte("")
	err := rlp.DecodeBytes(b, &rlpDecoded)
	if err != nil {
		return nil, err
	}

	j := make(map[string]interface{})
	json.Unmarshal(rlpDecoded, &j)

	ps, err := params.NewGovParamSetStrMap(j)
	if err != nil {
		return nil, err
	}

	params := make(map[string]Param, len(j))
	for k := range j {
		value, ok := ps.Get(governance.GovernanceKeyMap[k])
		if !ok {
			return nil, errors.New("key not found")
		}
		params[k] = Param{Name: k, Value: value}
	}

	return &GovernanceData{
		BlockNum: blockNum,
		Params:   params,
	}, nil
}
