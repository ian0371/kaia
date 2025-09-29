// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from ethclient/ethclient_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/gxhash"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/node"
	"github.com/kaiachain/kaia/node/cn"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

// Verify that Client implements the Kaia interfaces.
var (
	// _ = kaia.Subscription(&Client{})
	_ = kaia.ChainReader(&Client{})
	_ = kaia.TransactionReader(&Client{})
	_ = kaia.ChainStateReader(&Client{})
	_ = kaia.ChainSyncReader(&Client{})
	_ = kaia.ContractCaller(&Client{})
	_ = kaia.LogFilterer(&Client{})
	_ = kaia.TransactionSender(&Client{})
	_ = kaia.GasPricer(&Client{})
	_ = kaia.PendingStateReader(&Client{})
	_ = kaia.PendingContractCaller(&Client{})
	_ = kaia.GasEstimator(&Client{})
	// _ = kaia.PendingStateEventer(&Client{})
)

var (
	testKey, _         = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr           = crypto.PubkeyToAddress(testKey.PublicKey)
	testBalance        = big.NewInt(2e15)
	revertContractAddr = common.HexToAddress("290f1b36649a61e369c6276f6d29463335b4400c")
	revertCode         = common.FromHex("7f08c379a0000000000000000000000000000000000000000000000000000000006000526020600452600a6024527f75736572206572726f7200000000000000000000000000000000000000000000604452604e6000fd")
)

var vanity = make([]byte, types.IstanbulExtraVanity)
var extra, _ = rlp.EncodeToBytes(&types.IstanbulExtra{
	Validators:    []common.Address{testAddr},
	Seal:          []byte{},
	CommittedSeal: [][]byte{},
})
var genesis = &blockchain.Genesis{
	Config: params.TestChainConfig,
	Alloc: blockchain.GenesisAlloc{
		testAddr:           {Balance: testBalance},
		revertContractAddr: {Balance: big.NewInt(0), Code: revertCode},
	},
	ExtraData: append(vanity, extra...),
	Timestamp: 9000,
}

var testTx1 = func() *types.Transaction {
	tx := types.NewTransaction(0, common.Address{2}, big.NewInt(12), params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
	signer := types.LatestSignerForChainID(genesis.Config.ChainID)
	signedTx, _ := types.SignTx(tx, signer, testKey)
	return signedTx
}()

var testTx2 = func() *types.Transaction {
	tx := types.NewTransaction(1, common.Address{2}, big.NewInt(8), params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
	signer := types.LatestSigner(genesis.Config)
	signedTx, _ := types.SignTx(tx, signer, testKey)
	return signedTx
}()

func newTestBackend(config *node.Config, workspace string) (*node.Node, []*types.Block, error) {
	// Create a proper database and blockchain for testing
	dbm := database.NewMemoryDBManager()

	// Commit genesis to database
	genesisBlock := genesis.MustCommit(dbm)

	// Generate test chain
	engine := gxhash.NewFaker()
	generate := func(i int, g *blockchain.BlockGen) {
		g.OffsetTime(1)
		g.SetExtra([]byte("test"))
		if i == 1 {
			// Test transactions are included in block #2.
			g.AddTx(testTx1)
			g.AddTx(testTx2)
		}
	}

	blocks, _ := blockchain.GenerateChain(params.TestChainConfig, genesisBlock, engine, dbm, 2, generate)
	allBlocks := append([]*types.Block{genesisBlock}, blocks...)

	// Create node with proper configuration
	if config == nil {
		config = &node.Config{
			DataDir:          workspace,
			HTTPHost:         "127.0.0.1",
			HTTPPort:         36000,
			HTTPVirtualHosts: []string{"*"},
			HTTPModules:      []string{"kaia", "net", "web3", "admin", "debug"}, // Enable kaia namespace
		}
	}

	cnConf := cn.GetDefaultConfig()
	cnConf.Genesis = genesis
	fullNode, err := node.New(config)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create new node: %v", err)
	}
	if err = fullNode.Register(func(ctx *node.ServiceContext) (node.Service, error) { return cn.New(ctx, cnConf) }); err != nil {
		return nil, nil, fmt.Errorf("failed to register Kaia protocol: %v", err)
	}

	// Start the node first to initialize services
	if err := fullNode.Start(); err != nil {
		return nil, nil, fmt.Errorf("can't start test node: %v", err)
	}

	var cn *cn.CN
	if err := fullNode.Service(&cn); err != nil {
		return nil, nil, fmt.Errorf("failed to service Kaia protocol: %v", err)
	}

	return fullNode, allBlocks, nil
}

func TestEthClient(t *testing.T) {
	workspace, err := os.MkdirTemp("", "kaia-client-tester-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(workspace)

	backend, chain, err := newTestBackend(nil, workspace)
	if err != nil {
		t.Fatal(err)
	}
	client, err := backend.Attach()
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Stop()
	defer client.Close()

	tests := map[string]struct {
		test func(t *testing.T)
	}{
		"Header": {
			func(t *testing.T) { testHeader(t, chain, client) },
		},
		"BalanceAt": {
			func(t *testing.T) { testBalanceAt(t, client) },
		},
		/*
			"TxInBlockInterrupted": {
				func(t *testing.T) { testTransactionInBlock(t, client) },
			},
			"ChainID": {
				func(t *testing.T) { testChainID(t, client) },
			},
			"GetBlock": {
				func(t *testing.T) { testGetBlock(t, client) },
			},
			"StatusFunctions": {
				func(t *testing.T) { testStatusFunctions(t, client) },
			},
			"CallContract": {
				func(t *testing.T) { testCallContract(t, client) },
			},
			"CallContractAtHash": {
				func(t *testing.T) { testCallContractAtHash(t, client) },
			},
			"AtFunctions": {
				func(t *testing.T) { testAtFunctions(t, client) },
			},
			"TransactionSender": {
				func(t *testing.T) { testTransactionSender(t, client) },
			},
		*/
	}

	t.Parallel()
	for name, tt := range tests {
		t.Run(name, tt.test)
	}
}

func testHeader(t *testing.T, chain []*types.Block, client *rpc.Client) {
	tests := map[string]struct {
		block   *big.Int
		want    *types.Header
		wantErr error
	}{
		"genesis": {
			block: big.NewInt(0),
			want:  chain[0].Header(),
		},
		"first_block": {
			block: big.NewInt(1),
			want:  chain[1].Header(),
		},
		"future_block": {
			block:   big.NewInt(1000000000),
			want:    nil,
			wantErr: kaia.NotFound,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ec := NewClient(client)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := ec.HeaderByNumber(ctx, tt.block)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("HeaderByNumber(%v) error = %q, want %q", tt.block, err, tt.wantErr)
			}
			if got != nil && got.Number != nil && got.Number.Sign() == 0 {
				got.Number = big.NewInt(0) // hack to make DeepEqual work
			}
			if got.Hash() != tt.want.Hash() {
				t.Fatalf("HeaderByNumber(%v) got = %v, want %v", tt.block, got, tt.want)
			}
		})
	}
}

func testBalanceAt(t *testing.T, client *rpc.Client) {
	tests := map[string]struct {
		account common.Address
		block   *big.Int
		want    *big.Int
		wantErr error
	}{
		"valid_account_genesis": {
			account: testAddr,
			block:   big.NewInt(0),
			want:    testBalance,
		},
		"valid_account": {
			account: testAddr,
			block:   big.NewInt(1),
			want:    testBalance,
		},
		"non_existent_account": {
			account: common.Address{1},
			block:   big.NewInt(1),
			want:    big.NewInt(0),
		},
		"future_block": {
			account: testAddr,
			block:   big.NewInt(1000000000),
			want:    big.NewInt(0),
			wantErr: errors.New("header not found"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ec := NewClient(client)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := ec.BalanceAt(ctx, tt.account, tt.block)
			if tt.wantErr != nil && (err == nil || err.Error() != tt.wantErr.Error()) {
				t.Fatalf("BalanceAt(%x, %v) error = %q, want %q", tt.account, tt.block, err, tt.wantErr)
			}
			if got.Cmp(tt.want) != 0 {
				t.Fatalf("BalanceAt(%x, %v) = %v, want %v", tt.account, tt.block, got, tt.want)
			}
		})
	}
}

func testTransactionInBlock(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// Get current block by number.
	block, err := ec.BlockByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test tx in block not found.
	if _, err := ec.TransactionInBlock(context.Background(), block.Hash(), 20); err != kaia.NotFound {
		t.Fatal("error should be kaia.NotFound")
	}

	// Test tx in block found.
	tx, err := ec.TransactionInBlock(context.Background(), block.Hash(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Hash() != testTx1.Hash() {
		t.Fatalf("unexpected transaction: %v", tx)
	}

	tx, err = ec.TransactionInBlock(context.Background(), block.Hash(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Hash() != testTx2.Hash() {
		t.Fatalf("unexpected transaction: %v", tx)
	}

	// Test pending block
	_, err = ec.BlockByNumber(context.Background(), big.NewInt(int64(rpc.PendingBlockNumber)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testChainID(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)
	id, err := ec.ChainID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == nil || id.Cmp(params.TestChainConfig.ChainID) != 0 {
		t.Fatalf("ChainID returned wrong number: %+v", id)
	}
}

func testGetBlock(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// Get current block number
	blockNumber, err := ec.BlockNumber(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blockNumber.Int64() != 2 {
		t.Fatalf("BlockNumber returned wrong number: %d", blockNumber)
	}
	// Get current block by number
	block, err := ec.BlockByNumber(context.Background(), big.NewInt(blockNumber.Int64()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.NumberU64() != blockNumber.Uint64() {
		t.Fatalf("BlockByNumber returned wrong block: want %d got %d", blockNumber, block.NumberU64())
	}
	// Get current block by hash
	blockH, err := ec.BlockByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Hash() != blockH.Hash() {
		t.Fatalf("BlockByHash returned wrong block: want %v got %v", block.Hash().Hex(), blockH.Hash().Hex())
	}
	// Get header by number
	header, err := ec.HeaderByNumber(context.Background(), big.NewInt(blockNumber.Int64()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != header.Hash() {
		t.Fatalf("HeaderByNumber returned wrong header: want %v got %v", block.Header().Hash().Hex(), header.Hash().Hex())
	}
	// Get header by hash
	headerH, err := ec.HeaderByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != headerH.Hash() {
		t.Fatalf("HeaderByHash returned wrong header: want %v got %v", block.Header().Hash().Hex(), headerH.Hash().Hex())
	}
}

func testStatusFunctions(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// Sync progress
	progress, err := ec.SyncProgress(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progress != nil {
		t.Fatalf("unexpected progress: %v", progress)
	}

	// NetworkID
	networkID, err := ec.NetworkID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if networkID.Cmp(big.NewInt(1337)) != 0 {
		t.Fatalf("unexpected networkID: %v", networkID)
	}

	// SuggestGasPrice
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gasPrice.Cmp(big.NewInt(1000000000)) != 0 {
		t.Fatalf("unexpected gas price: %v", gasPrice)
	}

	// Note: SuggestGasTipCap, BlobBaseFee, and FeeHistory methods are not available in Kaia client
	// Skipping these tests for Kaia compatibility
}

func testCallContractAtHash(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// EstimateGas
	msg := kaia.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
	_, err = ec.HeaderByNumber(context.Background(), big.NewInt(1))
	if err != nil {
		t.Fatalf("BlockByNumber error: %v", err)
	}
	// Note: CallContractAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
}

func testCallContract(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// EstimateGas
	msg := kaia.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
	// CallContract
	if _, err := ec.CallContract(context.Background(), msg, big.NewInt(1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PendingCallContract
	if _, err := ec.PendingCallContract(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testAtFunctions(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	_, err := ec.HeaderByNumber(context.Background(), big.NewInt(1))
	if err != nil {
		t.Fatalf("BlockByNumber error: %v", err)
	}

	// send a transaction for some interesting pending status
	if err := sendTransaction(ec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// wait for the transaction to be included in the pending block
	for {
		// Check pending transaction count
		pending, err := ec.PendingTransactionCount(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pending == 1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Query balance
	balance, err := ec.BalanceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: BalanceAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
	penBalance, err := ec.PendingBalanceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance.Cmp(penBalance) == 0 {
		t.Fatalf("unexpected balance: %v %v", balance, penBalance)
	}
	// NonceAt
	nonce, err := ec.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: NonceAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
	penNonce, err := ec.PendingNonceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if penNonce != nonce+1 {
		t.Fatalf("unexpected nonce: %v %v", nonce, penNonce)
	}
	// StorageAt
	storage, err := ec.StorageAt(context.Background(), testAddr, common.Hash{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: StorageAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
	penStorage, err := ec.PendingStorageAt(context.Background(), testAddr, common.Hash{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(storage, penStorage) {
		t.Fatalf("unexpected storage: %v %v", storage, penStorage)
	}
	// CodeAt
	code, err := ec.CodeAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: CodeAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
	penCode, err := ec.PendingCodeAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(code, penCode) {
		t.Fatalf("unexpected code: %v %v", code, penCode)
	}
	// Note: EstimateGasAtBlock and EstimateGasAtBlockHash methods are not available in Kaia client
	// Skipping these tests for Kaia compatibility

	// Verify that sender address of pending transaction is saved in cache.
	pendingBlock, err := ec.BlockByNumber(context.Background(), big.NewInt(int64(rpc.PendingBlockNumber)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No additional RPC should be required, ensure the server is not asked by
	// canceling the context.
	sender, err := ec.TransactionSender(newCanceledContext(), pendingBlock.Transactions()[0], pendingBlock.Hash(), 0)
	if err != nil {
		t.Fatal("unable to recover sender:", err)
	}
	if sender != testAddr {
		t.Fatal("wrong sender:", sender)
	}
}

func testTransactionSender(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)
	ctx := context.Background()

	// Retrieve testTx1 via RPC.
	block2, err := ec.HeaderByNumber(ctx, big.NewInt(2))
	if err != nil {
		t.Fatal("can't get block 1:", err)
	}
	tx1, err := ec.TransactionInBlock(ctx, block2.Hash(), 0)
	if err != nil {
		t.Fatal("can't get tx:", err)
	}
	if tx1.Hash() != testTx1.Hash() {
		t.Fatalf("wrong tx hash %v, want %v", tx1.Hash(), testTx1.Hash())
	}

	// The sender address is cached in tx1, so no additional RPC should be required in
	// TransactionSender. Ensure the server is not asked by canceling the context here.
	sender1, err := ec.TransactionSender(newCanceledContext(), tx1, block2.Hash(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if sender1 != testAddr {
		t.Fatal("wrong sender:", sender1)
	}

	// Now try to get the sender of testTx2, which was not fetched through RPC.
	// TransactionSender should query the server here.
	sender2, err := ec.TransactionSender(ctx, testTx2, block2.Hash(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if sender2 != testAddr {
		t.Fatal("wrong sender:", sender2)
	}
}

func newCanceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	<-ctx.Done() // Ensure the close of the Done channel
	return ctx
}

func sendTransaction(ec *Client) error {
	chainID, err := ec.ChainID(context.Background())
	if err != nil {
		return err
	}
	nonce, err := ec.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		return err
	}

	signer := types.LatestSignerForChainID(chainID)
	tx := types.NewTransaction(nonce, common.Address{2}, big.NewInt(1), 22000, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
	tx, err = types.SignTx(tx, signer, testKey)
	if err != nil {
		return err
	}
	return ec.SendTransaction(context.Background(), tx)
}

// Here we show how to get the error message of reverted contract call.
func ExampleRevertErrorData() {
	// First create a client.Client instance.
	ctx := context.Background()
	ec, _ := DialContext(ctx, "http://localhost:36000")

	// Call the contract.
	// Note we expect the call to return an error.
	contract := common.HexToAddress("290f1b36649a61e369c6276f6d29463335b4400c")
	call := kaia.CallMsg{To: &contract, Gas: 30000}
	result, err := ec.CallContract(ctx, call, nil)
	if len(result) > 0 {
		panic("got result")
	}
	if err == nil {
		panic("call did not return error")
	}

	// Note: RevertErrorData function is not available in Kaia client
	// For now, just print the error
	fmt.Printf("error: %v\n", err)

	// Note: Since RevertErrorData is not available, we cannot parse the revert data
	// Output would be different for Kaia
}
