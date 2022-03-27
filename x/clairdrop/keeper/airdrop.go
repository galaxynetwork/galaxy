package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/clairdrop/types"
)

// SetParams set the params
func (k Keeper) EndAirdrop(ctx sdk.Context) error {

	err := k.ClawbackAirdrop(ctx)
	if err != nil {
		return err
	}

	err = k.FundRemainingsToCommunity(ctx)
	if err != nil {
		return err
	}

	k.ClearClaimables(ctx)

	return nil
}

func (k Keeper) ClawbackAirdrop(ctx sdk.Context) error {
	claimRecords := k.GetClaimRecords(ctx)
	developerAddress := k.mk.GetDeveloperAddress(ctx)
	for _, claimRecord := range claimRecords {
		addr, err := sdk.AccAddressFromBech32(claimRecord.Address)
		if err != nil {
			return err
		}

		acc := k.ak.GetAccount(ctx, addr)
		if acc == nil {
			continue
		}

		seq, err := k.ak.GetSequence(ctx, addr)
		if err != nil {
			return err
		}
		//if never make transaction
		if seq == 0 {
			//skip developer
			var skip bool
			for _, addr := range developerAddress {
				if addr == claimRecord.Address {
					skip = true
					continue
				}
			}
			if skip {
				continue
			}
			balance := k.bk.GetBalance(ctx, addr, types.DefaultClaimDenom)
			clawbackCoins := sdk.NewCoins(balance)
			err = k.dk.FundCommunityPool(ctx, clawbackCoins, addr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) FundRemainingsToCommunity(ctx sdk.Context) error {
	moduleAccAddr := k.ak.GetModuleAddress(types.ModuleName)
	amt := k.GetModuleAccountBalance(ctx)
	return k.dk.FundCommunityPool(ctx, sdk.NewCoins(amt), moduleAccAddr)
}
