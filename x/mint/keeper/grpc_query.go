package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/mint/types"
)

func (k Keeper) Minter(c context.Context, _ *types.QueryMinterRequest) (*types.QueryMinterResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := k.GetMinter(ctx)

	return &types.QueryMinterResponse{Minter: minter}, nil
}

func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}
