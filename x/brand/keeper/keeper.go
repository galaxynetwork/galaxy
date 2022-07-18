package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

type Keeper struct {
	storeKey    storetypes.StoreKey
	cdc         codec.BinaryCodec
	authKeeper  types.AccountKeeper
	distrKeeper types.DistrKeeper
	paramstore  paramtypes.Subspace
}

// NewKeeper returns a brand keeper. It handles:
// - creating/editing brands
func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, authKeeper types.AccountKeeper, distrKeeper types.DistrKeeper, paramstore paramtypes.Subspace) Keeper {

	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:    storeKey,
		cdc:         cdc,
		authKeeper:  authKeeper,
		distrKeeper: distrKeeper,
		paramstore:  paramstore,
	}
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// getBrandPrefixStore get prefix brands
func (k Keeper) getBrandByOwnerStore(ctx sdk.Context, owner sdk.AccAddress) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetPrefixBrandByOwnerKey(owner))
}
