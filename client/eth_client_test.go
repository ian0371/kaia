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

package client

import (
	"context"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestEthClient_MockServer(t *testing.T) {
	quitChan := make(chan struct{})
	defer close(quitChan)

	serverURL := MockHttpServer(t, quitChan)

	client, err := DialContext(context.Background(), serverURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Client connected to mock server")
	defer client.Close()

	kaiaHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	ethclient, err := DialContextEth(context.Background(), serverURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Eth client connected to mock server")
	ethHeader, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0x3b624db9bc6547b908e2e78460d2849047b6d28c0c078f09d6a0472ab0e57d0c", kaiaHeader.Hash().Hex())
	assert.Equal(t, "0x1c6ef781e4f30626053500c374498f78e3138128603e6f9c92bff0292613c5bb", ethHeader.Hash().Hex())
}

func TestEthClient_AnvilServer(t *testing.T) {
	serverURL := "http://127.0.0.1:8545"
	ethclient, err := DialContextEth(context.Background(), serverURL)
	if err != nil {
		t.Fatal(err)
	}

	// Test if server actually responds with a simple call
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if chainId, err := ethclient.ChainID(ctx); err != nil {
		t.Log("anvil server is not responding:", err)
		t.Skip("skip this test")
		return
	} else if chainId.Cmp(big.NewInt(1337)) != 0 {
		t.Fatal("the server must have chain id 1337, but got", chainId)
		return
	}

	t.Log("Eth client connected to anvil server")

	_, err = ethclient.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	header, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	cnt, err := ethclient.TransactionCount(context.Background(), header.Hash())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint(0), cnt)

	testAddrInitialBalance := big.NewInt(1e18)
	tx := func() *types.Transaction {
		richKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
		richAddr := crypto.PubkeyToAddress(richKey.PublicKey)
		if err != nil {
			t.Fatal(err)
		}
		nonce, err := ethclient.NonceAt(context.Background(), richAddr, nil)
		if err != nil {
			t.Fatal(err)
		}
		tx := types.NewTransaction(nonce, testAddr, testAddrInitialBalance, params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
		signer := types.LatestSignerForChainID(genesis.Config.ChainID)
		signedTx, _ := types.SignTx(tx, signer, richKey)
		return signedTx
	}()
	assert.Equal(t, header.ParentHash, common.Hash{})
	assert.NotEqual(t, header.Hash(), common.Hash{})

	hash, err := ethclient.SendRawTransaction(context.Background(), tx)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	receipt, err := ethclient.TransactionReceiptRpcOutput(context.Background(), hash)
	if err != nil {
		t.Fatal(err)
	}
	status, err := strconv.ParseUint(receipt["status"].(string), 0, 64)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, types.ReceiptStatusSuccessful, uint(status), "tx %s failed", hash.Hex())

	balance, err := ethclient.BalanceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, balance.Cmp(testAddrInitialBalance) >= 0)

	nonce, err := ethclient.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	dynamicTx := func() *types.Transaction {
		tx := types.NewTx(&types.TxInternalDataEthereumDynamicFee{
			ChainID:      genesis.Config.ChainID,
			AccountNonce: nonce,
			Recipient:    &testAddr,
			Amount:       big.NewInt(10),
			GasLimit:     25000,
			GasFeeCap:    big.NewInt(50e9),
			GasTipCap:    big.NewInt(25e9),
		})
		signer := types.LatestSignerForChainID(genesis.Config.ChainID)
		signedTx, _ := types.SignTx(tx, signer, testKey)
		return signedTx
	}()
	assert.Equal(t, header.ParentHash, common.Hash{})
	assert.NotEqual(t, header.Hash(), common.Hash{})
	_, err = ethclient.SendRawTransaction(context.Background(), dynamicTx)
	if err != nil {
		t.Fatal(err)
	}

	nonce++
	deployTx := func() *types.Transaction {
		// contract Storage { uint256 number = 1337; * @dev Return value @return value of 'number' */ function retrieve() public view returns (uint256){ return number; } }
		bytecode := common.Hex2Bytes("60806040526105395f553480156013575f5ffd5b5060af80601f5f395ff3fe6080604052348015600e575f5ffd5b50600436106026575f3560e01c80632e64cec114602a575b5f5ffd5b60306044565b604051603b91906062565b60405180910390f35b5f5f54905090565b5f819050919050565b605c81604c565b82525050565b5f60208201905060735f8301846055565b9291505056fea2646970667358221220bbed5c2a1719068dca0cf4da53d280029c147463a0b8f8319bc3494906ad27a964736f6c634300081e0033")
		tx := types.NewContractCreation(nonce, big.NewInt(0), 1e6, big.NewInt(25e9), bytecode)
		signer := types.LatestSignerForChainID(genesis.Config.ChainID)
		signedTx, _ := types.SignTx(tx, signer, testKey)
		return signedTx
	}()
	hash, err = ethclient.SendRawTransaction(context.Background(), deployTx)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	receipt, err = ethclient.TransactionReceiptRpcOutput(context.Background(), hash)
	if err != nil {
		t.Fatal(err)
	}
	status, err = strconv.ParseUint(receipt["status"].(string), 0, 64)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, types.ReceiptStatusSuccessful, uint(status), "tx %s failed", hash.Hex())

	contractAddr := crypto.CreateAddress(testAddr, nonce)
	calldata := common.Hex2Bytes("2e64cec1") // retrieve()(uint256)
	ret, err := ethclient.CallContract(context.Background(), kaia.CallMsg{To: &contractAddr, Data: calldata}, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000539", common.Bytes2Hex(ret))

	accesslist, _, _, err := ethclient.CreateAccessList(context.Background(), kaia.CallMsg{To: &contractAddr, Data: calldata})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, accesslist.StorageKeys())

	storage, err := ethclient.StorageAt(context.Background(), contractAddr, (*accesslist)[0].StorageKeys[0], nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000539", common.Bytes2Hex(storage))
}
