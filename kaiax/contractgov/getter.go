package contractgov

import (
	"math/big"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/common"
	govcontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/gov"
)

func (c *contractGovModule) EffectiveParamSet(blockNum uint64) (ParamSet, error) {
	pset, err := c.contractGetAllParamsAt(blockNum)
	if err != nil {
		return ParamSet{}, err
	}
	return pset, nil
}

func (c *contractGovModule) contractGetAllParamsAt(blockNum uint64) (ParamSet, error) {
	chain := c.Chain
	if chain == nil {
		return ParamSet{}, errContractEngineNotReady
	}

	config := c.ChainConfig
	if !config.IsKoreForkEnabled(new(big.Int).SetUint64(blockNum)) {
		return ParamSet{}, errContractEngineNotReady
	}

	addr, err := c.contractAddrAt(blockNum)
	if err != nil {
		return ParamSet{}, err
	}
	if common.EmptyAddress(addr) {
		logger.Trace("ContractEngine disabled: GovParamContract address not set")
		return ParamSet{}, nil
	}

	caller := backends.NewBlockchainContractBackend(chain, nil, nil)
	contract, _ := govcontract.NewGovParamCaller(addr, caller)

	names, values, err := contract.GetAllParamsAt(nil, new(big.Int).SetUint64(blockNum))
	if err != nil {
		logger.Warn("ContractEngine disabled: getAllParams call failed", "err", err)
		return ParamSet{}, nil
	}

	if len(names) != len(values) {
		logger.Warn("ContractEngine disabled: getAllParams result invalid", "len(names)", len(names), "len(values)", len(values))
		return ParamSet{}, nil
	}

	ret := ParamSet{}
	for i := 0; i < len(names); i++ {
		param := Params[names[i]]
		cv, err := param.Canonicalizer(values[i])
		if err != nil {
			return ParamSet{}, err
		}
		err = ret.Set(names[i], cv)
		if err != nil {
			return ParamSet{}, err
		}
	}

	return ret, nil
}

func (c *contractGovModule) contractAddrAt(blockNum uint64) (common.Address, error) {
	headerParams, err := c.hgm.EffectiveParamSet(blockNum)
	if err != nil {
		return common.Address{}, errParamsAtFail
	}

	return headerParams.GovParamContract, nil
}
