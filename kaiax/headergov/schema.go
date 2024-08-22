package headergov

import (
	"encoding/json"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	voteDataBlockNumsKey = []byte("voteDataBlockNums")
	govDataBlockNumsKey  = []byte("governanceDataBlockNums")
	govParamSetKey       = []byte("governanceParamSet-")
)

type StoredVoteBlockNums []uint64
type StoredGovBlockNums []uint64

func makeKey(prefix []byte, num uint64) []byte {
	byteKey := common.Int64ToByteLittleEndian(num)
	return append(prefix, byteKey...)
}

func ReadVoteDataBlocks(db database.Database) *StoredVoteBlockNums {
	b, err := db.Get(voteDataBlockNumsKey)
	if err != nil || len(b) == 0 {
		return nil
	}

	ret := new(StoredVoteBlockNums)
	if err := json.Unmarshal(b, ret); err != nil {
		logger.Error("Invalid voteDataBlocks JSON", "err", err)
		return nil
	}
	return ret
}

func WriteVoteDataBlocks(db database.Database, voteDataBlockNums *StoredVoteBlockNums) {
	b, err := json.Marshal(voteDataBlockNums)
	if err != nil {
		logger.Error("Failed to marshal voteDataBlocks", "err", err)
		return
	}

	if err := db.Put(voteDataBlockNumsKey, b); err != nil {
		logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}

func ReadGovDataBlocks(db database.Database) *StoredGovBlockNums {
	b, err := db.Get(govDataBlockNumsKey)
	if err != nil || len(b) == 0 {
		return nil
	}

	ret := new(StoredGovBlockNums)
	if err := json.Unmarshal(b, ret); err != nil {
		logger.Error("Invalid govDataBlocks JSON", "err", err)
		return nil
	}
	return ret
}

func WriteGovDataBlocks(db database.Database, govData *StoredGovBlockNums) {
	b, err := json.Marshal(govData)
	if err != nil {
		logger.Error("Failed to marshal govDataBlocks", "err", err)
		return
	}

	if err := db.Put(govDataBlockNumsKey, b); err != nil {
		logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}

func ReadGovParams(db database.Database, num uint64) *params.GovParamSet {
	b, err := db.Get(makeKey(govParamSetKey, num))
	if err != nil || len(b) == 0 {
		return nil
	}

	ps := new(params.GovParamSet)
	if err := json.Unmarshal(b, ps); err != nil {
		logger.Error("Invalid govParams JSON", "err", err)
		return nil
	}
	return ps
}

func WriteGovParams(db database.Database, num uint64, ps *params.GovParamSet) {
	b, err := json.Marshal(ps)
	if err != nil {
		logger.Error("Failed to marshal govParams", "err", err)
		return
	}

	if err := db.Put(makeKey(govParamSetKey, num), b); err != nil {
		logger.Crit("Failed to write govParams", "err", err)
	}
}
