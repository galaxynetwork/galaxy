package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/galaxies-labs/galaxy/x/brand/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
	require := suite.Require()
	brandKeeper := suite.app.BrandKeeper
	ctx := suite.ctx

	ownerA := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())
	ownerB := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())

	brandA := types.NewBrand("brandida", ownerA, types.NewBrandDescription("name", "", ""))
	brandB := types.NewBrand("brandidb", ownerA, types.NewBrandDescription("name", "", ""))
	brandC := types.NewBrand("brandidc", ownerB, types.NewBrandDescription("name", "", ""))
	invalidBrand := types.NewBrand("_brandidc", ownerB, types.NewBrandDescription("", "", ""))

	genStateA := types.NewGenesisState(types.Brands{brandA, brandB, brandC}, types.DefaultParams())
	invalidGenStateA := types.NewGenesisState(types.Brands{invalidBrand}, types.DefaultParams())
	invalidGenStateB := types.NewGenesisState(types.Brands{brandA, brandB, brandC}, types.NewParams(sdk.NewCoin("ustake", sdk.NewInt(0))))

	require.NoError(types.ValidateGenesis(*genStateA))
	require.Error(types.ValidateGenesis(*invalidGenStateA))
	require.Error(types.ValidateGenesis(*invalidGenStateB))

	require.NoError(brandKeeper.InitGenesis(ctx, *genStateA))
	require.Error(brandKeeper.InitGenesis(ctx, *invalidGenStateA))
	//panic when initial invalid parameter space set
	require.Panics(func() {
		brandKeeper.InitGenesis(ctx, *invalidGenStateB)
	})

	exportedGenStateA := brandKeeper.ExportGenesis(ctx)
	require.Equal(genStateA, &exportedGenStateA)
	require.NotEqual(invalidGenStateA, &exportedGenStateA)

	require.NoError(brandKeeper.InitGenesis(ctx, exportedGenStateA))
}
