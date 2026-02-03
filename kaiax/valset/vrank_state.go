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

package valset

import (
	"crypto/sha512"
	"strconv"
	"strings"
	"sync"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/tyler-smith/go-bip32"
	"golang.org/x/crypto/pbkdf2"
)

// VrankStatus is the KIP-227 validator/candidate status.
type VrankStatus int

const (
	ValActive   VrankStatus = iota // active validator; participates in consensus
	CandTesting                    // candidate under evaluation; sends CandidatePrepare
)

// Mock constants per spec §3: N addresses from mnemonic, first numVals are ValActive, rest CandTesting.
const (
	mockNumCNs   = 6
	mockNumVals  = 4
	mockNumCands = 2
)

func init() {
	if mockNumCNs != mockNumVals+mockNumCands {
		panic("valset vrank mock: numCNs must equal numVals + numCands")
	}
}

// ReadVrankStates returns the VRank status per address for the given block (KIP-227 §3).
// MVP: mock implementation returning a fixed map derived from the well-known mnemonic
// "(test x11, junk)" and Ethereum default BIP-39 path m/44'/60'/0'/0/i.
// First mockNumVals addresses are ValActive, next mockNumCands are CandTesting.
func ReadVrankStates(blockNum uint64) map[common.Address]VrankStatus {
	return readVrankStatesMock()
}

var (
	mockVrankOnce sync.Once
	mockVrankMap  map[common.Address]VrankStatus
)

// Mnemonic and path per spec §3: "We use the well-known mnemonic (test x11, junk) addresses for mock.
// Given N, generate N HD private keys based on the mnemonic and the Ethereum default BIP-39 path."
const (
	mockMnemonic = "test test test test test test test test test test test junk"
	mockPath     = "m/44'/60'/0'/0"
)

// generateAddressesFromMnemonic derives num addresses from mnemonic and path using the same
// logic as cmd/homi/common.GenerateKeysFromMnemonic (BIP-32/BIP-44). Path format: e.g. "m/44'/60'/0'/0".
func generateAddressesFromMnemonic(num int, mnemonic, path string) []common.Address {
	var key *bip32.Key
	for _, level := range strings.Split(path, "/") {
		if level == "" {
			continue
		}
		if level == "m" {
			seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"), 2048, 64, sha512.New)
			var err error
			key, err = bip32.NewMasterKey(seed)
			if err != nil {
				panic("valset vrank mock: NewMasterKey: " + err.Error())
			}
			continue
		}
		var n uint64
		var err error
		if strings.HasSuffix(level, "'") {
			n, err = strconv.ParseUint(level[:len(level)-1], 10, 32)
			if err != nil {
				panic("valset vrank mock: parse path " + level + ": " + err.Error())
			}
			n += 0x80000000
		} else {
			n, err = strconv.ParseUint(level, 10, 32)
			if err != nil {
				panic("valset vrank mock: parse path " + level + ": " + err.Error())
			}
		}
		key, err = key.NewChildKey(uint32(n))
		if err != nil {
			panic("valset vrank mock: NewChildKey " + level + ": " + err.Error())
		}
	}

	addrs := make([]common.Address, 0, num)
	for i := 0; i < num; i++ {
		derived, err := key.NewChildKey(uint32(i))
		if err != nil {
			panic("valset vrank mock: child " + strconv.Itoa(i) + ": " + err.Error())
		}
		// Same as cmd/homi/common: hex key (no "0x") for HexToECDSA
		nodekey := hexutil.Encode(derived.Key)[2:]
		priv, err := crypto.HexToECDSA(nodekey)
		if err != nil {
			// derived.Key may be 33 bytes (leading 0x00); ToECDSA expects 32
			keyBytes := derived.Key
			if len(keyBytes) == 33 && keyBytes[0] == 0 {
				keyBytes = keyBytes[1:]
			}
			priv, err = crypto.ToECDSA(keyBytes)
			if err != nil {
				panic("valset vrank mock: ToECDSA: " + err.Error())
			}
		}
		addrs = append(addrs, crypto.PubkeyToAddress(priv.PublicKey))
	}
	return addrs
}

func readVrankStatesMock() map[common.Address]VrankStatus {
	mockVrankOnce.Do(func() {
		addrs := generateAddressesFromMnemonic(mockNumCNs, mockMnemonic, mockPath)
		mockVrankMap = make(map[common.Address]VrankStatus, len(addrs))
		for i, addr := range addrs {
			if i < mockNumVals {
				mockVrankMap[addr] = ValActive
			} else {
				mockVrankMap[addr] = CandTesting
			}
		}
	})
	return mockVrankMap
}
