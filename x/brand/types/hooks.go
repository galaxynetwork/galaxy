package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple brand hooks, all hook functions are run in array sequence
var _ BrandHooks = &MultiBrandHooks{}

type MultiBrandHooks []BrandHooks

func NewMultiBrandHooks(hooks ...BrandHooks) MultiBrandHooks {
	return hooks
}

func (hooks MultiBrandHooks) AfterBrandCreated(ctx sdk.Context, brandID string) error {
	for i := range hooks {
		if err := hooks[i].AfterBrandCreated(ctx, brandID); err != nil {
			return err
		}
	}

	return nil
}

func (hooks MultiBrandHooks) AfterBrandOwnerChanged(ctx sdk.Context, brandID string, newOwner sdk.AccAddress, originOwner sdk.AccAddress) error {
	for i := range hooks {
		if err := hooks[i].AfterBrandOwnerChanged(ctx, brandID, newOwner, originOwner); err != nil {
			return err
		}
	}

	return nil
}
