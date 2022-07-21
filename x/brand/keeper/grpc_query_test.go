package keeper_test

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

func (suite *KeeperTestSuite) TestBrands() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	res, err := queryClient.Brands(sdk.WrapSDKContext(ctx), &types.QueryBrandsRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Nil(res.Pagination.NextKey)
	suite.Require().Len(res.Brands, 0)

	var brands types.Brands
	var maxLen = 11

	for i := 0; i < maxLen; i++ {
		_, _, addr := testdata.KeyTestPubAddr()
		brand := types.NewBrand(strings.Join([]string{"brand", strconv.Itoa(i)}, ""), addr, types.NewBrandDescription("name", "", ""))
		suite.Require().NoError(brand.Validate())
		suite.Require().NoError(app.BrandKeeper.SetBrand(ctx, brand))
		brands = append(brands, brand)
	}

	pageReq := &query.PageRequest{
		Key:        nil,
		Limit:      10,
		CountTotal: false,
	}

	res, err = queryClient.Brands(sdk.WrapSDKContext(ctx), &types.QueryBrandsRequest{Pagination: pageReq})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Brands, 10)
	suite.Require().NotNil(res.Pagination.NextKey)

	pageReq = &query.PageRequest{
		Key:        res.Pagination.NextKey,
		Limit:      10,
		CountTotal: true,
	}
	res, err = queryClient.Brands(sdk.WrapSDKContext(ctx), &types.QueryBrandsRequest{Pagination: pageReq})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Brands, maxLen-10)
	suite.Require().Nil(res.Pagination.NextKey)
}

func (suite *KeeperTestSuite) TestBrand() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	var invalidBrandID = "invalidbrandid"

	_, _, addr := testdata.KeyTestPubAddr()
	brand := types.NewBrand("brandid", addr, types.NewBrandDescription("name", "", ""))
	suite.Require().NoError(brand.Validate())

	res, err := queryClient.Brand(sdk.WrapSDKContext(ctx), &types.QueryBrandRequest{})
	suite.Require().Error(err)

	res, err = queryClient.Brand(sdk.WrapSDKContext(ctx), &types.QueryBrandRequest{BrandId: brand.Id})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Empty(res.Brand)

	suite.Require().NoError(app.BrandKeeper.SetBrand(ctx, brand))

	res, err = queryClient.Brand(sdk.WrapSDKContext(ctx), &types.QueryBrandRequest{BrandId: brand.Id})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotNil(res.Brand)
	suite.Require().Equal(res.Brand.Id, brand.Id)

	res, err = queryClient.Brand(sdk.WrapSDKContext(ctx), &types.QueryBrandRequest{BrandId: invalidBrandID})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Empty(res.Brand)
}
