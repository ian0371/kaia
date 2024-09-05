package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/rlp"
)

type GovData struct {
	Params map[string]interface{} // canonicalized value
}

func (g *GovData) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for name, value := range g.Params {
		if bigInt, ok := value.(*big.Int); ok {
			tmp[name] = bigInt.String()
		} else {
			tmp[name] = value
		}
	}

	return json.Marshal(tmp)
}

func (g *GovData) Serialize() ([]byte, error) {
	j, err := g.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(j)
}

func DeserializeHeaderGov(b []byte, blockNum uint64) (*GovData, error) {
	rlpDecoded := []byte("")
	err := rlp.DecodeBytes(b, &rlpDecoded)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal(rlpDecoded, &ret)
	if err != nil {
		return nil, err
	}

	for name, value := range ret {
		param, ok := Params[name]
		if !ok {
			return nil, errors.New("invalid param")
		}

		cv, err := param.Canonicalizer(value)
		if err != nil {
			return nil, err
		}
		ret[name] = cv
	}

	return &GovData{
		Params: ret,
	}, nil
}
