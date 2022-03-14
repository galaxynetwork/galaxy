package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
	SetModuleAccount(sdk.Context, types.ModuleAccountI)
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
}

type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

type DistributionKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
