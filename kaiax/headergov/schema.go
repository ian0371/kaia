package headergov

import (
	"encoding/json"

	"github.com/kaiachain/kaia/storage/database"
)

var (
	voteDataBlockNumsKey = []byte("voteDataBlockNums")
	govDataBlockNumsKey  = []byte("governanceDataBlockNums")
)

type StoredVoteBlockNums []uint64
type StoredGovBlockNums []uint64

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
