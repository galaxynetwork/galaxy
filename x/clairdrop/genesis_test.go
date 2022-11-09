package clairdrop_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxynetwork/galaxy/app"
	"github.com/galaxynetwork/galaxy/x/clairdrop"
	"github.com/galaxynetwork/galaxy/x/clairdrop/types"
	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var now = time.Now().UTC()
var acc1 = sdk.AccAddress([]byte("addr1~"))
var acc2 = sdk.AccAddress([]byte("addr2~"))
var claimRecords = []types.ClaimRecord{
	{
		Address:               acc1.String(),
		InitalClaimableAmount: sdk.Coins{sdk.NewInt64Coin(types.DefaultClaimDenom, 1_000_000)},
		ActionCompleted:       []bool{true, false, true, true},
	},
	{
		Address:               acc2.String(),
		InitalClaimableAmount: sdk.Coins{sdk.NewInt64Coin(types.DefaultClaimDenom, 4_000_000)},
		ActionCompleted:       []bool{false, false, false, false},
	},
}

func TestClaimInitGenesis(t *testing.T) {

	tests := []types.GenesisState{
		{
			ModuleAccountBalance: sdk.NewInt64Coin(types.DefaultClaimDenom, 1_000_000+4_000_000+1),
			Params: types.Params{
				ClairdropStartTime: now,
				ClairdropEndTime:   now.Add(time.Hour * 3),
			},
			ClaimRecords: claimRecords,
		},
		{
			ModuleAccountBalance: sdk.NewInt64Coin(types.DefaultClaimDenom, 1_000_000+4_000_000),
			Params: types.Params{
				ClairdropStartTime: now,
				ClairdropEndTime:   now.Add(time.Hour * 3),
			},
			ClaimRecords: claimRecords,
		}, {
			ModuleAccountBalance: sdk.NewInt64Coin(types.DefaultClaimDenom, 1_000_000+4_000_000),
			Params: types.Params{
				ClairdropStartTime: time.Time{},
				ClairdropEndTime:   time.Time{},
			},
			ClaimRecords: claimRecords,
		},
	}
	for i, genesis := range tests {
		err := types.ValidateGenesis(genesis)
		if i != 0 {
			require.NoError(t, err)
			continue
		}

		app := app.Setup(false)
		ctx := app.BaseApp.NewContext(false, tmproto.Header{})
		ctx = ctx.WithBlockTime(now.Add(time.Second))

		clairdrop.InitGenesis(ctx, app.ClairdropKeeper, genesis)

		coin := app.ClairdropKeeper.GetModuleAccountBalance(ctx)
		require.Equal(t, coin.String(), genesis.ModuleAccountBalance.String())

		params := app.ClairdropKeeper.GetParams(ctx)
		if genesis.Params.ClairdropStartTime.IsZero() {
			require.NotEqual(t, params, genesis.Params)
		} else {
			require.Equal(t, params, genesis.Params)

		}

		claimRecords := app.ClairdropKeeper.GetClaimRecords(ctx)
		require.Equal(t, claimRecords, genesis.ClaimRecords)
	}
}
