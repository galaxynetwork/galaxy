package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

func (suite *KeeperTestSuite) TestStoreBrand() {
	require := suite.Require()
	brandKeeper := suite.app.BrandKeeper
	ctx := suite.ctx

	ownerA := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())
	ownerB := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())
	ownerC := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())

	brandA := types.NewBrand("brandida", ownerA, types.NewBrandDescription("name", "", ""))
	brandB := types.NewBrand("brandidb", ownerA, types.NewBrandDescription("name", "", ""))
	brandC := types.NewBrand("brandidc", ownerB, types.NewBrandDescription("name", "", ""))
	invalidBrand := types.NewBrand("_brandidc", ownerB, types.NewBrandDescription("", "", ""))

	//basic validation
	require.NoError(brandA.Validate())
	require.NoError(brandB.Validate())
	require.NoError(brandC.Validate())
	require.Error(invalidBrand.Validate())

	//set brands
	require.NoError(brandKeeper.SetBrand(ctx, brandA))
	require.NoError(brandKeeper.SetBrand(ctx, brandB))
	require.NoError(brandKeeper.SetBrand(ctx, brandC))

	//set brand ids index by onwer
	brandKeeper.SetBrandByOwner(ctx, brandA.Id, ownerA)
	brandKeeper.SetBrandByOwner(ctx, brandB.Id, ownerA)
	brandKeeper.SetBrandByOwner(ctx, brandC.Id, ownerB)

	///has brands
	require.True(brandKeeper.HasBrand(ctx, brandA.Id))
	require.True(brandKeeper.HasBrand(ctx, brandB.Id))
	require.True(brandKeeper.HasBrand(ctx, brandC.Id))
	require.False(brandKeeper.HasBrand(ctx, invalidBrand.Id))

	///get brands
	brand, exist := brandKeeper.GetBrand(ctx, brandA.Id)
	require.True(exist)
	require.Equal(brand.Id, brandA.Id)
	brand, exist = brandKeeper.GetBrand(ctx, brandB.Id)
	require.True(exist)
	require.Equal(brand.Id, brandB.Id)
	brand, exist = brandKeeper.GetBrand(ctx, brandC.Id)
	require.True(exist)
	require.Equal(brand.Id, brandC.Id)
	brand, exist = brandKeeper.GetBrand(ctx, invalidBrand.Id)
	require.False(exist)
	require.Empty(brand)

	require.Len(brandKeeper.GetBrands(ctx), 3)

	require.Len(brandKeeper.GetBrandsByOwner(ctx, ownerA.String()), 2)
	require.Len(brandKeeper.GetBrandsByOwner(ctx, ownerB.String()), 1)
	require.Len(brandKeeper.GetBrandsByOwner(ctx, ownerC.String()), 0)

	// after swap owner of brand
	brand, exist = brandKeeper.GetBrand(ctx, brandA.Id)
	require.Equal(brand.Owner, ownerA.String())

	brandKeeper.DeleteBrandByOwner(ctx, brandA.Id, ownerA)
	brandA.Owner = ownerC.String()
	require.NoError(brandKeeper.SetBrand(ctx, brandA))
	brandKeeper.SetBrandByOwner(ctx, brandA.Id, ownerC)

	brand, exist = brandKeeper.GetBrand(ctx, brandA.Id)
	require.Equal(brand.Owner, ownerC.String())

	require.Len(brandKeeper.GetBrands(ctx), 3)

	require.Len(brandKeeper.GetBrandsByOwner(ctx, ownerA.String()), 1)
	require.Len(brandKeeper.GetBrandsByOwner(ctx, ownerB.String()), 1)
	require.Len(brandKeeper.GetBrandsByOwner(ctx, ownerC.String()), 1)
}
