package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

// GetTotalSupplyOfClass returns the totalSupply information of the specified class
func (k Keeper) GetTotalSupplyOfClass(ctx sdk.Context, brandID, id string) (uint64, error) {
	bz := k.getClassSupplyStore(ctx, brandID).Get([]byte(id))
	if bz == nil {
		return 0, fmt.Errorf("invalid class dereference")
	}

	var supply types.Supply
	if err := k.cdc.Unmarshal(bz, &supply); err != nil {
		return 0, err
	}

	return supply.TotalSupply, nil
}

func (k Keeper) incrSupplyOfClass(ctx sdk.Context, brandID, id string) error {
	store := k.getClassSupplyStore(ctx, brandID)
	bz := store.Get([]byte(id))
	if bz == nil {
		return fmt.Errorf("invalid class dereference")
	}

	var supply types.Supply
	if err := k.cdc.Unmarshal(bz, &supply); err != nil {
		return err
	}

	supply.IncreaseSupply()

	if bz, err := k.cdc.Marshal(&supply); err != nil {
		return err
	} else {
		store.Set([]byte(id), bz)
	}

	return nil
}

func (k Keeper) decrSupplyOfClass(ctx sdk.Context, brandID, id string) error {
	store := k.getClassSupplyStore(ctx, brandID)
	bz := store.Get([]byte(id))
	if bz == nil {
		return fmt.Errorf("invalid class dereference")
	}

	var supply types.Supply
	if err := k.cdc.Unmarshal(bz, &supply); err != nil {
		return err
	}

	supply.DecreaseSupply()

	if bz, err := k.cdc.Marshal(&supply); err != nil {
		return err
	} else {
		store.Set([]byte(id), bz)
	}

	return nil
}

func (k Keeper) getSequenceOfClass(ctx sdk.Context, brandID, id string) (uint64, error) {
	bz := k.getClassSupplyStore(ctx, brandID).Get([]byte(id))

	if bz == nil {
		return 0, fmt.Errorf("invalid class dereference")
	}

	var supply types.Supply
	if err := k.cdc.Unmarshal(bz, &supply); err != nil {
		return 0, err
	}

	return supply.Sequence, nil
}

func (k Keeper) getClassSupplyStore(ctx sdk.Context, brandID string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.GetClassSupplyStoreKey(brandID))
}

func (k Keeper) initializeClassSupply(ctx sdk.Context, brandID, id string) error {
	supply := types.DefaultSupply()

	bz, err := k.cdc.Marshal(&supply)
	if err != nil {
		return err
	}

	k.getClassSupplyStore(ctx, brandID).Set([]byte(id), bz)

	return nil
}
