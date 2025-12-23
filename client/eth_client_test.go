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
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain/types"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IEthClient interface {
	Close()
	SetHeader(key, value string)
	// BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	// BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	// HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	// HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*EthBlock, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*EthBlock, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*EthHeader, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*EthHeader, error)
	// TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *api.EthRPCTransaction, isPending bool, err error)
	TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	// TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*api.EthRPCTransaction, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionReceiptRpcOutput(ctx context.Context, txHash common.Hash) (r map[string]interface{}, err error)
	// SyncProgress(ctx context.Context) (*kaia.SyncProgress, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (kaia.Subscription, error)
	// AuctionSubscribeNewHead(ctx context.Context, ch chan<- *types.Header)
	NetworkID(ctx context.Context) (*big.Int, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	FilterLogs(ctx context.Context, q kaia.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q kaia.FilterQuery, ch chan<- types.Log) (kaia.Subscription, error)
	// AuctionSubscribeFilterLogs(ctx context.Context, q kaia.FilterQuery, ch chan<- types.Log)
	PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
	PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	PendingTransactionCount(ctx context.Context) (uint, error)
	// AuctionSubscribeFullPendingTransactions(ctx context.Context, ch chan<- *types.Transaction)
	// AuctionSubscribeFullPendingTransactionsRaw(ctx context.Context, ch chan<- map[string]any)
	// AuctionSubscribePendingTransactions(ctx context.Context, ch chan<- common.Hash)
	CallContract(ctx context.Context, msg kaia.CallMsg, blockNumber *big.Int) ([]byte, error)
	// AuctionCallContract(ctx context.Context, msg kaia.CallMsg, blockNumber *big.Int)
	PendingCallContract(ctx context.Context, msg kaia.CallMsg) ([]byte, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg kaia.CallMsg) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	// SendAuctionTx(ctx context.Context, bidInput auction_impl.BidInput)
	SendRawTransaction(ctx context.Context, tx *types.Transaction) (common.Hash, error)
	// SendUnsignedTransaction(ctx context.Context, unsignedTx api.SendTxArgs) (common.Hash, error)
	// ImportRawKey(ctx context.Context, key string, password string) (common.Address, error)
	// UnlockAccount(ctx context.Context, address common.Address, password string, time uint)
	BlockNumber(ctx context.Context) (*big.Int, error)
	ChainID(ctx context.Context) (*big.Int, error)
	// AddPeer(ctx context.Context, url string) (bool, error)
	// RemovePeer(ctx context.Context, url string) (bool, error)
	CreateAccessList(ctx context.Context, msg kaia.CallMsg) (*types.AccessList, uint64, string, error)
}

// Verify that EthClient implements the Kaia interfaces.
var (
	// _ = kaia.Subscription(&EthClient{})
	// _ = kaia.ChainReader(&EthClient{}) // returns EthBlock thus not compatible
	// _ = kaia.TransactionReader(&EthClient{})
	_ = kaia.ChainStateReader(&EthClient{})
	// _ = kaia.ChainSyncReader(&EthClient{})
	_ = kaia.ContractCaller(&EthClient{})
	_ = kaia.LogFilterer(&EthClient{})
	_ = kaia.TransactionSender(&EthClient{})
	_ = kaia.GasPricer(&EthClient{})
	_ = kaia.PendingStateReader(&EthClient{})
	_ = kaia.PendingContractCaller(&EthClient{})
	_ = kaia.GasEstimator(&EthClient{})
	_ = bind.ContractBackend(&EthClient{})
	_ = bind.DeployBackend(&EthClient{})
	_ = IEthClient(&EthClient{})
	// _ = kaia.PendingStateEventer(&Client{})

	anvilRichKey, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	anvilRichAddr   = crypto.PubkeyToAddress(anvilRichKey.PublicKey)
)

type AnvilTestSuite struct {
	suite.Suite
	ethClient *EthClient
	serverURL string
	cmd       *exec.Cmd
	cancel    func()

	tx *types.Transaction
}

func (s *AnvilTestSuite) SetupSuite() {
	var err error
	var nonce uint64
	defer func() {
		if err != nil {
			s.TearDownSuite()
			s.T().Skip(err)
			return
		}
		s.T().Logf("Anvil server started on %s", s.serverURL)
	}()

	if err = s.launchAnvilServer(); err != nil {
		return
	}
	if s.ethClient, err = tryConnectEth(s.serverURL); err != nil {
		return
	}
	if nonce, err = s.ethClient.NonceAt(context.Background(), anvilRichAddr, nil); err != nil {
		return
	}

	unsignedTx := types.NewTransaction(nonce, testAddr, big.NewInt(1e18), params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
	signer := types.LatestSignerForChainID(genesisConfig.ChainID)
	s.tx, _ = types.SignTx(unsignedTx, signer, anvilRichKey)
	if _, err = s.ethClient.SendRawTransaction(context.Background(), s.tx); err != nil {
		return
	}
	time.Sleep(1 * time.Second)
	s.T().Logf("Signed transaction: %s", s.tx.Hash().Hex())
}

func (s *AnvilTestSuite) TearDownSuite() {
	if s.ethClient != nil {
		s.ethClient.Close()
	}
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *AnvilTestSuite) launchAnvilServer() error {
	randPort := rand.Intn(30000) + 20000
	// Check if anvil is installed
	_, err := exec.LookPath("anvil")
	if err != nil {
		return errors.New("anvil not found in PATH")
	}

	// Start anvil in background
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.cmd = exec.CommandContext(ctx, "anvil", "--chain-id", "1337", "--host", "127.0.0.1", "--port", strconv.Itoa(randPort))
	err = s.cmd.Start()
	if err != nil {
		return errors.New("failed to start anvil")
	}

	// Give anvil time to start up
	time.Sleep(2 * time.Second)

	s.serverURL = fmt.Sprintf("http://127.0.0.1:%d", randPort)
	return nil
}

func TestAnvilTestSuite(t *testing.T) {
	suite.Run(t, new(AnvilTestSuite))
}

func (s *AnvilTestSuite) TestBlockchainAccess() {
	// BlockByHash, BlockByNumber, HeaderByHash, HeaderByNumber,
	// TransactionByHash, TransactionSender, TransactionCount, TransactionInBlock, TransactionReceipt, TransactionReceiptRpcOutput
	ethclient := s.ethClient

	block, err := ethclient.BlockByNumber(context.Background(), big.NewInt(1))
	require.NoError(s.T(), err)
	assert.Equal(s.T(), uint64(1), block.Header().Number.Uint64())

	block, err = ethclient.BlockByHash(context.Background(), block.Header().Hash())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), uint64(1), block.Header().Number.Uint64())

	bn, err := ethclient.BlockNumber(context.Background())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), uint64(1), bn.Uint64())

	header, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(1))
	require.NoError(s.T(), err)
	assert.Equal(s.T(), uint64(1), header.Number.Uint64())

	apiTx, _, err := ethclient.TransactionByHash(context.Background(), s.tx.Hash())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), s.tx.Hash(), apiTx.Hash)

	cnt, err := ethclient.TransactionCount(context.Background(), header.Hash())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), uint(1), cnt)

	apiTx, err = ethclient.TransactionInBlock(context.Background(), header.Hash(), 0)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), apiTx.Hash, apiTx.Hash)

	receipt, err := ethclient.TransactionReceipt(context.Background(), s.tx.Hash())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), types.ReceiptStatusSuccessful, receipt.Status)

	receiptMap, err := ethclient.TransactionReceiptRpcOutput(context.Background(), s.tx.Hash())
	require.NoError(s.T(), err)

	status, err := strconv.ParseUint(receiptMap["status"].(string), 0, 64)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), types.ReceiptStatusSuccessful, uint(status), "tx %s failed", s.tx.Hash().Hex())
}

func (s *AnvilTestSuite) TestBalanceAt() {
	// balance, err := ethclient.BalanceAt(context.Background(), testAddr, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// assert.True(t, balance.Cmp(testAddrInitialBalance) >= 0)

	// nonce, err := ethclient.NonceAt(context.Background(), testAddr, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// dynamicTx := func() *types.Transaction {
	// 	tx := types.NewTx(&types.TxInternalDataEthereumDynamicFee{
	// 		ChainID:      genesisConfig.ChainID,
	// 		AccountNonce: nonce,
	// 		Recipient:    &testAddr,
	// 		Amount:       big.NewInt(10),
	// 		GasLimit:     25000,
	// 		GasFeeCap:    big.NewInt(50e9),
	// 		GasTipCap:    big.NewInt(25e9),
	// 	})
	// 	signer := types.LatestSignerForChainID(genesisConfig.ChainID)
	// 	signedTx, _ := types.SignTx(tx, signer, testKey)
	// 	return signedTx
	// }()
	// assert.Equal(t, header.ParentHash, common.Hash{})
	// assert.NotEqual(t, header.Hash(), common.Hash{})
	// _, err = ethclient.SendRawTransaction(context.Background(), dynamicTx)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// nonce++
	// deployTx := func() *types.Transaction {
	// 	// contract Storage { uint256 number = 1337; * @dev Return value @return value of 'number' */ function retrieve() public view returns (uint256){ return number; } }
	// 	bytecode := common.Hex2Bytes("60806040526105395f553480156013575f5ffd5b5060af80601f5f395ff3fe6080604052348015600e575f5ffd5b50600436106026575f3560e01c80632e64cec114602a575b5f5ffd5b60306044565b604051603b91906062565b60405180910390f35b5f5f54905090565b5f819050919050565b605c81604c565b82525050565b5f60208201905060735f8301846055565b9291505056fea2646970667358221220bbed5c2a1719068dca0cf4da53d280029c147463a0b8f8319bc3494906ad27a964736f6c634300081e0033")
	// 	tx := types.NewContractCreation(nonce, big.NewInt(0), 1e6, big.NewInt(25e9), bytecode)
	// 	signer := types.LatestSignerForChainID(genesisConfig.ChainID)
	// 	signedTx, _ := types.SignTx(tx, signer, testKey)
	// 	return signedTx
	// }()
	// hash, err = ethclient.SendRawTransaction(context.Background(), deployTx)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// time.Sleep(1 * time.Second)
	// receipt, err = ethclient.TransactionReceiptRpcOutput(context.Background(), hash)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// status, err = strconv.ParseUint(receipt["status"].(string), 0, 64)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// assert.Equal(t, types.ReceiptStatusSuccessful, uint(status), "tx %s failed", hash.Hex())

	// contractAddr := crypto.CreateAddress(testAddr, nonce)
	// calldata := common.Hex2Bytes("2e64cec1") // retrieve()(uint256)
	// ret, err := ethclient.CallContract(context.Background(), kaia.CallMsg{To: &contractAddr, Data: calldata}, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000539", common.Bytes2Hex(ret))

	// accesslist, _, _, err := ethclient.CreateAccessList(context.Background(), kaia.CallMsg{To: &contractAddr, Data: calldata})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// assert.Equal(t, 1, accesslist.StorageKeys())

	// storage, err := ethclient.StorageAt(context.Background(), contractAddr, (*accesslist)[0].StorageKeys[0], nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000539", common.Bytes2Hex(storage))
}

// func TestEthClient_AnvilServerWithCleanup(t *testing.T) {
// 	// Launch anvil server with 30s timeout
// 	serverURL, cleanup := launchAnvilServer(t)
// 	defer cleanup() // Ensure cleanup happens when test ends

// 	// Connect to anvil server
// 	ethclient, err := tryConnectEth(serverURL)
// 	if err != nil {
// 		t.Skip("Could not connect Eth client to anvil server:", err)
// 		return
// 	}
// 	defer ethclient.Close()

// 	// Test basic functionality
// 	ethHeader, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0))
// 	if err != nil {
// 		t.Fatal("Failed to get genesis block:", err)
// 	}

// 	t.Log("Genesis block hash:", ethHeader.Hash().Hex())
// 	assert.Equal(t, uint64(0), ethHeader.Number.Uint64())

// 	// Test contract deployment
// 	bytecode := common.Hex2Bytes("608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029")

// 	// Use anvil's default funded account
// 	deployTx := types.NewContractCreation(0, big.NewInt(0), 1000000, big.NewInt(1e9), bytecode)
// 	signer := types.LatestSignerForChainID(big.NewInt(1337))
// 	signedDeployTx, err := types.SignTx(deployTx, signer, anvilRichKey)
// 	if err != nil {
// 		t.Fatal("Failed to sign deploy tx:", err)
// 	}

// 	hash, err := ethclient.SendRawTransaction(context.Background(), signedDeployTx)
// 	if err != nil {
// 		t.Fatal("Failed to deploy contract:", err)
// 	}

// 	t.Log("Contract deployed with tx hash:", hash.Hex())
// }

func (s *AnvilTestSuite) TestKaiaClient() {
	kaiaClient, err := tryConnect(s.serverURL)
	if err != nil {
		s.T().Skip("Could not connect Kaia client to anvil server:", err)
		return
	}
	defer kaiaClient.Close()

	_, err = kaiaClient.HeaderByNumber(context.Background(), big.NewInt(0))
	assert.Equal(s.T(), err.Error(), "Method not found")
}
