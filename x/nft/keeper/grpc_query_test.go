package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

// Classes queries all Classes
func (suite *KeeperTestSuite) TestClasses() {
	keeper, ctx, queryClient := suite.app.NFTKeeper, suite.ctx, suite.queryClient
	wrapCtx := sdk.WrapSDKContext(ctx)

	//invalid brandID
	req := &types.QueryClassesRequest{BrandId: "-"}
	res, err := queryClient.Classes(wrapCtx, req)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	req = &types.QueryClassesRequest{}
	res, err = queryClient.Classes(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Nil(res.Pagination.NextKey)
	suite.Require().Len(res.Classes, 0)

	brandID := "brandid"
	brandID2 := "brandid2"
	brandID3 := "brandid3"
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID, "classid", 10_000, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID, "classid2", 1000, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID2, "classid", 100, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID2, "classid2", 100, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID, "classid3", 10, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID, "classid4", 1, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID3, "classid", 9999, types.NewClassDescription("", "", "", ""))))

	pageReq := &query.PageRequest{
		Key:        nil,
		Limit:      3,
		CountTotal: false,
	}

	req = &types.QueryClassesRequest{Pagination: pageReq}
	res, err = queryClient.Classes(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotNil(res.Pagination.NextKey)
	suite.Require().Len(res.Classes, 3)

	pageReq = &query.PageRequest{
		Key:        res.Pagination.NextKey,
		CountTotal: true,
	}

	req = &types.QueryClassesRequest{Pagination: pageReq}
	res, err = queryClient.Classes(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Nil(res.Pagination.NextKey)
	suite.Require().Len(res.Classes, 4)

	req = &types.QueryClassesRequest{BrandId: brandID}
	res, err = queryClient.Classes(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Nil(res.Pagination.NextKey)
	suite.Require().Len(res.Classes, 4)

	req = &types.QueryClassesRequest{BrandId: brandID2}
	res, err = queryClient.Classes(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Nil(res.Pagination.NextKey)
	suite.Require().Len(res.Classes, 2)

	req = &types.QueryClassesRequest{BrandId: brandID3}
	res, err = queryClient.Classes(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Nil(res.Pagination.NextKey)
	suite.Require().Len(res.Classes, 1)
}

func (suite *KeeperTestSuite) TestClass() {
	keeper, ctx, queryClient := suite.app.NFTKeeper, suite.ctx, suite.queryClient
	wrapCtx := sdk.WrapSDKContext(ctx)

	//invalid brandID
	req := &types.QueryClassRequest{BrandId: "", ClassId: "classid"}
	res, err := queryClient.Class(wrapCtx, req)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	//invalid classID
	req = &types.QueryClassRequest{BrandId: "brandid", ClassId: ""}
	res, err = queryClient.Class(wrapCtx, req)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	brandID := "brandid"
	brandID2 := "brandid2"
	brandID3 := "brandid3"

	classID := "classid"
	classID2 := "classid2"
	classID3 := "classid3"
	classID4 := "classid4"

	// not exist class
	req = &types.QueryClassRequest{BrandId: brandID, ClassId: classID}
	res, err = queryClient.Class(wrapCtx, req)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID, classID, 10_000, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID, classID2, 1000, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID2, classID, 100, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID2, classID2, 100, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID, classID3, 10, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID, classID4, 1, types.NewClassDescription("", "", "", ""))))
	suite.Require().NoError(keeper.SaveClass(ctx, types.NewClass(brandID3, classID, 9999, types.NewClassDescription("", "", "", ""))))

	req = &types.QueryClassRequest{BrandId: brandID, ClassId: classID}
	res, err = queryClient.Class(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(types.GetClassUniqueID(brandID, classID), types.GetClassUniqueID(res.Class.BrandId, res.Class.Id))

	req = &types.QueryClassRequest{BrandId: brandID3, ClassId: classID}
	res, err = queryClient.Class(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(types.GetClassUniqueID(brandID3, classID), types.GetClassUniqueID(res.Class.BrandId, res.Class.Id))

	//empty classID
	req = &types.QueryClassRequest{BrandId: brandID3, ClassId: classID2}
	res, err = queryClient.Class(wrapCtx, req)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	//empty brandID
	req = &types.QueryClassRequest{BrandId: "random", ClassId: classID}
	res, err = queryClient.Class(wrapCtx, req)
	suite.Require().Error(err)
	suite.Require().Nil(res)
}
