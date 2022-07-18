package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

// Implements BrandHooks interface
var (
	_ types.BrandHooks = Keeper{}
)

func (keeper Keeper) AfterBrandCreated(ctx sdk.Context, brandID string) error {
	if keeper.hooks != nil {
		return keeper.hooks.AfterBrandCreated(ctx, brandID)
	}
	return nil
}

func (keeper Keeper) AfterBrandOwnerChanged(ctx sdk.Context, brandID string, newOwner sdk.AccAddress, originOwner sdk.AccAddress) error {
	if keeper.hooks != nil {
		return keeper.hooks.AfterBrandOwnerChanged(ctx, brandID, newOwner, originOwner)
	}
	return nil
}
