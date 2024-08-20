package types

import (
	"encoding/json"
	"errors"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

type param struct {
	name  string
	value interface{}
}

type voteData struct {
	blockNum uint64
	voter    common.Address
	param    param
}

type governanceData struct {
	blockNum uint64
	params   []param
}

type allParamsHistory map[string]*PartitionList[param]

type HeaderGovernanceReader struct {
	votes  []voteData
	govs   []governanceData
	epoch  uint64
	koreHf uint64
}

type ChainReader interface {
	GetHeaderByNumber(number uint64) *types.Header
}

func NewHeaderGovernanceReader(chain ChainReader, db database.Database, koreHf uint64) *HeaderGovernanceReader {
	voteBlocks := gov.ReadVoteDataBlocks(db)
	votes := make([]voteData, 0)
	govs := make([]governanceData, 0)
	if voteBlocks != nil {
		for _, blockNum := range *voteBlocks {
			header := chain.GetHeaderByNumber(blockNum)
			parsedVote, err := parseHeaderVote(header.Vote, blockNum)
			if err != nil {
				gov.Logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			}

			votes = append(votes, *parsedVote)
		}
	}

	govBlocks := gov.ReadGovDataBlocks(db)
	if govBlocks != nil {
		for _, blockNum := range *govBlocks {
			header := chain.GetHeaderByNumber(blockNum)
			parsedGov, err := parseHeaderGov(header.Governance, blockNum)
			if err != nil {
				gov.Logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			}

			govs = append(govs, *parsedGov)
		}
	}

	return &HeaderGovernanceReader{
		votes:  votes,
		govs:   govs,
		epoch:  params.DefaultEpoch,
		koreHf: koreHf,
	}
}

func (h *HeaderGovernanceReader) AddVote(num uint64, vote voteData) {
	h.votes = append(h.votes, vote)
}

func (h *HeaderGovernanceReader) AddGov(num uint64, g governanceData) error {
	if num%h.epoch != 0 {
		return errors.New("gov block num must be epoch aligned")
	}
	h.govs = append(h.govs, g)
	return nil
}

func (h *HeaderGovernanceReader) EffectiveParams(num uint64) map[string]param {
	allParamsHistory := getAllParamsHistory(h.govs, h.epoch, h.koreHf)
	ret := make(map[string]param)
	for k, v := range allParamsHistory {
		ret[k] = v.GetItem(uint(num))
	}
	return ret
}

func getAllParamsHistory(govs []governanceData, epoch uint64, kaiaHf uint64) allParamsHistory {
	ret := make(allParamsHistory)

	for _, g := range govs {
		for _, p := range g.params {
			ia, ok := ret[p.name]
			if !ok {
				ret[p.name] = &PartitionList[param]{}
				ia = ret[p.name]
			}

			activation := uint64(0)
			if g.blockNum < kaiaHf {
				activation = headerActivationBlockPreKore(g.blockNum, epoch)
			} else {
				activation = headerActivationBlockPostKore(g.blockNum, epoch)
			}
			ia.AddRecord(uint(activation), p)
		}
	}

	return ret
}

func headerActivationBlockPreKore(num, epoch uint64) uint64 {
	if num == 0 {
		return 0
	}
	return num + epoch + 1
}

func headerActivationBlockPostKore(num, epoch uint64) uint64 {
	if num == 0 {
		return 0
	}
	return num + epoch
}

func parseHeaderVote(b []byte, blockNum uint64) (*voteData, error) {
	var v struct {
		Validator common.Address
		Key       string
		Value     interface{}
	}

	err := rlp.DecodeBytes(b, &v)
	if err != nil {
		return nil, err
	}

	return &voteData{
		blockNum: blockNum,
		voter:    v.Validator,
		param:    param{name: v.Key, value: v.Value},
	}, nil
}

func parseHeaderGov(b []byte, blockNum uint64) (*governanceData, error) {
	rlpDecoded := []byte("")
	err := rlp.DecodeBytes(b, &rlpDecoded)
	if err != nil {
		return nil, err
	}

	j := make(map[string]interface{})
	json.Unmarshal(rlpDecoded, &j)

	params := make([]param, 0, len(j))
	for k, v := range j {
		params = append(params, param{name: k, value: v})
	}

	return &governanceData{
		blockNum: blockNum,
		params:   params,
	}, nil
}
