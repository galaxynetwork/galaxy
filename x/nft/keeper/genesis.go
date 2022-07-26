package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	brandtypes "github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {

	for _, classEntry := range state.ClassEntries {
		if !k.brandKeeper.HasBrand(ctx, classEntry.Class.BrandId) {
			return sdkerrors.Wrapf(brandtypes.ErrNotFoundBrand, "for brandID: %s", classEntry.Class.BrandId)
		}

		if err := k.SaveClass(ctx, classEntry.Class); err != nil {
			return err
		}
		if err := k.setSequenceOfClass(ctx, classEntry.Class.BrandId, classEntry.Class.Id, classEntry.NextSequence); err != nil {
			return err
		}
	}

	for _, entry := range state.Entries {
		owner, _ := sdk.AccAddressFromBech32(entry.Owner)
		for _, nft := range entry.Nfts {

			if nextSequence, err := k.getSequenceOfClass(ctx, nft.BrandId, nft.ClassId); err != nil {
				return err
			} else {
				if nft.Id >= nextSequence {
					return sdkerrors.Wrapf(types.ErrInvalidNFTID, "out of sequence for brandID: %s, classID: %s, id: %d", nft.BrandId, nft.ClassId, nft.Id)
				}
			}

			if exist := k.HasNFT(ctx, nft.BrandId, nft.ClassId, nft.Id); exist {
				return sdkerrors.Wrapf(types.ErrExistNFT, "for brandID: %s, classID: %s, id: %d", nft.BrandId, nft.ClassId, nft.Id)
			}

			if err := k.setNFT(ctx, nft); err != nil {
				return err
			}

			if err := k.incrOnlySupplyOfClass(ctx, nft.BrandId, nft.ClassId); err != nil {
				return err
			}

			k.setOwner(ctx, nft.BrandId, nft.ClassId, nft.Id, owner)
		}
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	classes := k.GetClasses(ctx)

	var classEntries types.ClassEntries
	nftMap := make(map[string][]types.NFT)

	for _, class := range classes {
		seq, err := k.getSequenceOfClass(ctx, class.BrandId, class.Id)
		if err != nil {
			return nil, err
		}

		classEntries = append(classEntries, types.ClassEntry{
			Class:        class,
			NextSequence: seq,
		})
		nfts := k.GetNFTsOfClass(ctx, class.BrandId, class.Id)
		for _, nft := range nfts {
			owner := k.GetOwner(ctx, nft.BrandId, nft.ClassId, nft.Id)
			arr, ok := nftMap[owner.String()]
			if !ok {
				arr = make([]types.NFT, 0)
			}
			nftMap[owner.String()] = append(arr, nft)
		}
	}

	owners := make([]string, 0, len(nftMap))
	for owner := range nftMap {
		owners = append(owners, owner)
	}
	sort.Strings(owners)

	entries := make([]types.Entry, 0, len(nftMap))
	for _, owner := range owners {
		entries = append(entries, types.Entry{
			Owner: owner,
			Nfts:  nftMap[owner],
		})
	}

	return &types.GenesisState{
		ClassEntries: classEntries,
		Entries:      entries,
	}, nil
}
