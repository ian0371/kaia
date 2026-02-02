// Copyright 2025 The Kaia Authors
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

package types

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
)

// VrankPayload is the decoded content of header.Vrank (KIP-227 §4).
// header.Vrank = nil when both pfReport and cfReport are nil; else RLP([pfReport, cfReport]).
type VrankPayload struct {
	// PfReport is the list of proposer addresses that caused round change in the previous block (round 0..N-1 before commit). Nil when round = 0.
	PfReport []common.Address
	// CfReport is the list of candidate addresses who failed (e.g. did not send valid VRankCandidate within timeout). Nil when there are no failures.
	CfReport []common.Address
}

// EncodeVrankPayload RLP-encodes the payload as [pfReport, cfReport] (KIP-227 §4).
// When both are nil, callers should set header.Vrank = nil instead of encoding.
// Empty slices encode as RLP empty lists; result is valid for header.Vrank at/after fork.
func EncodeVrankPayload(p *VrankPayload) ([]byte, error) {
	if p == nil {
		p = &VrankPayload{}
	}
	return rlp.EncodeToBytes(p)
}

// DecodeVrankPayload decodes header.Vrank bytes into PfReport and CfReport.
// Returns error if data is not a valid two-element RLP list (address list, address list).
func DecodeVrankPayload(data []byte) (*VrankPayload, error) {
	if len(data) == 0 {
		return &VrankPayload{}, nil
	}
	var p VrankPayload
	if err := rlp.DecodeBytes(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
