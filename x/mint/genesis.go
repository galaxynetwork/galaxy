package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxynetwork/galaxy/x/mint/keeper"
	"github.com/galaxynetwork/galaxy/x/mint/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, ak types.AccountKeeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	k.SetMinter(ctx, genState.Minter)
	ak.GetModuleAccount(ctx, types.ModuleName)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.Minter = k.GetMinter(ctx)
	return genesis
}
