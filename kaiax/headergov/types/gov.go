package types

import (
	"encoding/json"
	"errors"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

type Param struct {
	Name  string
	Value interface{} // canonical value
}

type VoteData struct {
	BlockNum uint64
	Voter    common.Address
	Param    Param
}

type GovernanceData struct {
	BlockNum uint64
	Params   map[string]Param
}

type AllParamsHistory map[string]*PartitionList[Param] // p1 -> [(activation1, value1), ...], p2 -> ...

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

func (h *GovernanceCache) AddVote(num uint64, vote VoteData) {
	h.Votes = append(h.Votes, vote)
}

func (h *GovernanceCache) AddGov(num uint64, g GovernanceData) {
	h.Govs = append(h.Govs, g)
}

func GetAllParamsHistory(govs []GovernanceData, epoch uint64, koreHf uint64) AllParamsHistory {
	ret := make(AllParamsHistory)

	for _, g := range govs {
		for paramName, p := range g.Params {
			activation := uint64(0)
			if g.BlockNum < koreHf {
				activation = headerGovActivationBlockPreKore(g.BlockNum, epoch)
			} else {
				activation = headerGovActivationBlockPostKore(g.BlockNum, epoch)
			}

			if _, exists := ret[paramName]; !exists {
				ret[paramName] = &PartitionList[Param]{}
			}
			ret[paramName].AddRecord(uint(activation), p)
		}
	}

	return ret
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

func (vote *VoteData) Serialize() ([]byte, error) {
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

func DeserializeHeaderVote(b []byte, blockNum uint64) (*VoteData, error) {
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

func (g *GovernanceData) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for _, v := range g.Params {
		tmp[v.Name] = v.Value
	}

	return json.Marshal(tmp)
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
