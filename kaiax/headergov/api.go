package headergov

import (
	"errors"

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
	pset, err := api.s.EffectiveParams(blockNumber + 1)
	if err != nil {
		return "", err
	}

	// TODO: check node address
	gMode := pset.GovernanceModeInt()
	if gMode == params.GovernanceMode_Single || true {
		return "", errPermissionDenied
	}
	err = api.s.VerifyVote(key, val)
	if err != nil {
		return "", err
	}

	return "", errInvalidKeyValue
}
