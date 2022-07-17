package keeper

import (
	"context"

	"github.com/galaxies-labs/galaxy/x/brand/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the x/brand MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateBrand defines a method for creating a new Brand.
func (ms msgServer) CreateBrand(context.Context, *types.MsgCreateBrand) (*types.MsgCreateBrandResponse, error) {
	return nil, nil
}

// EditBrand defines a method for editing an existing brand.
func (ms msgServer) EditBrand(context.Context, *types.MsgEditBrand) (*types.MsgEditBrandResponse, error) {
	return nil, nil
}

// TransferOwnershipBrand defines a method for transfer ownership of existing brand
func (ms msgServer) TransferOwnershipBrand(context.Context, *types.MsgTransferOwnershipBrand) (*types.MsgTransferOwnershipBrandResponse, error) {
	return nil, nil
}
