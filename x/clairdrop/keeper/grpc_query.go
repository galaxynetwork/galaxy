package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/clairdrop/types"
)

func (k Keeper) ModuleAccountBalance(c context.Context, _ *types.QueryModuleAccountBalanceRequest) (*types.QueryModuleAccountBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	balance := k.GetModuleAccountBalance(ctx)
	return &types.QueryModuleAccountBalanceResponse{ModuleAccountBalance: sdk.NewCoins(balance)}, nil
}

func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) TotalClaimable(c context.Context, req *types.QueryTotalClaimableRequest) (*types.QueryTotalClaimableResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	coins, err := k.GetUserTotalClaimable(ctx, addr)

	return &types.QueryTotalClaimableResponse{
		Coins: coins,
	}, err
}

func (k Keeper) ClaimRecord(c context.Context, req *types.QueryClaimRecordRequest) (*types.QueryClaimRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	claimRecord, err := k.GetClaimRecord(ctx, addr)

	return &types.QueryClaimRecordResponse{
		ClaimRecord: claimRecord,
	}, err
}

func (k Keeper) ClaimableForAction(c context.Context, req *types.QueryClaimableForActionRequest) (*types.QueryClaimableForActionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	coins, err := k.GetClaimableAmountForAction(ctx, addr, req.Action)

	return &types.QueryClaimableForActionResponse{
		Coins: coins,
	}, err
}
