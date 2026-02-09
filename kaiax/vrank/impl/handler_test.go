// Copyright 2026 The Kaia Authors
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

package impl

import (
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestVRankModule(t *testing.T) {
	var (
		nodeKey, _ = crypto.GenerateKey()
		candKey, _ = crypto.GenerateKey()
		candAddr   = crypto.PubkeyToAddress(candKey.PublicKey)
		m          = NewVRankModule()
		valset     = mock_valset.NewMockValsetModule(gomock.NewController(t))
	)

	m.Init(&InitOpts{
		Valset:      valset,
		NodeKey:     nodeKey,
		ChainConfig: params.TestKaiaConfig("permissionless"),
	})
	assert.Equal(t, m.prepreparedTime, time.Time{})

	valset.EXPECT().GetCouncil(gomock.Any()).Return([]common.Address{m.nodeId}, nil).AnyTimes()
	valset.EXPECT().GetCandidates(gomock.Any()).Return([]common.Address{candAddr}, nil).AnyTimes()
	valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(m.nodeId, nil).AnyTimes()

	b := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
	m.HandleIstanbulPreprepare(b, &istanbul.View{Round: common.Big0})
	assert.NotEqual(t, m.prepreparedTime, time.Time{})
}
