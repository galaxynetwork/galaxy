package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
	minttypes "github.com/galaxies-labs/galaxy/x/mint/types"
)

func (suite *KeeperTestSuite) TestCreateBrand() {
	app, ctx, msgServer, queryClient := suite.app, suite.ctx, suite.msgServer, suite.queryClient
	wrapCtx := sdk.WrapSDKContext(ctx)

	bak1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())
	bak2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())

	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, bak1))
	amount := sdk.NewInt(1_000_000_000_000)
	initialMintedCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultBrandCreationFeeDenom, amount))
	suite.Require().NoError(
		app.MintKeeper.MintCoins(ctx, initialMintedCoins),
	)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, bak1, initialMintedCoins)

	tests := []struct {
		msg        *types.MsgCreateBrand
		expectPass bool
	}{
		{types.NewMsgCreateBrand("brandId", bak1.String(), types.NewBrandDescription("name", "", "")), true},
		//empty creation fee
		{types.NewMsgCreateBrand("brandId22", bak2.String(), types.NewBrandDescription("name", "", "")), false},
		//same brandi d
		{types.NewMsgCreateBrand("brandId", bak1.String(), types.NewBrandDescription("name", "", "")), false},
		//invalid brand description
		{types.NewMsgCreateBrand("brandIdb", bak1.String(), types.NewBrandDescription("", "", "")), false},
		//empty brand id
		{types.NewMsgCreateBrand("", bak1.String(), types.NewBrandDescription("name", "", "")), false},
		//invalid brand id
		{types.NewMsgCreateBrand("_brandId", bak1.String(), types.NewBrandDescription("name", "", "")), false},
		//invalid owner address
		{types.NewMsgCreateBrand("brandIdbc", "bech32", types.NewBrandDescription("name", "", "")), false},
	}

	balance := app.BankKeeper.GetBalance(ctx, bak1, initialMintedCoins.GetDenomByIndex(0))
	suite.Require().Equal(balance, initialMintedCoins[0])

	for _, test := range tests {
		res, err := msgServer.CreateBrand(wrapCtx, test.msg)
		if !test.expectPass {
			suite.Require().Error(err)

		} else {
			suite.Require().NoError(err)
			suite.Require().NotNil(res.BrandAddress)

			acc, err := sdk.AccAddressFromBech32(res.BrandAddress)
			suite.Require().NoError(err)
			suite.Require().True(app.AccountKeeper.HasAccount(ctx, acc))

			res, err := queryClient.BrandsByOwner(wrapCtx, &types.QueryBrandsByOwnerRequest{Owner: test.msg.Owner})
			suite.Require().NoError(err)
			suite.Require().Len(res.Brands, 1)
			suite.Require().Equal(res.Brands[0].Id, test.msg.Id)

			newBalance := app.BankKeeper.GetBalance(ctx, bak1, initialMintedCoins.GetDenomByIndex(0))
			suite.Require().Equal(
				balance.Sub(app.BrandKeeper.GetParams(ctx).BrandCreationFee),
				newBalance,
			)

			balance = newBalance
		}
	}
}

func (suite *KeeperTestSuite) TestEditBrand() {
	app, ctx, msgServer, queryClient := suite.app, suite.ctx, suite.msgServer, suite.queryClient
	wrapCtx := sdk.WrapSDKContext(ctx)

	bak1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())
	bak2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())

	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, bak1))
	initialMintedCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultBrandCreationFeeDenom, sdk.NewInt(1_000_000_000_000)))
	suite.Require().NoError(
		app.MintKeeper.MintCoins(ctx, initialMintedCoins),
	)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, bak1, initialMintedCoins)

	msgBrandA := types.NewMsgCreateBrand("brandId", bak1.String(), types.NewBrandDescription("name", "", ""))
	msgBrandB := types.NewMsgCreateBrand("brandId22", bak2.String(), types.NewBrandDescription("name", "", ""))

	_, err := msgServer.CreateBrand(wrapCtx, msgBrandA)
	suite.Require().NoError(err)

	_, err = msgServer.CreateBrand(wrapCtx, msgBrandB)
	suite.Require().Error(err)

	tests := []struct {
		originMsg  *types.MsgCreateBrand
		msg        *types.MsgEditBrand
		expectPass bool
	}{
		{msgBrandB, types.NewMsgEditBrand("brandId22", bak1.String(), types.NewBrandDescription("changed", "", "")), false},
		{msgBrandA, types.NewMsgEditBrand("brandId", bak2.String(), types.NewBrandDescription("changed", "", "")), false},
		{msgBrandA, types.NewMsgEditBrand("brandId", bak1.String(), types.NewBrandDescription("", "", "")), false},
		{msgBrandA, types.NewMsgEditBrand("brandId", bak1.String(), types.NewBrandDescription("changed", "", "")), true},
	}

	for _, test := range tests {
		_, err := msgServer.EditBrand(wrapCtx, test.msg)
		if test.expectPass {
			suite.Require().NoError(err)
			res, err := queryClient.Brand(wrapCtx, &types.QueryBrandRequest{BrandId: test.msg.Id})
			suite.Require().NoError(err)

			suite.Require().Equal(res.Brand.Id, test.originMsg.Id)
			suite.Require().Equal(res.Brand.Owner, test.originMsg.Owner)
			suite.Require().Equal(res.Brand.Description, test.msg.Description)
			suite.Require().NotEqual(res.Brand.Description, test.originMsg.Description)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestTransferOwnershipBrand() {
	app, ctx, msgServer, queryClient := suite.app, suite.ctx, suite.msgServer, suite.queryClient
	wrapCtx := sdk.WrapSDKContext(ctx)

	bak1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())
	bak2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())
	bak3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().String())

	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, bak1))
	initialMintedCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultBrandCreationFeeDenom, sdk.NewInt(1_000_000_000_000)))
	suite.Require().NoError(
		app.MintKeeper.MintCoins(ctx, initialMintedCoins),
	)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, bak1, initialMintedCoins)

	msgBrandA := types.NewMsgCreateBrand("brandId", bak1.String(), types.NewBrandDescription("name", "", ""))
	msgBrandB := types.NewMsgCreateBrand("brandId22", bak2.String(), types.NewBrandDescription("name", "", ""))

	_, err := msgServer.CreateBrand(wrapCtx, msgBrandA)
	suite.Require().NoError(err)

	_, err = msgServer.CreateBrand(wrapCtx, msgBrandB)
	suite.Require().Error(err)

	tests := []struct {
		msg        *types.MsgTransferOwnershipBrand
		expectPass bool
	}{
		{types.NewMsgTransferOwnershipBrand(msgBrandA.Id, bak2.String(), bak3.String()), false},
		{types.NewMsgTransferOwnershipBrand(msgBrandB.Id, msgBrandA.Owner, bak3.String()), false},
		{types.NewMsgTransferOwnershipBrand(msgBrandA.Id, msgBrandA.Owner, bak3.String()), true},
	}

	for _, test := range tests {
		_, err := msgServer.TransferOwnershipBrand(wrapCtx, test.msg)
		if test.expectPass {
			suite.Require().NoError(err)

			res, err := queryClient.Brand(wrapCtx, &types.QueryBrandRequest{BrandId: test.msg.Id})
			suite.Require().NoError(err)

			suite.Require().NotEqual(res.Brand.Owner, test.msg.Owner)
			suite.Require().Equal(res.Brand.Owner, test.msg.DestOwner)

			res2, err2 := queryClient.BrandsByOwner(wrapCtx, &types.QueryBrandsByOwnerRequest{Owner: test.msg.Owner})
			suite.Require().NoError(err2)
			suite.Require().NotNil(res2)
			suite.Require().Len(res2.Brands, 0)

			res2, err2 = queryClient.BrandsByOwner(wrapCtx, &types.QueryBrandsByOwnerRequest{Owner: test.msg.DestOwner})
			suite.Require().NoError(err2)
			suite.Require().NotNil(res2)
			suite.Require().Len(res2.Brands, 1)

		} else {
			suite.Require().Error(err)
		}
	}
}
