package headergov

import (
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
			Namespace: "kaia",
			Version:   "1.0",
			Service:   newHeaderGovAPI(s),
			Public:    true,
		},
	}
}

type headerGovAPI struct {
	s *HeaderGovModule
}

func newHeaderGovAPI(s *HeaderGovModule) *headerGovAPI {
	return &headerGovAPI{s}
}

func (api *headerGovAPI) Vote(key string, val interface{}) (string, error) {
	blockNumber := api.s.Chain.CurrentBlock().NumberU64()
	gp, err := api.s.EffectiveParams(blockNumber + 1)
	if err != nil {
		return "", err
	}

	gMode := gp.GovernanceMode
	if gMode == params.GovernanceMode_Single && api.s.NodeAddress != gp.GoverningNode {
		return "", errPermissionDenied
	}

	err = api.s.VerifyVote(&VoteData{
		BlockNum: blockNumber,
		Voter:    api.s.NodeAddress,
		Name:     key,
		Value:    val,
	}, gp)
	if err != nil {
		return "", err
	}

	// TODO: check if val is in the validator set for addval, removeval
	if key == "governance.removevalidator" {
		if val.(common.Address) == api.s.NodeAddress {
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

	api.s.MyVotes = append(api.s.MyVotes, VoteData{Name: key, Value: val})
	return "Your vote is prepared. It will be put into the block header or applied when your node generates a block as a proposer. Note that your vote may be duplicate.", nil
}
