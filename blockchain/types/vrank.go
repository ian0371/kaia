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
// header.Vrank contains only cfReport: RLP(cfReport) = RLP list of failing candidate addresses.
// pfReport is not stored in the header; it is kept in consensus state only.
type VrankPayload struct {
	// CfReport is the list of candidate addresses who failed (e.g. did not send valid VRankCandidate within timeout). Nil when there are no failures.
	CfReport []common.Address
}

// EncodeVrankPayload RLP-encodes the payload as a single list of addresses (cfReport).
// When CfReport is nil or empty, callers may set header.Vrank = nil instead of encoding.
func EncodeVrankPayload(p *VrankPayload) ([]byte, error) {
	if p == nil {
		p = &VrankPayload{}
	}
	return rlp.EncodeToBytes(p.CfReport)
}

// DecodeVrankPayload decodes header.Vrank bytes into CfReport.
// data must be a valid RLP list of addresses.
func DecodeVrankPayload(data []byte) (*VrankPayload, error) {
	if len(data) == 0 {
		return &VrankPayload{}, nil
	}
	var cfReport []common.Address
	if err := rlp.DecodeBytes(data, &cfReport); err != nil {
		return nil, err
	}
	return &VrankPayload{CfReport: cfReport}, nil
}
