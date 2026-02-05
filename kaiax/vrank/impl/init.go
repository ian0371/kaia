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
//
// KIP-227: VRank module implementation (P2P + collection).

package impl

import (
	"crypto/ecdsa"
	"slices"
	"sync"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

const (
	candidatePrepareDeadlineMs = 200

	VRankPreprepareMsg = 0x17
	VRankCandidateMsg  = 0x18
)

var (
	_ vrank.VRankModule = &VRankModule{}

	logger = log.NewModuleLogger(log.KaiaxVrank)
)

type ProtocolManager interface {
	FindCNPeers(map[common.Address]bool) map[common.Address]consensus.Peer
}

type InitOpts struct {
	Valset      valset.ValsetModule
	PrivateKey  *ecdsa.PrivateKey
	ChainConfig *params.ChainConfig
	EventMux    *event.TypeMux
	Pm          ProtocolManager
}

type VRankModule struct {
	InitOpts

	nodeAddress   common.Address
	preprepareSub *event.TypeMuxSubscription

	candidateMu               sync.Mutex
	candidateCollection       map[uint64]map[common.Address]candidateEntry
	candidatePrepareStartTime map[uint64]time.Time
	stopCh                    chan struct{}
}

type candidateEntry struct {
	Sig     []byte
	Elapsed time.Duration
}

// NewVRankModule creates a new VRank module.
func NewVRankModule() *VRankModule {
	return &VRankModule{
		candidateCollection:       make(map[uint64]map[common.Address]candidateEntry),
		candidatePrepareStartTime: make(map[uint64]time.Time),
		stopCh:                    make(chan struct{}),
	}
}

// Init initializes the module with the given options.
func (v *VRankModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Valset == nil || opts.PrivateKey == nil || opts.ChainConfig == nil || opts.EventMux == nil || opts.Pm == nil {
		return errInitNil
	}
	v.InitOpts = *opts
	v.nodeAddress = crypto.PubkeyToAddress(opts.PrivateKey.PublicKey)
	v.preprepareSub = opts.EventMux.Subscribe(istanbul.PrepreparedEvent{})
	return nil
}

func (v *VRankModule) amIProposer(blockNum, round uint64) bool {
	proposer, err := v.Valset.GetProposer(blockNum, round)
	if err != nil {
		logger.Error("GetProposer failed", "blockNum", blockNum, "round", round)
		return false
	}

	return proposer == v.nodeAddress
}

func (v *VRankModule) amIValidator(blockNum uint64) bool {
	validators, err := v.Valset.GetCouncil(blockNum)
	if err != nil || validators == nil {
		logger.Error("GetCandidates failed", "blockNum", blockNum)
		return false
	}

	return slices.Contains(validators, v.nodeAddress)
}

func (v *VRankModule) amICandidate(blockNum uint64) bool {
	candidates, err := v.Valset.GetCandidates(blockNum)
	if err != nil || candidates == nil {
		logger.Error("GetCandidates failed", "blockNum", blockNum)
		return false
	}

	return candidates.Contains(v.nodeAddress)
}

func (v *VRankModule) BuildCfReportForBlock(blockNum uint64) []common.Address {
	v.candidateMu.Lock()
	defer v.candidateMu.Unlock()
	collected := v.candidateCollection[blockNum]
	var expectedCandidates []common.Address
	if candidates, err := v.Valset.GetCandidates(blockNum); err == nil && candidates != nil {
		expectedCandidates = candidates.List()
	}
	var cfReport []common.Address
	for _, addr := range expectedCandidates {
		entry, ok := collected[addr]
		if !ok || entry.Elapsed > candidatePrepareDeadlineMs {
			cfReport = append(cfReport, addr)
		}
	}
	delete(v.candidateCollection, blockNum)
	delete(v.candidatePrepareStartTime, blockNum)
	return cfReport
}

func (v *VRankModule) Start() error {
	logger.Info("VRankModule started")
	go v.handlePreprepareEvent()
	return nil
}

func (v *VRankModule) Stop() {
	logger.Info("VRankModule stopped")
	v.stopCh <- struct{}{}
}
