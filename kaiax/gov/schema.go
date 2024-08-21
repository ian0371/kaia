package gov

import (
	"encoding/json"

	"github.com/kaiachain/kaia/storage/database"
)

var (
	voteDataBlockNumsKey = []byte("voteDataBlocks")
	govDataBlockNumsKey  = []byte("governanceDataBlocks")
)

type StoredVoteData []uint64
type StoredGovData []uint64

func ReadVoteDataBlocks(db database.Database) *StoredVoteData {
	b, err := db.Get(voteDataBlockNumsKey)
	if err != nil || len(b) == 0 {
		return nil
	}

	ret := new(StoredVoteData)
	if err := json.Unmarshal(b, ret); err != nil {
		Logger.Error("Invalid voteDataBlocks JSON", "err", err)
		return nil
	}
	return ret
}

func WriteVoteDataBlocks(db database.Database, voteDataBlockNums *StoredVoteData) {
	b, err := json.Marshal(voteDataBlockNums)
	if err != nil {
		Logger.Error("Failed to marshal voteDataBlocks", "err", err)
		return
	}

	if err := db.Put(voteDataBlockNumsKey, b); err != nil {
		Logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}

func ReadGovDataBlocks(db database.Database) *StoredGovData {
	b, err := db.Get(govDataBlockNumsKey)
	if err != nil || len(b) == 0 {
		return nil
	}

	ret := new(StoredGovData)
	if err := json.Unmarshal(b, ret); err != nil {
		Logger.Error("Invalid govDataBlocks JSON", "err", err)
		return nil
	}
	return ret
}

func WriteGovDataBlocks(db database.Database, govData *StoredGovData) {
	b, err := json.Marshal(govData)
	if err != nil {
		Logger.Error("Failed to marshal govDataBlocks", "err", err)
		return
	}

	if err := db.Put(govDataBlockNumsKey, b); err != nil {
		Logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}
