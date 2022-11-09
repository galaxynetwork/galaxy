package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxynetwork/galaxy/x/clairdrop/types"
)

func (k Keeper) GetClaimRecord(ctx sdk.Context, addr sdk.AccAddress) (types.ClaimRecord, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordStorePrefix))

	if !prefixStore.Has(addr) {
		return types.ClaimRecord{}, nil
	}

	bz := prefixStore.Get(addr)
	claimRecord := types.ClaimRecord{}

	err := k.cdc.Unmarshal(bz, &claimRecord)

	if err != nil {
		return types.ClaimRecord{}, err
	}

	return claimRecord, nil

}

func (k Keeper) SetClaimRecord(ctx sdk.Context, claimRecord types.ClaimRecord) error {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordStorePrefix))

	bz, err := k.cdc.Marshal(&claimRecord)

	if err != nil {
		return err
	}

	addr, err := sdk.AccAddressFromBech32(claimRecord.Address)

	if err != nil {
		return err
	}

	prefixStore.Set(addr, bz)
	return nil
}

func (k Keeper) GetClaimRecords(ctx sdk.Context) []types.ClaimRecord {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordStorePrefix))

	iterator := prefixStore.Iterator(
		nil, nil,
	)

	defer iterator.Close()

	claimRecords := []types.ClaimRecord{}

	for ; iterator.Valid(); iterator.Next() {
		claimRecord := types.ClaimRecord{}

		err := k.cdc.Unmarshal(iterator.Value(), &claimRecord)
		if err != nil {
			panic(err)
		}

		claimRecords = append(claimRecords, claimRecord)
	}
	return claimRecords
}

func (k Keeper) SetClaimRecords(ctx sdk.Context, claimRecords []types.ClaimRecord) error {
	for _, claimRecord := range claimRecords {
		err := k.SetClaimRecord(ctx, claimRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) ClearClaimables(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.ClaimRecordStorePrefix))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
	}
}

func (k Keeper) GetClaimableAmountForAction(ctx sdk.Context, addr sdk.AccAddress, action types.ClaimAction) (sdk.Coins, error) {
	claimRecord, err := k.GetClaimRecord(ctx, addr)
	if err != nil {
		return nil, err
	}

	if claimRecord.Address == "" {
		return sdk.Coins{}, nil
	}

	if claimRecord.ActionCompleted[action] {
		return sdk.Coins{}, nil
	}

	params := k.GetParams(ctx)

	if ctx.BlockTime().Before(params.ClairdropStartTime) || ctx.BlockTime().After(params.ClairdropEndTime) {
		return sdk.Coins{}, nil
	}

	claimableCoins := sdk.Coins{}

	for _, coin := range claimRecord.InitalClaimableAmount {
		claimableCoins = claimableCoins.Add(
			sdk.NewCoin(coin.Denom, coin.Amount.QuoRaw(int64(len(types.ClaimAction_name)))),
		)
	}

	return claimableCoins, nil
}

func (k Keeper) GetUserTotalClaimable(ctx sdk.Context, addr sdk.AccAddress) (sdk.Coins, error) {
	claimRecord, err := k.GetClaimRecord(ctx, addr)
	if err != nil {
		return sdk.Coins{}, err
	}
	if claimRecord.Address == "" {
		return sdk.Coins{}, nil
	}

	totalClaimable := sdk.Coins{}

	for action := range types.ClaimAction_name {
		claimableForAction, err := k.GetClaimableAmountForAction(ctx, addr, types.ClaimAction(action))
		if err != nil {
			return sdk.Coins{}, err
		}
		totalClaimable = totalClaimable.Add(claimableForAction...)
	}
	return totalClaimable, nil
}

func (k Keeper) ClaimForAction(ctx sdk.Context, addr sdk.AccAddress, action types.ClaimAction) (sdk.Coins, error) {
	claimableAmount, err := k.GetClaimableAmountForAction(ctx, addr, action)
	if err != nil {
		return claimableAmount, err
	}
	if claimableAmount.Empty() {
		return claimableAmount, nil
	}

	claimRecord, err := k.GetClaimRecord(ctx, addr)

	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, claimableAmount)
	if err != nil {
		return nil, err
	}

	claimRecord.ActionCompleted[action] = true

	err = k.SetClaimRecord(ctx, claimRecord)
	if err != nil {
		return claimableAmount, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, addr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, claimableAmount.String()),
		),
	})

	return claimableAmount, nil
}
