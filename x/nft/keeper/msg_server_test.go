package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	brandtypes "github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

// Classes queries all Classes
func (suite *KeeperTestSuite) CreateClass() {
	app, ctx, msgServer := suite.app, suite.ctx, suite.msgServer
	wrapCtx := sdk.WrapSDKContext(ctx)

	addr := sdk.AccAddress("addr1...")
	addr2 := sdk.AccAddress("addr2...")
	addr3 := sdk.AccAddress("addr3...")

	brandID := "brandid"
	brandID2 := "brandid2"

	desc := types.NewClassDescription("", "", "", "")

	msgA := types.NewMsgCreateClass(brandID, "classid", addr.String(), 10_000, desc)
	msgB := types.NewMsgCreateClass(brandID2, "classid", addr2.String(), 10_000, desc)

	res, err := msgServer.CreateClass(wrapCtx, msgA)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	res, err = msgServer.CreateClass(wrapCtx, msgB)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	app.BrandKeeper.SetBrand(ctx, brandtypes.NewBrand(brandID, addr, brandtypes.NewBrandDescription("name", "", "")))
	app.BrandKeeper.SetBrand(ctx, brandtypes.NewBrand(brandID2, addr2, brandtypes.NewBrandDescription("name", "", "")))

	res, err = msgServer.CreateClass(wrapCtx, msgA)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	res, err = msgServer.CreateClass(wrapCtx, msgB)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	//invalid creator and owner
	res, err = msgServer.CreateClass(
		wrapCtx,
		types.NewMsgCreateClass(msgB.BrandId, msgB.Id+"2", addr3.String(), 10_000, desc),
	)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	//duplicate uniqueID
	res, err = msgServer.CreateClass(wrapCtx, msgA)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	res, err = msgServer.CreateClass(wrapCtx, msgB)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	res, err = msgServer.CreateClass(
		wrapCtx,
		types.NewMsgCreateClass(msgB.BrandId, msgB.Id+"2", "ddawd", 10_000, desc),
	)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	res, err = msgServer.CreateClass(
		wrapCtx,
		types.NewMsgCreateClass(msgB.BrandId, msgB.Id+"2", addr2.String(), 10_001, desc),
	)
	suite.Require().Error(err)
	suite.Require().Nil(res)
}

func (suite *KeeperTestSuite) EditClass() {
	app, ctx, msgServer := suite.app, suite.ctx, suite.msgServer
	wrapCtx := sdk.WrapSDKContext(ctx)

	addr := sdk.AccAddress("addr1...")
	addr2 := sdk.AccAddress("addr2...")

	brandID := "brandid"
	brandID2 := "brandid2"

	desc := types.NewClassDescription("", "", "", "")

	msgA := types.NewMsgEditClass(brandID, "classid", addr.String(), 10_000, desc)
	msgB := types.NewMsgEditClass(brandID2, "classid", addr2.String(), 10_000, desc)

	res, err := msgServer.EditClass(wrapCtx, msgA)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	res, err = msgServer.EditClass(wrapCtx, msgB)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	//without checking brand
	app.NFTKeeper.SaveClass(ctx, types.NewClass(msgA.BrandId, msgA.Id, msgA.FeeBasisPoints, msgA.Description))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(msgB.BrandId, msgB.Id, msgB.FeeBasisPoints, msgB.Description))

	res, err = msgServer.EditClass(wrapCtx, msgA)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	res, err = msgServer.EditClass(wrapCtx, msgB)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	msgC := msgA
	msgC.Editor = sdk.AccAddress("randome...").String()
	res, err = msgServer.EditClass(wrapCtx, msgA)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	msgD := msgB
	msgD.Editor = msgA.Editor
	res, err = msgServer.EditClass(wrapCtx, msgB)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	res, err = msgServer.EditClass(
		wrapCtx,
		types.NewMsgEditClass(msgB.BrandId, msgB.Id+"2", "ddawd", 10_000, desc),
	)
	suite.Require().Error(err)
	suite.Require().Nil(res)

	res, err = msgServer.EditClass(
		wrapCtx,
		types.NewMsgEditClass(msgB.BrandId, msgB.Id+"2", msgB.Editor, 10_001, desc),
	)
	suite.Require().Error(err)
	suite.Require().Nil(res)
}
