package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

// InitGenesis new brand genesis
func (keeper Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) error {
	return nil
}

// ExportGenesis returns a GenesisState for a given context.
func (keeper Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return nil
}
