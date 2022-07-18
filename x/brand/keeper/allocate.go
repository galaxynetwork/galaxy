package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AllocateFunds Not used yet.
func (keeper Keeper) AllocateFunds(ctx sdk.Context) error {
	brands := keeper.GetBrands(ctx)

	for _, brand := range brands {
		addr, _ := sdk.AccAddressFromBech32(brand.BrandAddress)
		ownerAddr, _ := sdk.AccAddressFromBech32(brand.Owner)

		coins := keeper.bankKeeper.GetAllBalances(ctx, addr)

		if err := keeper.bankKeeper.SendCoins(ctx, addr, ownerAddr, coins); err != nil {
			return fmt.Errorf("failed to allocate brand assets to owner: %s", err)
		}
	}

	return nil
}
