package clairdrop

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxynetwork/galaxy/x/clairdrop/keeper"
	"github.com/galaxynetwork/galaxy/x/clairdrop/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if genState.Params.ClairdropStartTime.Equal(time.Time{}) {
		genState.Params.ClairdropStartTime = ctx.BlockTime()
		genState.Params.ClairdropEndTime = ctx.BlockTime().Add(time.Hour * 24 * 30 * 8)
	}
	k.SetParams(ctx, genState.Params)
	k.CreateModuleAccount(ctx, genState.ModuleAccountBalance)
	err := k.SetClaimRecords(ctx, genState.ClaimRecords)
	if err != nil {
		panic(
			err,
		)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.ModuleAccountBalance = k.GetModuleAccountBalance(ctx)
	genesis.ClaimRecords = k.GetClaimRecords(ctx)
	return genesis
}
