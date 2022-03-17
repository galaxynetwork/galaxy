package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/galaxies-labs/galaxy/x/clairdrop/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramStore paramtypes.Subspace

	ak types.AccountKeeper
	bk types.BankKeeper
	dk types.DistributionKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	dk types.DistributionKeeper,

) Keeper {

	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramStore: ps,
		ak:         ak,
		bk:         bk,
		dk:         dk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetModuleAccountBalance(ctx sdk.Context) sdk.Coin {
	moduleAccAddr := k.ak.GetModuleAddress(types.ModuleName)
	return k.bk.GetBalance(ctx, moduleAccAddr, types.DefaultClaimDenom)
}

func (k Keeper) CreateModuleAccount(ctx sdk.Context, amount sdk.Coin) {
	moduleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)

	k.ak.SetModuleAccount(ctx, moduleAcc)

	mintCoins := sdk.NewCoins(amount)

	existingModuleAcctBalance := k.bk.GetBalance(
		ctx,
		k.ak.GetModuleAddress(types.ModuleName),
		amount.Denom,
	)

	if existingModuleAcctBalance.IsPositive() {

		actual := existingModuleAcctBalance.Add(amount)

		ctx.Logger().Info(fmt.Sprintf(
			"WARNING! There is a bug in claims on InitGenesis, that you are subject to."+
				" You likely expect the claims module account balance to be %d %s, but it will actually be %d %s due to this bug.",
			amount.Amount.Int64(), amount.Denom, actual.Amount.Int64(), actual.Denom))
	}

	if err := k.bk.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		panic(err)
	}

}
