package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
}

type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

type DistrKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

// Event Hooks
// These can be utilized to communicate between a brand keeper and another
// keeper which must take particular actions when brand change
// state. The second keeper must implement this interface, which then the
// staking keeper can call.

// BrandHooks event hooks for brand object (noalias)
type BrandHooks interface {
	AfterBrandCreated(ctx sdk.Context, brandID string) error                                                           // Must be called when a brand is created
	AfterBrandOwnerChanged(ctx sdk.Context, brandID string, newOwner sdk.AccAddress, originOwner sdk.AccAddress) error // Must be called when a brand owner is changed
}

// BrandHooksWrapper is a wrapper for modules to inject BrandHooks using depinject.
type BrandHooksWrapper struct{ BrandHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (BrandHooksWrapper) IsOnePerModuleType() {}
