package impl

import (
	"math/big"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// handlePreprepare is executed by validators and the proposer
func (v *VRankModule) handlePreprepare(block *types.Block, view *istanbul.View) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return
	}

	blockNum := block.NumberU64()
	logger.Info("HandlePreprepare",
		"blockNum", blockNum, "view", view,
		"AmIValidator", v.amIValidator(blockNum),
		"AmIProposer", v.amIProposer(blockNum, view.Round.Uint64()),
	)
	if v.amIValidator(blockNum) {
		v.candidatePrepareStartTime[blockNum] = time.Now()
	}
	if v.amIProposer(blockNum, view.Round.Uint64()) {
		vrankPreprepare := &vrank.VRankPreprepare{Block: block}
		v.BroadcastVRankPreprepare(vrankPreprepare)
	}
}

func (v *VRankModule) handlePreprepareEvent() {
	for {
		select {
		case <-v.stopCh:
			logger.Info("VRankModule stopped")
			return
		case ev, ok := <-v.preprepareSub.Chan():
			if !ok || ev.Data == nil {
				logger.Warn("Drop an empty message from subscription channel")
				return
			}

			data, ok := ev.Data.(istanbul.PrepreparedEvent)
			if !ok || data.Block == nil || data.View == nil {
				logger.Warn("PrepreparedEvent message arrived. but type is not istanbul.PrepreparedEvent or some data is nil")
				return
			}
			v.handlePreprepare(data.Block, data.View)
		}
	}
}

// HandleVRankPreprepare is executed by candidates
func (v *VRankModule) HandleVRankPreprepare(preprepare *vrank.VRankPreprepare) {
	if preprepare == nil || preprepare.Block == nil {
		logger.Error("Unexpected nil")
		return
	}
	block := preprepare.Block
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return
	}
	if v.amICandidate(block.NumberU64()) {
		logger.Info("HandleVRankPreprepare",
			"blockNum", block.NumberU64(),
		)
		blockHash := block.Hash()
		sig, err := crypto.Sign(crypto.Keccak256(blockHash.Bytes()), v.PrivateKey)
		if err != nil {
			logger.Error("Sign failed", "blockNum", block.NumberU64())
			return
		}
		msg := &vrank.VRankCandidate{
			BlockNumber: block.NumberU64(),
			Round:       0,
			BlockHash:   blockHash,
			Sig:         sig,
		}
		v.BroadcastVRankCandidate(msg)
	}
}

// HandleVRankCandidate is executed by validators
func (v *VRankModule) HandleVRankCandidate(msg *vrank.VRankCandidate) {
	if msg == nil {
		logger.Error("Unexpected nil")
		return
	}
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(msg.BlockNumber)) {
		return
	}

	elapsed := time.Since(v.candidatePrepareStartTime[msg.BlockNumber])
	logger.Info("HandleVRankCandidate",
		"blockNum", msg.BlockNumber,
		"elapsed", elapsed,
		"AmIValidator", v.amIValidator(msg.BlockNumber),
	)
	if v.amIValidator(msg.BlockNumber) {
		cand, err := istanbul.GetSignatureAddress(crypto.Keccak256(msg.BlockHash.Bytes()), msg.Sig)
		if err != nil {
			logger.Error("GetSignatureAddress failed", "blockNum", msg.BlockNumber, "blockHash", msg.BlockHash, "sig", msg.Sig)
			return
		}
		v.addCandidateResponse(cand, msg, elapsed)
	}
}

func (v *VRankModule) addCandidateResponse(cand common.Address, msg *vrank.VRankCandidate, elapsed time.Duration) {
	logger.Info("addCandidateResponse", "blockNum", msg.BlockNumber, "cand", cand, "elapsed", elapsed)
	v.candidateMu.Lock()
	defer v.candidateMu.Unlock()
	if _, ok := v.candidatePrepareStartTime[msg.BlockNumber]; !ok {
		return
	}
	if v.candidateCollection[msg.BlockNumber] == nil {
		return
	}
	if _, exists := v.candidateCollection[msg.BlockNumber][cand]; exists {
		return
	}
	v.candidateCollection[msg.BlockNumber][cand] = candidateEntry{Sig: msg.Sig, Elapsed: elapsed}
}

func (v *VRankModule) BroadcastVRankPreprepare(vrankPreprepare *vrank.VRankPreprepare) {
	block := vrankPreprepare.Block
	candidates, err := v.Valset.GetCandidates(block.NumberU64())
	if err != nil || candidates == nil {
		logger.Error("GetCandidates failed", "blockNum", block.NumberU64())
		return
	}
	v.broadcast(candidates.List(), vrankPreprepare)
}

func (v *VRankModule) BroadcastVRankCandidate(vrankCandidate *vrank.VRankCandidate) {
	validators, err := v.Valset.GetCouncil(vrankCandidate.BlockNumber)
	if err != nil || validators == nil {
		logger.Error("GetCouncil failed", "blockNum", vrankCandidate.BlockNumber)
		return
	}

	v.broadcast(validators, vrankCandidate)
}

func (v *VRankModule) broadcast(targets []common.Address, msg any) {
	peerTargets := make(map[common.Address]bool)
	for _, target := range targets {
		peerTargets[target] = true
	}

	peers := v.Pm.FindCNPeers(peerTargets)
	for addr, p := range peers {
		switch msg.(type) {
		case *vrank.VRankPreprepare:
			msg := msg.(*vrank.VRankPreprepare)
			logger.Info("broadcast VRankPreprepareMsg", "blockNum", msg.Block.NumberU64(), "to", addr)
			go p.Send(VRankPreprepareMsg, msg)
		case *vrank.VRankCandidate:
			msg := msg.(*vrank.VRankCandidate)
			logger.Info("broadcast VRankCandidateMsg", "blockNum", msg.BlockNumber, "to", addr)
			go p.Send(VRankCandidateMsg, msg)
		}
	}
}
