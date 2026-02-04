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
	"math/big"
	"sync"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

const (
	candidatePrepareDeadlineMs = 200
)

var logger = log.NewModuleLogger(log.KaiaxVrank)

var _ vrank.VRankModule = &VRankModule{}

type ProtocolManager interface {
	FindPeers(map[common.Address]bool) map[common.Address]consensus.Peer
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

	nodeAddress  common.Address
	subscription *event.TypeMuxSubscription

	candidateMu               sync.Mutex
	candidateCollection       map[uint64]map[common.Address]candidateEntry
	candidatePrepareStartTime map[uint64]time.Time
	stopCh                    chan struct{}
}

type candidateEntry struct {
	Sig         []byte
	ArrivalTime time.Time
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
	v.subscription = opts.EventMux.Subscribe(istanbul.PrepreparedEvent{})
	return nil
}

func (v *VRankModule) AmICandidate(blockNum uint64) bool {
	candidates, err := v.Valset.GetCandidates(blockNum)
	if err != nil || candidates == nil || !candidates.Contains(v.nodeAddress) {
		logger.Error("GetCandidates failed", "blockNum", blockNum)
		return false
	}

	return candidates.Contains(v.nodeAddress)
}

func (v *VRankModule) SendVRankPreprepare(block *types.Block) {
}

func (v *VRankModule) SendVRankCandidate(block *types.Block) {
}

func (v *VRankModule) HandlePreprepare(block *types.Block, view *istanbul.View) {
	logger.Warn("Preprepare message arrived",
		"blockNumber", block.Number(),
		"blockHash", block.Hash(),
		"round", view.Round.Uint64(),
	)
}

func (v *VRankModule) HandleVRankPreprepare(addr common.Address, payload []byte) {
	block := new(types.Block)
	if err := rlp.DecodeBytes(payload, block); err != nil {
		return
	}
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return
	}
	blockNum := block.NumberU64()
	if !v.AmICandidate(blockNum) { // only candidate handles VRankPreprepare
		return
	}

	/*
		blockHash := block.Hash()
		sig, err := crypto.Sign(blockHash.Bytes(), v.PrivateKey)
		if err != nil {
			return
		}
		vcp := &istanbul.VRankCandidatePayload{
			BlockNumber: blockNum,
			Round:       0,
			BlockHash:   blockHash,
			Sig:         sig,
		}
		vPayload, err := istanbul.EncodeVRankCandidatePayload(vcp)
		if err != nil {
			return
		}
		v.candidateBroadcastFeed.Send(vPayload)
	*/
}

func (v *VRankModule) HandleVRankCandidate(addr common.Address, payload []byte) {
	p, err := istanbul.DecodeVRankCandidatePayload(payload)
	if err != nil {
		return
	}
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(p.BlockNumber)) {
		return
	}
	if !v.AmICandidate(p.BlockNumber) {
		return
	}

	signer, err := istanbul.GetSignatureAddress(p.BlockHash.Bytes(), p.Sig)
	if err != nil {
		return
	}
	v.addCandidateResponse(signer, p.BlockNumber, uint64(p.Round), p.BlockHash, p.Sig)
}

func (v *VRankModule) addCandidateResponse(signer common.Address, blockNum, round uint64, blockHash common.Hash, sig []byte) {
	nextProposer, err := v.Valset.GetProposer(blockNum+1, 0)
	if err != nil || nextProposer != v.nodeAddress {
		return
	}
	v.candidateMu.Lock()
	defer v.candidateMu.Unlock()
	if _, ok := v.candidatePrepareStartTime[blockNum]; !ok {
		return
	}
	if v.candidateCollection[blockNum] == nil {
		return
	}
	if _, exists := v.candidateCollection[blockNum][signer]; exists {
		return
	}
	v.candidateCollection[blockNum][signer] = candidateEntry{Sig: sig, ArrivalTime: time.Now()}
}

func (v *VRankModule) BuildCfReportForBlock(blockNum uint64) []common.Address {
	v.candidateMu.Lock()
	defer v.candidateMu.Unlock()
	start, ok := v.candidatePrepareStartTime[blockNum]
	if !ok {
		return nil
	}
	deadline := start.Add(candidatePrepareDeadlineMs * time.Millisecond)
	collected := v.candidateCollection[blockNum]
	var expectedCandidates []common.Address
	if candidates, err := v.Valset.GetCandidates(blockNum); err == nil && candidates != nil {
		expectedCandidates = candidates.List()
	}
	var cfReport []common.Address
	for _, addr := range expectedCandidates {
		entry, ok := collected[addr]
		if !ok || entry.ArrivalTime.After(deadline) {
			cfReport = append(cfReport, addr)
		}
	}
	delete(v.candidateCollection, blockNum)
	delete(v.candidatePrepareStartTime, blockNum)
	return cfReport
}

func (v *VRankModule) handleMsg() {
	for {
		select {
		case <-v.stopCh:
			logger.Warn("VRankModule stopped")
			return
		case ev, ok := <-v.subscription.Chan():
			if !ok || ev.Data == nil {
				logger.Warn("Drop an empty message from subscription channel")
				return
			}

			data, ok := ev.Data.(istanbul.PrepreparedEvent)
			if !ok || data.Block == nil || data.View == nil {
				logger.Warn("PrepreparedEvent message arrived. but type is not istanbul.PrepreparedEvent or some data is nil")
				return
			}
			logger.Warn("PrepreparedEvent message arrived. calling HandlePreprepare")
			v.HandlePreprepare(data.Block, data.View)
		}
	}
}

func (v *VRankModule) Start() error {
	logger.Info("VRankModule started")
	go v.handleMsg()
	return nil
}

func (v *VRankModule) Stop() {
	logger.Info("VRankModule stopped")
	v.stopCh <- struct{}{}
}
