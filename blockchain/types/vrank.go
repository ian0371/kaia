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

// VRankPayload is the decoded content of header.VRank (KIP-227 §4).
// header.VRank contains only cfReport: RLP(cfReport) = RLP list of failing candidate addresses.
// pfReport is not stored in the header; it is kept in consensus state only.
type VRankPayload struct {
	// CfReport is the list of candidate addresses who failed (e.g. did not send valid VRankCandidate within timeout). Nil when there are no failures.
	CfReport []common.Address
}

// EncodeVRankPayload RLP-encodes the payload as a single list of addresses (cfReport).
// When CfReport is nil or empty, callers may set header.VRank = nil instead of encoding.
func EncodeVRankPayload(p *VRankPayload) ([]byte, error) {
	if p == nil {
		p = &VRankPayload{}
	}
	return rlp.EncodeToBytes(p.CfReport)
}

// DecodeVRankPayload decodes header.VRank bytes into CfReport.
// data must be a valid RLP list of addresses.
func DecodeVRankPayload(data []byte) (*VRankPayload, error) {
	if len(data) == 0 {
		return &VRankPayload{}, nil
	}
	var cfReport []common.Address
	if err := rlp.DecodeBytes(data, &cfReport); err != nil {
		return nil, err
	}
	return &VRankPayload{CfReport: cfReport}, nil
}
