package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/mint/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramStore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramStore.SetParamSet(ctx, &params)
}
