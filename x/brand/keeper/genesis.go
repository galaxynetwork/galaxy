package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

// InitGenesis new brand genesis
func (keeper Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) error {
	keeper.SetParams(ctx, genState.Params)
	for _, brand := range genState.Brands {
		if len(strings.TrimSpace(brand.Id)) == 0 {
			brand.BrandAddress = types.NewBrandAddress(brand.Id).String()
		}
		if err := brand.Validate(); err != nil {
			return err
		}

		keeper.SetBrand(ctx, brand)

		acc, _ := sdk.AccAddressFromBech32(brand.Owner)
		keeper.SetBrandByOwner(ctx, brand.Id, acc)
	}
	return nil
}

// ExportGenesis returns a GenesisState for a given context.
func (keeper Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	var genState types.GenesisState
	genState.Params = keeper.GetParams(ctx)
	genState.Brands = keeper.GetBrands(ctx)
	return genState
}
