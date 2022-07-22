package keeper_test

import "github.com/galaxies-labs/galaxy/x/nft/types"

func (suite *KeeperTestSuite) TestClass() {
	keeper, ctx := suite.app.NFTKeeper, suite.ctx

	brandId := "brandid"
	brandId2 := "brandId2"

	classA := types.NewClass(brandId, "classsid", 10_000, types.NewClassDescription("", "", "", ""))
	classA2 := types.NewClass(brandId, "classsid", 10_000, types.NewClassDescription("", "", "", ""))
	classB := types.NewClass(brandId, "classsid2", 10_000, types.NewClassDescription("", "", "", ""))
	classC := types.NewClass(brandId2, "classsid", 10_000, types.NewClassDescription("", "", "", ""))
	unsavedClassA := types.NewClass("random", "random", 10_000, types.NewClassDescription("", "", "", ""))

	//save
	suite.Require().NoError(keeper.SaveClass(ctx, classA))
	suite.Require().NoError(keeper.SaveClass(ctx, classB))
	suite.Require().NoError(keeper.SaveClass(ctx, classC))
	suite.Require().Error(keeper.SaveClass(ctx, classA2))
	suite.Require().Error(keeper.SaveClass(ctx, classB))

	//set
	suite.Require().NoError(keeper.SetClass(ctx, classA))
	suite.Require().NoError(keeper.SetClass(ctx, classB))
	suite.Require().NoError(keeper.SetClass(ctx, classC))

	//get
	savedClasses := types.Classes{classA, classB, classC}
	for _, class := range savedClasses {
		class, exist := keeper.GetClass(ctx, class.BrandId, class.Id)
		suite.Require().True(exist)
		suite.Require().NotNil(class)
		suite.Require().NotEmpty(class)
	}

	class, exist := keeper.GetClass(ctx, unsavedClassA.BrandId, unsavedClassA.Id)
	suite.Require().False(exist)
	suite.Require().NotNil(class)
	suite.Require().Empty(class)

	//has
	for _, class := range savedClasses {
		suite.Require().True(keeper.HasClass(ctx, class.BrandId, class.Id))
	}
	suite.Require().False(keeper.HasClass(ctx, unsavedClassA.BrandId, unsavedClassA.Id))

	//get all
	suite.Require().Equal(len(keeper.GetClasses(ctx)), len(savedClasses))

	//get all of brand
	suite.Require().Equal(len(keeper.GetClassesOfBrand(ctx, brandId)), 2)
	suite.Require().Equal(len(keeper.GetClassesOfBrand(ctx, brandId2)), 1)
	suite.Require().Equal(len(keeper.GetClassesOfBrand(ctx, "random")), 0)
}
