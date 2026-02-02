// Copyright 2024 The Kaia Authors
// KIP-227 §4: header.VRank encode/decode tests.

package types

import (
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeVrankPayload_Empty(t *testing.T) {
	// Empty payload: RLP([])
	p := &VRankPayload{}
	enc, err := EncodeVRankPayload(p)
	assert.NoError(t, err)
	assert.NotEmpty(t, enc)

	dec, err := DecodeVRankPayload(enc)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Empty(t, dec.CfReport)
}

func TestEncodeDecodeVrankPayload_Nil(t *testing.T) {
	enc, err := EncodeVRankPayload(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, enc)

	dec, err := DecodeVRankPayload(enc)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Empty(t, dec.CfReport)
}

func TestEncodeDecodeVrankPayload_CfReportOnly(t *testing.T) {
	// cfReport = list of candidate addresses who failed (no signatures)
	addrs := []common.Address{
		common.HexToAddress("0x15d34AAf54267DB7D7cC839724318F2730aC377B"),
		common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819D0A4DC"),
	}
	p := &VRankPayload{CfReport: addrs}
	enc, err := EncodeVRankPayload(p)
	assert.NoError(t, err)
	assert.NotEmpty(t, enc)

	dec, err := DecodeVRankPayload(enc)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Len(t, dec.CfReport, 2)
	assert.Equal(t, addrs[0], dec.CfReport[0])
	assert.Equal(t, addrs[1], dec.CfReport[1])
}

func TestDecodeVrankPayload_EmptyBytes(t *testing.T) {
	dec, err := DecodeVRankPayload(nil)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Empty(t, dec.CfReport)

	dec, err = DecodeVRankPayload([]byte{})
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Empty(t, dec.CfReport)
}

func TestDecodeVrankPayload_InvalidRLP(t *testing.T) {
	_, err := DecodeVRankPayload([]byte{0xff, 0xff})
	assert.Error(t, err)
}
