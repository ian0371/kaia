package gov

import (
	"encoding/json"

	"github.com/kaiachain/kaia/storage/database"
)

var (
	voteDataBlocksKey = []byte("voteDataBlocks")
	govDataBlocksKey  = []byte("governanceDataBlocks")
)

type storedVoteData []uint64
type storedGovData []uint64

func ReadVoteDataBlocks(db database.Database) *storedVoteData {
	b, err := db.Get(voteDataBlocksKey)
	if err != nil || len(b) == 0 {
		return nil
	}

	ret := new(storedVoteData)
	if err := json.Unmarshal(b, ret); err != nil {
		Logger.Error("Invalid voteDataBlocks JSON", "err", err)
		return nil
	}
	return ret
}

func WriteVoteDataBlocks(db database.Database, voteData *storedVoteData) {
	b, err := json.Marshal(voteData)
	if err != nil {
		Logger.Error("Failed to marshal voteDataBlocks", "err", err)
		return
	}

	if err := db.Put(voteDataBlocksKey, b); err != nil {
		Logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}

func ReadGovDataBlocks(db database.Database) *storedGovData {
	b, err := db.Get(govDataBlocksKey)
	if err != nil || len(b) == 0 {
		return nil
	}

	ret := new(storedGovData)
	if err := json.Unmarshal(b, ret); err != nil {
		Logger.Error("Invalid govDataBlocks JSON", "err", err)
		return nil
	}
	return ret
}

func WriteGovDataBlocks(db database.Database, govData *storedGovData) {
	b, err := json.Marshal(govData)
	if err != nil {
		Logger.Error("Failed to marshal govDataBlocks", "err", err)
		return
	}

	if err := db.Put(govDataBlocksKey, b); err != nil {
		Logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}
