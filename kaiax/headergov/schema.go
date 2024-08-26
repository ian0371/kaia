package headergov

import (
	"encoding/json"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	voteDataBlockNumsKey = []byte("voteDataBlockNums")
	govDataBlockNumsKey  = []byte("governanceDataBlockNums")
	govParamKey          = []byte("governanceParam-")
)

type StoredVoteBlockNums []uint64
type StoredGovBlockNums []uint64

func makeKey(prefix []byte, num uint64) []byte {
	byteKey := common.Int64ToByteLittleEndian(num)
	return append(prefix, byteKey...)
}

func ReadVoteDataBlockNums(db database.Database) *StoredVoteBlockNums {
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

func WriteVoteDataBlockNums(db database.Database, voteDataBlockNums *StoredVoteBlockNums) {
	b, err := json.Marshal(voteDataBlockNums)
	if err != nil {
		logger.Error("Failed to marshal voteDataBlocks", "err", err)
		return
	}

	if err := db.Put(voteDataBlockNumsKey, b); err != nil {
		logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}

func ReadGovDataBlockNums(db database.Database) *StoredGovBlockNums {
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

func WriteGovDataBlockNums(db database.Database, govData *StoredGovBlockNums) {
	b, err := json.Marshal(govData)
	if err != nil {
		logger.Error("Failed to marshal govDataBlocks", "err", err)
		return
	}

	if err := db.Put(govDataBlockNumsKey, b); err != nil {
		logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}

func ReadGovernanceParam(db database.Database, num uint64) *GovernanceParam {
	b, err := db.Get(makeKey(govParamKey, num))
	if err != nil || len(b) == 0 {
		return nil
	}

	ps := new(GovernanceParam)
	if err := json.Unmarshal(b, ps); err != nil {
		logger.Error("Invalid GovParamSet JSON", "num", num, "err", err)
		return nil
	}
	return ps
}

func WriteGovernanceParam(db database.Database, num uint64, ps *GovernanceParam) {
	b, err := json.Marshal(ps)
	if err != nil {
		logger.Error("Failed to marshal govParams", "err", err)
		return
	}

	if err := db.Put(makeKey(govParamKey, num), b); err != nil {
		logger.Crit("Failed to write GovParamSet", "num", num, "err", err)
	}
}
