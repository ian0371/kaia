package contractgov

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	headergov_mock "github.com/kaiachain/kaia/headergov/mocks"
	headergov_types "github.com/kaiachain/kaia/kaiax/headergov/types"
	"github.com/kaiachain/kaia/node/cn/filters"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEffectiveParamSet(t *testing.T) {
	var (
		mockCtrl = gomock.NewController(t)
		chain    = mocks.NewMockBlockChain(mockCtrl)
		mockHGM  = new(headergov_mock.MockHeaderGovModule)
	)

	cgm := &ContractGovModule{}
	cgm.Init(&InitOpts{
		Chain:       chain,
		ChainConfig: &params.ChainConfig{},
		hgm:         mockHGM,
	})

	testCases := []struct {
		name           string
		blockNum       uint64
		koreForkActive bool
		contractAddr   common.Address
		expectedParams ParamSet
		expectedError  error
	}{
		{
			name:           "Kore fork not active",
			blockNum:       100,
			koreForkActive: false,
			expectedError:  errContractEngineNotReady,
		},
		{
			name:           "Kore fork active, contract address not set",
			blockNum:       200,
			koreForkActive: true,
			contractAddr:   common.Address{},
			expectedParams: ParamSet{},
			expectedError:  nil,
		},
		{
			name:           "Kore fork active, contract address set",
			blockNum:       300,
			koreForkActive: true,
			contractAddr:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedParams: ParamSet{
				MaxValidators: 100,
				MinStake:      big.NewInt(1000000),
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cgm.ChainConfig.KoreForkBlock = big.NewInt(150)

			chain.EXPECT().Config().Return(cgm.ChainConfig)
			mockHGM.EXPECT().EffectiveParams(tc.blockNum).Return(headergov_types.ParamSet{GovParamContract: tc.contractAddr}, nil)

			if tc.koreForkActive && !tc.contractAddr.IsZero() {
				chain.On("GetHeaderByNumber", tc.blockNum).Return(&types.Header{})
				chain.On("State").Return(&state.StateDB{}, nil)

				mockBackend := new(mockContractBackend)
				mockBackend.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return([]byte{}, nil)
				backends.NewBlockchainContractBackend = func(chain backends.BlockChainForCaller, tp backends.TxPoolForCaller, es *filters.EventSystem) *backends.BlockchainContractBackend {
					return mockBackend
				}
			}

			params, err := cgm.EffectiveParamSet(tc.blockNum)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedParams, params)
			}

			chain.AssertExpectations(t)
			mockHGM.AssertExpectations(t)
		})
	}
}
