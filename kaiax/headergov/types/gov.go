package types

import (
	"encoding/json"
	"errors"

	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

type Param struct {
	Name  string
	Value interface{} // canonical value
}

type GovernanceData struct {
	BlockNum uint64
	Params   map[string]Param
}

type GovernanceCache struct {
	Votes []VoteData
	Govs  []GovernanceData
}

func (h *GovernanceCache) VoteBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for _, vote := range h.Votes {
		blockNums = append(blockNums, vote.BlockNum)
	}
	return blockNums
}

func (h *GovernanceCache) GovBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for _, vote := range h.Govs {
		blockNums = append(blockNums, vote.BlockNum)
	}
	return blockNums
}

// TODO: sort
func (h *GovernanceCache) AddVote(num uint64, vote VoteData) {
	h.Votes = append(h.Votes, vote)
}

// TODO: sort
func (h *GovernanceCache) AddGov(num uint64, g GovernanceData) {
	h.Govs = append(h.Govs, g)
}

func GetAllParamsHistory(govs []GovernanceData) (PartitionList[*params.GovParamSet], error) {
	ret := PartitionList[*params.GovParamSet]{}

	effectiveParams := &params.GovParamSet{}
	for _, g := range govs {
		ps, err := g.ToParamSet()
		if err != nil {
			return ret, err
		}
		effectiveParams = params.NewGovParamSetMerged(effectiveParams, ps)
		copied := *effectiveParams
		ret.AddRecord(uint(g.BlockNum), &copied)
	}

	return ret, nil
}

func headerGovActivationBlockPreKore(govDataBlockNum, epoch uint64) uint64 {
	if govDataBlockNum == 0 {
		return 0
	}
	return govDataBlockNum + epoch + 1
}

func headerGovActivationBlockPostKore(govDataBlockNum, epoch uint64) uint64 {
	if govDataBlockNum == 0 {
		return 0
	}
	return govDataBlockNum + epoch
}

func (g *GovernanceData) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for _, v := range g.Params {
		tmp[v.Name] = v.Value
	}

	return json.Marshal(tmp)
}

func (g *GovernanceData) ToParamSet() (*params.GovParamSet, error) {
	tmp := make(map[string]interface{})
	for _, v := range g.Params {
		tmp[v.Name] = v.Value
	}
	return params.NewGovParamSetStrMap(tmp)
}

func (g *GovernanceData) Serialize() ([]byte, error) {
	j, err := g.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(j)
}

func DeserializeHeaderGov(b []byte, blockNum uint64) (*GovernanceData, error) {
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
