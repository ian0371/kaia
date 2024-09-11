package contractgov

import (
	"math/big"
	"testing"

	oldgomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	govcontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/gov"
	"github.com/kaiachain/kaia/crypto"
	mock_headergov "github.com/kaiachain/kaia/kaiax/headergov/mocks"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createSimulateBackend(t *testing.T) ([]*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *govcontract.GovParam) {
	// Create accounts and simulated blockchain
	accounts := []*bind.TransactOpts{}
	alloc := blockchain.GenesisAlloc{}
	for i := 0; i < 1; i++ {
		key, _ := crypto.GenerateKey()
		account := bind.NewKeyedTransactor(key)
		account.GasLimit = 10000000
		accounts = append(accounts, account)
		alloc[account.From] = blockchain.GenesisAccount{Balance: big.NewInt(params.KAIA)}
	}
	config := &params.ChainConfig{}
	config.SetDefaults()
	config.UnitPrice = 25e9
	config.IstanbulCompatibleBlock = common.Big0
	config.LondonCompatibleBlock = common.Big0
	config.EthTxTypeCompatibleBlock = common.Big0
	config.MagmaCompatibleBlock = common.Big0
	config.KoreCompatibleBlock = common.Big0

	sim := backends.NewSimulatedBackendWithDatabase(database.NewMemoryDBManager(), alloc, config)

	// Deploy contract
	owner := accounts[0]
	address, tx, contract, err := govcontract.DeployGovParam(owner, sim)
	require.Nil(t, err)
	sim.Commit()

	receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	return accounts, sim, address, contract
}

func prepareContractGovModule(t *testing.T, bc *blockchain.BlockChain, addr common.Address) *contractGovModule {
	mockHGM := mock_headergov.NewMockHeaderGovModule(gomock.NewController(t))
	cgm := &contractGovModule{}
	cgm.Init(&InitOpts{
		Chain:       bc,
		ChainConfig: &params.ChainConfig{KoreCompatibleBlock: big.NewInt(100)},
		hgm:         mockHGM,
	})
	mockHGM.EXPECT().EffectiveParamSet(oldgomock.Any()).Return(ParamSet{GovParamContract: addr}, nil).AnyTimes()
	return cgm
}

func TestEffectiveParamSet(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlError)
	accounts, sim, addr, gp := createSimulateBackend(t)
	cgm := prepareContractGovModule(t, sim.BlockChain(), addr)

	{
		activation := big.NewInt(1000)
		val := []byte{0, 0, 0, 0, 0, 0, 0, 25}
		tx, err := gp.SetParam(accounts[0], "governance.unitprice", true, val, activation)
		require.Nil(t, err)
		sim.Commit()

		receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

		pset, err := cgm.EffectiveParamSet(1000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(25), pset.UnitPrice)
	}

	{
		activation := big.NewInt(2000)
		val := []byte{0, 0, 0, 0, 0, 0, 0, 125}
		tx, err := gp.SetParam(accounts[0], "governance.unitprice", true, val, activation)
		require.Nil(t, err)
		sim.Commit()

		receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

		pset, err := cgm.EffectiveParamSet(2000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(125), pset.UnitPrice)
	}
}
