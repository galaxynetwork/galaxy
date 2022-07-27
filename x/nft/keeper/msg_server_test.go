package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	brandtypes "github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

func (suite *KeeperTestSuite) TestCreateClass() {

	app, ctx, msgServer := suite.app, suite.ctx, suite.msgServer
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB := "brandIDA", "brandIDB"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb")

	desc := types.NewClassDescription("name", "", "", "")
	//invalid arguments
	//details message in sdk.BasicValidate
	tests := []struct {
		msg *types.MsgCreateClass
	}{
		{types.NewMsgCreateClass("", classIDA, ownerA.String(), 10_000, desc)},
		{types.NewMsgCreateClass(brandIDA, "", ownerA.String(), 10_000, desc)},
		{types.NewMsgCreateClass(brandIDA, classIDA, ownerA.String(), 10_001, desc)},
		{types.NewMsgCreateClass(brandIDA, classIDA, ownerA.String(), 0, types.NewClassDescription("", "", "", ""))},
		{types.NewMsgCreateClass(brandIDA, classIDA, "", 10_000, desc)},
	}

	for _, test := range tests {
		res, err := msgServer.CreateClass(wrapCtx, test.msg)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}

	//owner should be same for brand
	classes := []struct {
		b string
		c string
		o sdk.AccAddress
	}{
		{brandIDA, classIDA, ownerA}, {brandIDA, classIDB, ownerA}, {brandIDA, classIDC, ownerA},
		{brandIDB, classIDA, ownerB}, {brandIDB, classIDB, ownerB},
	}

	for _, d := range classes {
		hasBrand := app.BrandKeeper.HasBrand(ctx, d.b)

		msg := types.NewMsgCreateClass(d.b, d.c, d.o.String(), 10_000, desc)

		if !hasBrand {
			_, err := msgServer.CreateClass(wrapCtx, msg)
			suite.Require().Equal(err, brandtypes.ErrNotFoundBrand)

			suite.Require().NoError(
				app.BrandKeeper.SetBrand(ctx, brandtypes.NewBrand(d.b, d.o, brandtypes.NewBrandDescription("name", "", ""))),
			)
		}

		msg.Creator = sdk.AccAddress("randomaddress").String()
		_, err := msgServer.CreateClass(wrapCtx, msg)
		suite.Require().Equal(err, types.ErrUnauthorized)

		msg.Creator = d.o.String()
		_, err = msgServer.CreateClass(wrapCtx, msg)
		suite.Require().NoError(err)

		_, err = msgServer.CreateClass(wrapCtx, msg)
		suite.Require().Error(err)
	}
}

func (suite *KeeperTestSuite) TestEditClass() {
	app, ctx, msgServer := suite.app, suite.ctx, suite.msgServer
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB := "brandIDA", "brandIDB"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb")

	desc := types.NewClassDescription("name", "", "", "")
	//invalid arguments
	//details message in sdk.BasicValidate
	tests := []struct {
		msg *types.MsgEditClass
	}{
		{types.NewMsgEditClass("", classIDA, ownerA.String(), 10_000, desc)},
		{types.NewMsgEditClass(brandIDA, "", ownerA.String(), 10_000, desc)},
		{types.NewMsgEditClass(brandIDA, classIDA, ownerA.String(), 10_001, desc)},
		{types.NewMsgEditClass(brandIDA, classIDA, ownerA.String(), 0, types.NewClassDescription("", "", "", ""))},
		{types.NewMsgEditClass(brandIDA, classIDA, "", 10_000, desc)},
	}

	for _, test := range tests {
		res, err := msgServer.EditClass(wrapCtx, test.msg)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}

	//owner should be same for brand
	classes := []struct {
		b string
		c string
		o sdk.AccAddress
	}{
		{brandIDA, classIDA, ownerA}, {brandIDA, classIDB, ownerA}, {brandIDA, classIDC, ownerA},
		{brandIDB, classIDA, ownerB}, {brandIDB, classIDB, ownerB},
	}

	for _, d := range classes {
		hasBrand := app.BrandKeeper.HasBrand(ctx, d.b)

		msg := types.NewMsgEditClass(d.b, d.c, d.o.String(), 10_000, desc)

		if !hasBrand {
			_, err := msgServer.EditClass(wrapCtx, msg)
			suite.Require().Equal(err, brandtypes.ErrNotFoundBrand)

			suite.Require().NoError(
				app.BrandKeeper.SetBrand(ctx, brandtypes.NewBrand(d.b, d.o, brandtypes.NewBrandDescription("name", "", ""))),
			)
		}

		_, err := msgServer.EditClass(wrapCtx, msg)
		suite.Require().Equal(err, types.ErrNotFoundClass)

		savedClass := types.NewClass(msg.BrandId, msg.Id, msg.FeeBasisPoints, msg.Description)
		suite.Require().NoError(app.NFTKeeper.SaveClass(ctx, savedClass))

		class, exist := app.NFTKeeper.GetClass(ctx, msg.BrandId, msg.Id)
		suite.Require().True(exist)
		suite.Require().Equal(class, savedClass)

		savedClass = types.NewClass(msg.BrandId, msg.Id, 123, types.NewClassDescription("newname", "details", "https://galaxy", "ipfs://image"))
		msg.Description = savedClass.Description
		msg.FeeBasisPoints = savedClass.FeeBasisPoints
		_, err = msgServer.EditClass(wrapCtx, msg)
		suite.Require().NoError(err)

		suite.Require().NotEqual(class, savedClass)
		class, exist = app.NFTKeeper.GetClass(ctx, msg.BrandId, msg.Id)
		suite.Require().True(exist)
		suite.Require().Equal(class, savedClass)

		msg.Editor = sdk.AccAddress("randomaddress").String()
		_, err = msgServer.EditClass(wrapCtx, msg)
		suite.Require().Equal(err, types.ErrUnauthorized)
	}
}

func (suite *KeeperTestSuite) TestMintNFT() {
	app, ctx, msgServer := suite.app, suite.ctx, suite.msgServer
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB := "brandIDA", "brandIDB"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB, ownerC, recipient := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb"), sdk.AccAddress("ownerc"), sdk.AccAddress("recipient")

	//invalid arguments
	tests := []struct {
		msg *types.MsgMintNFT
	}{
		{types.NewMsgMintNFT("", classIDA, "ipfs://nft", "", ownerA.String(), recipient.String())},
		{types.NewMsgMintNFT(brandIDA, "", "ipfs://nft", "", ownerA.String(), recipient.String())},
		{types.NewMsgMintNFT(brandIDA, classIDA, "", "", ownerA.String(), recipient.String())},
		{types.NewMsgMintNFT(brandIDA, classIDA, "ipfs://nft", "", "", recipient.String())},
		{types.NewMsgMintNFT(brandIDA, classIDA, "ipfs://nft", "", ownerA.String(), "")},
	}

	for _, test := range tests {
		res, err := msgServer.MintNFT(wrapCtx, test.msg)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}

	nftData := []struct {
		num int
		b   string
		c   string
		o   sdk.AccAddress
	}{
		{3, brandIDA, classIDA, ownerA}, {10, brandIDA, classIDA, ownerB}, {2, brandIDA, classIDB, ownerA}, {2, brandIDA, classIDC, ownerA},
		{5, brandIDB, classIDA, ownerB}, {3, brandIDB, classIDB, ownerB}, {10, brandIDB, classIDB, ownerC},
	}

	classMap := map[string]int{}

	for _, d := range nftData {
		cuID := types.GetClassUniqueID(d.b, d.c)
		if _, ok := classMap[cuID]; !ok {
			classMap[cuID] = 0
		}

		for i := 0; i < d.num; i++ {
			classMap[cuID]++
			hasBrand := app.BrandKeeper.HasBrand(ctx, d.b)

			msg := types.NewMsgMintNFT(d.b, d.c, "ipfs://nft", "", recipient.String(), recipient.String())

			if !hasBrand {
				_, err := msgServer.MintNFT(wrapCtx, msg)
				suite.Require().Equal(err, brandtypes.ErrNotFoundBrand)

				suite.Require().NoError(
					app.BrandKeeper.SetBrand(ctx, brandtypes.NewBrand(msg.BrandId, recipient, brandtypes.NewBrandDescription("name", "", ""))),
				)
			}

			hasClass := app.NFTKeeper.HasClass(ctx, msg.BrandId, msg.ClassId)
			if !hasClass {
				_, err := msgServer.MintNFT(wrapCtx, msg)
				suite.Require().Equal(err, types.ErrNotFoundClass)

				suite.Require().NoError(
					app.NFTKeeper.SaveClass(ctx, types.NewClass(msg.BrandId, msg.ClassId, 0, types.NewClassDescription("", "", "", ""))),
				)
			}

			nmsg := *msg
			nmsg.Minter = sdk.AccAddress("randomaddress").String()
			_, err := msgServer.MintNFT(wrapCtx, &nmsg)
			suite.Require().Equal(err, types.ErrUnauthorized)

			res, err := msgServer.MintNFT(wrapCtx, msg)
			suite.Require().NoError(err)
			suite.Require().NotZero(res.Id)
		}
	}

	for uid, len := range classMap {
		brandID, classID := types.ParseClassUniqueID(uid)
		nfts := app.NFTKeeper.GetNFTsOfClass(ctx, brandID, classID)
		suite.Require().Len(nfts, len)
	}

}

func (suite *KeeperTestSuite) TestUpdateNFT() {
	app, ctx, msgServer := suite.app, suite.ctx, suite.msgServer
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB := "brandIDA", "brandIDB"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB, ownerC := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb"), sdk.AccAddress("ownerc")

	//invalid arguments
	tests := []struct {
		msg *types.MsgUpdateNFT
	}{
		{types.NewMsgUpdateNFT("", classIDA, 1, "", ownerA.String())},
		{types.NewMsgUpdateNFT(brandIDA, "", 1, "", ownerA.String())},
		{types.NewMsgUpdateNFT(brandIDA, classIDA, 0, "", ownerA.String())},
		{types.NewMsgUpdateNFT(brandIDA, classIDA, 1, "", "")},
	}

	for _, test := range tests {
		res, err := msgServer.UpdateNFT(wrapCtx, test.msg)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}

	//ignore set brand first
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDB, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDC, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDB, 0, types.NewClassDescription("", "", "", "")))

	nftData := []struct {
		num int
		b   string
		c   string
		o   sdk.AccAddress
	}{
		{3, brandIDA, classIDA, ownerA}, {10, brandIDA, classIDA, ownerB}, {2, brandIDA, classIDB, ownerA}, {2, brandIDA, classIDC, ownerA},
		{5, brandIDB, classIDA, ownerB}, {3, brandIDB, classIDB, ownerB}, {10, brandIDB, classIDB, ownerC},
	}

	//mint nfts
	for _, d := range nftData {
		for i := 0; i < d.num; i++ {
			nft, err := app.NFTKeeper.GenNFT(ctx, d.b, d.c, "ipfs://nft", "")
			suite.Require().NoError(err)

			msg := types.NewMsgUpdateNFT(nft.BrandId, nft.ClassId, nft.Id, "", d.o.String())
			_, err = msgServer.UpdateNFT(wrapCtx, msg)
			suite.Require().Equal(err, types.ErrNotFoundNFT)

			suite.Require().NoError(app.NFTKeeper.MintNFT(ctx, nft, d.o))

			msg.VarUri = "ipfs://customer_uri"
			res, err := msgServer.UpdateNFT(wrapCtx, msg)
			suite.Require().NoError(err)
			suite.Require().NotNil(res)

			msg.Sender = sdk.AccAddress("randomowner").String()
			_, err = msgServer.UpdateNFT(wrapCtx, msg)
			suite.Require().Equal(err, types.ErrUnauthorized)
		}
	}
}

func (suite *KeeperTestSuite) TestBurnNFT() {
	app, ctx, msgServer := suite.app, suite.ctx, suite.msgServer
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB := "brandIDA", "brandIDB"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB, ownerC := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb"), sdk.AccAddress("ownerc")

	//invalid arguments
	tests := []struct {
		msg *types.MsgBurnNFT
	}{
		{types.NewMsgBurnNFT("", classIDA, 1, ownerA.String())},
		{types.NewMsgBurnNFT(brandIDA, "", 1, ownerA.String())},
		{types.NewMsgBurnNFT(brandIDA, classIDA, 0, ownerA.String())},
		{types.NewMsgBurnNFT(brandIDA, classIDA, 1, "")},
	}

	for _, test := range tests {
		res, err := msgServer.BurnNFT(wrapCtx, test.msg)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}

	//ignore set brand first
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDB, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDC, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDB, 0, types.NewClassDescription("", "", "", "")))

	nftData := []struct {
		num int
		b   string
		c   string
		o   sdk.AccAddress
	}{
		{3, brandIDA, classIDA, ownerA}, {10, brandIDA, classIDA, ownerB}, {2, brandIDA, classIDB, ownerA}, {2, brandIDA, classIDC, ownerA},
		{5, brandIDB, classIDA, ownerB}, {3, brandIDB, classIDB, ownerB}, {10, brandIDB, classIDB, ownerC},
	}

	for _, d := range nftData {
		for i := 0; i < d.num; i++ {
			nft, err := app.NFTKeeper.GenNFT(ctx, d.b, d.c, "ipfs://nft", "")
			suite.Require().NoError(err)

			msg := types.NewMsgBurnNFT(nft.BrandId, nft.ClassId, nft.Id, d.o.String())
			_, err = msgServer.BurnNFT(wrapCtx, msg)
			suite.Require().Equal(err, types.ErrNotFoundNFT)

			suite.Require().NoError(app.NFTKeeper.MintNFT(ctx, nft, d.o))

			msg.Sender = sdk.AccAddress("randomowner").String()
			_, err = msgServer.BurnNFT(wrapCtx, msg)
			suite.Require().Equal(err, types.ErrUnauthorized)

			msg.Sender = d.o.String()
			res, err := msgServer.BurnNFT(wrapCtx, msg)
			suite.Require().NoError(err)
			suite.Require().NotNil(res)

			suite.Require().True(
				app.NFTKeeper.BurnedNFT(ctx, nft.BrandId, nft.ClassId, nft.Id),
			)
		}
	}
}

func (suite *KeeperTestSuite) TestTransferNFT() {
	app, ctx, msgServer := suite.app, suite.ctx, suite.msgServer
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB := "brandIDA", "brandIDB"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB, ownerC, recipient := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb"), sdk.AccAddress("ownerc"), sdk.AccAddress("recipient")

	//invalid arguments
	tests := []struct {
		msg *types.MsgTransferNFT
	}{
		{types.NewMsgTransferNFT("", classIDA, 1, ownerA.String(), recipient.String())},
		{types.NewMsgTransferNFT(brandIDA, "", 1, ownerA.String(), recipient.String())},
		{types.NewMsgTransferNFT(brandIDA, classIDA, 0, ownerA.String(), recipient.String())},
		{types.NewMsgTransferNFT(brandIDA, classIDA, 1, "", recipient.String())},
		{types.NewMsgTransferNFT(brandIDA, classIDA, 1, ownerA.String(), "")},
	}

	for _, test := range tests {
		res, err := msgServer.TransferNFT(wrapCtx, test.msg)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}

	//ignore set brand first
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDB, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDC, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDB, 0, types.NewClassDescription("", "", "", "")))

	nftData := []struct {
		num int
		b   string
		c   string
		o   sdk.AccAddress
	}{
		{3, brandIDA, classIDA, ownerA}, {10, brandIDA, classIDA, ownerB}, {2, brandIDA, classIDB, ownerA}, {2, brandIDA, classIDC, ownerA},
		{5, brandIDB, classIDA, ownerB}, {3, brandIDB, classIDB, ownerB}, {10, brandIDB, classIDB, ownerC},
	}

	for _, d := range nftData {
		for i := 0; i < d.num; i++ {
			nft, err := app.NFTKeeper.GenNFT(ctx, d.b, d.c, "ipfs://nft", "")
			suite.Require().NoError(err)

			msg := types.NewMsgTransferNFT(nft.BrandId, nft.ClassId, nft.Id, d.o.String(), recipient.String())
			_, err = msgServer.TransferNFT(wrapCtx, msg)
			suite.Require().Equal(err, types.ErrNotFoundNFT)

			suite.Require().NoError(app.NFTKeeper.MintNFT(ctx, nft, d.o))

			msg.Sender = sdk.AccAddress("randomowner").String()
			_, err = msgServer.TransferNFT(wrapCtx, msg)
			suite.Require().Equal(err, types.ErrUnauthorized)

			suite.Require().NotEqual(d.o, recipient)
			suite.Require().Equal(d.o, app.NFTKeeper.GetOwner(ctx, nft.BrandId, nft.ClassId, nft.Id))

			msg.Sender = d.o.String()
			res, err := msgServer.TransferNFT(wrapCtx, msg)
			suite.Require().NoError(err)
			suite.Require().NotNil(res)

			suite.Require().Equal(recipient, app.NFTKeeper.GetOwner(ctx, nft.BrandId, nft.ClassId, nft.Id))

			//check transfer same receipient is passed
			msg.Sender = msg.Recipient
			_, err = msgServer.TransferNFT(wrapCtx, msg)
			suite.Require().NoError(err)
		}
	}
}
