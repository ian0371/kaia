package headergov

import (
	"encoding/json"
	"errors"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
)

var (
	errUnknownBlock           = errors.New("Unknown block")
	errNotAvailableInThisMode = errors.New("In current governance mode, voting power is not available")
	errSetDefaultFailure      = errors.New("Failed to set a default value")
	errPermissionDenied       = errors.New("You don't have the right to vote")
	errRemoveSelf             = errors.New("You can't vote on removing yourself")
	errInvalidKeyValue        = errors.New("Your vote couldn't be placed. Please check your vote's key and value")
	errInvalidLowerBound      = errors.New("lowerboundbasefee cannot be set exceeding upperboundbasefee")
	errInvalidUpperBound      = errors.New("upperboundbasefee cannot be set lower than lowerboundbasefee")
)

func (s *HeaderGovModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "governance",
			Version:   "1.0",
			Service:   newHeaderGovAPI(s),
			Public:    true,
		},
	}
}

type headerGovAPI struct {
	h *HeaderGovModule
}

func newHeaderGovAPI(s *HeaderGovModule) *headerGovAPI {
	return &headerGovAPI{s}
}

func (api *headerGovAPI) Vote(key string, val interface{}) (string, error) {
	blockNumber := api.h.Chain.CurrentBlock().NumberU64()
	gp, err := api.h.EffectiveParams(blockNumber + 1)
	if err != nil {
		return "", err
	}

	gMode := gp.GovernanceMode
	if gMode == params.GovernanceMode_Single && api.h.NodeAddress != gp.GoverningNode {
		return "", errPermissionDenied
	}

	// TODO: add string handler
	if _, ok := val.(float64); ok {
		val = uint64(val.(float64))
	}

	err = api.h.VerifyVote(&VoteData{
		Voter: api.h.NodeAddress,
		Name:  key,
		Value: val,
	})
	if err != nil {
		return "", err
	}

	// TODO: check if val is in the validator set for addval, removeval
	if key == "governance.removevalidator" {
		if val.(common.Address) == api.h.NodeAddress {
			return "", errRemoveSelf
		}
	}

	if key == "kip71.lowerboundbasefee" {
		if val.(uint64) > gp.UpperBoundBaseFee {
			return "", errInvalidLowerBound
		}
	}
	if key == "kip71.upperboundbasefee" {
		if val.(uint64) < gp.LowerBoundBaseFee {
			return "", errInvalidUpperBound
		}
	}

	api.h.PushMyVotes(VoteData{Voter: api.h.NodeAddress, Name: key, Value: val})
	return "(kaiax) Your vote is prepared. It will be put into the block header or applied when your node generates a block as a proposer. Note that your vote may be duplicate.", nil
}

func (api *headerGovAPI) IdxCache() []uint64 {
	return api.h.cache.GovBlockNums()
}

type MyVotesAPI struct {
	BlockNum uint64
	Casted   bool
	Key      string
	Value    interface{}
}

func (api *headerGovAPI) MyVotes() []MyVotesAPI {
	epochIdx := calcEpochIdx(api.h.Chain.CurrentBlock().NumberU64(), api.h.epoch)
	votesInEpoch := api.h.getVotesInEpoch(epochIdx)

	ret := make([]MyVotesAPI, 0)
	for blockNum, vote := range votesInEpoch {
		if vote.Voter == api.h.NodeAddress {
			ret = append(ret, MyVotesAPI{
				BlockNum: blockNum,
				Casted:   true,
				Key:      vote.Name,
				Value:    vote.Value,
			})
		}
	}

	for _, vote := range api.h.myVotes {
		ret = append(ret, MyVotesAPI{
			BlockNum: 0,
			Casted:   false,
			Key:      vote.Name,
			Value:    vote.Value,
		})
	}

	return ret
}

func (api *headerGovAPI) PendingVotes() []VoteData {
	epochIdx := calcEpochIdx(api.h.Chain.CurrentBlock().NumberU64(), api.h.epoch)
	votesInEpoch := api.h.getVotesInEpoch(epochIdx)

	ret := make([]VoteData, 0)
	for _, vote := range votesInEpoch {
		if vote.Voter == api.h.NodeAddress {
			ret = append(ret, vote)
		}
	}

	return ret
}

func (api *headerGovAPI) NodeAddress() common.Address {
	return api.h.NodeAddress
}

func (api *headerGovAPI) GetParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	return api.getParams(num)
}

func (api *headerGovAPI) getParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = api.h.Chain.CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	gp, err := api.h.EffectiveParams(blockNumber)
	if err != nil {
		return nil, err
	}
	return gp.ToStrMap()
}

func (api *headerGovAPI) Status() (string, error) {
	type PublicCache struct {
		GroupedVotes map[uint64]VotesInEpoch   `json:"groupedVotes"`
		Governances  map[uint64]GovernanceData `json:"governances"`
		GovHistory   GovernanceHistory         `json:"govHistory"`
	}
	publicCache := PublicCache{
		GroupedVotes: api.h.cache.GroupedVotes(),
		Governances:  api.h.cache.Govs(),
		GovHistory:   api.h.cache.GetGovernanceHistory(),
	}

	cacheJson, err := json.Marshal(publicCache)
	if err != nil {
		logger.Error("kaiax: Failed to marshal cache", "err", err)
		return "", err
	}

	return string(cacheJson), nil
}
