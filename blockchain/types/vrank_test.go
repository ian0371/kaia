// Copyright 2024 The Kaia Authors
// KIP-227 §4: header.Vrank encode/decode tests.

package types

import (
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeVrankPayload_Empty(t *testing.T) {
	// Empty payload: RLP([[], []])
	p := &VrankPayload{}
	enc, err := EncodeVrankPayload(p)
	assert.NoError(t, err)
	assert.NotEmpty(t, enc)

	dec, err := DecodeVrankPayload(enc)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Empty(t, dec.PfReport)
	assert.Empty(t, dec.CfReport)
}

func TestEncodeDecodeVrankPayload_Nil(t *testing.T) {
	enc, err := EncodeVrankPayload(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, enc)

	dec, err := DecodeVrankPayload(enc)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Empty(t, dec.PfReport)
	assert.Empty(t, dec.CfReport)
}

func TestEncodeDecodeVrankPayload_PfReportOnly(t *testing.T) {
	addrs := []common.Address{
		common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
	}
	p := &VrankPayload{PfReport: addrs}
	enc, err := EncodeVrankPayload(p)
	assert.NoError(t, err)
	assert.NotEmpty(t, enc)

	dec, err := DecodeVrankPayload(enc)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Len(t, dec.PfReport, 2)
	assert.Equal(t, addrs[0], dec.PfReport[0])
	assert.Equal(t, addrs[1], dec.PfReport[1])
	assert.Empty(t, dec.CfReport)
}

func TestEncodeDecodeVrankPayload_CfReportOnly(t *testing.T) {
	// cfReport = list of candidate addresses who failed (no signatures)
	addrs := []common.Address{
		common.HexToAddress("0x15d34AAf54267DB7D7cC839724318F2730aC377B"),
		common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819D0A4DC"),
	}
	p := &VrankPayload{CfReport: addrs}
	enc, err := EncodeVrankPayload(p)
	assert.NoError(t, err)
	assert.NotEmpty(t, enc)

	dec, err := DecodeVrankPayload(enc)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Empty(t, dec.PfReport)
	assert.Len(t, dec.CfReport, 2)
	assert.Equal(t, addrs[0], dec.CfReport[0])
	assert.Equal(t, addrs[1], dec.CfReport[1])
}

func TestEncodeDecodeVrankPayload_Full(t *testing.T) {
	pfAddrs := []common.Address{common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")}
	cfAddrs := []common.Address{common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819D0A4DC")}
	p := &VrankPayload{PfReport: pfAddrs, CfReport: cfAddrs}
	enc, err := EncodeVrankPayload(p)
	assert.NoError(t, err)
	assert.NotEmpty(t, enc)

	dec, err := DecodeVrankPayload(enc)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Len(t, dec.PfReport, 1)
	assert.Equal(t, pfAddrs[0], dec.PfReport[0])
	assert.Len(t, dec.CfReport, 1)
	assert.Equal(t, cfAddrs[0], dec.CfReport[0])
}

func TestDecodeVrankPayload_EmptyBytes(t *testing.T) {
	dec, err := DecodeVrankPayload(nil)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Empty(t, dec.PfReport)
	assert.Empty(t, dec.CfReport)

	dec, err = DecodeVrankPayload([]byte{})
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Empty(t, dec.PfReport)
	assert.Empty(t, dec.CfReport)
}

func TestDecodeVrankPayload_InvalidRLP(t *testing.T) {
	_, err := DecodeVrankPayload([]byte{0xff, 0xff})
	assert.Error(t, err)
}
