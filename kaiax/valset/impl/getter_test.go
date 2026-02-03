// Copyright 2024 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package impl

import (
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

// TestGetCandidates checks that GetCandidates returns exactly the CandTesting
// addresses from the readVRankStates mock (KIP-227 §4).
// Expected addresses: indices 4 and 5 from mnemonic "test test test ... junk", path m/44'/60'/0'/0/i.
func TestGetCandidates(t *testing.T) {
	v := &ValsetModule{}

	expectedCandTesting := []common.Address{
		common.HexToAddress("0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65"),
		common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc"),
	}

	for _, blockNum := range []uint64{0, 1, 100} {
		candidates, err := v.GetCandidates(blockNum)
		assert.NoError(t, err)
		assert.NotNil(t, candidates)
		assert.Equal(t, 2, candidates.Len(), "blockNum=%d", blockNum)
		list := candidates.List()
		assert.ElementsMatch(t, expectedCandTesting, list, "blockNum=%d", blockNum)
	}
}
