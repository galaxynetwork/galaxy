package keeper_test

import (
	"fmt"
	"strings"

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

	brandID, brandID2, brandID3 := "brandid", "brandid2", "brandid3"
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

	brandID, brandID2, brandID3 := "brandid", "brandid2", "brandid3"
	classID, classID2, classID3, classID4 := "classid", "classid2", "classid3", "classid4"

	//invalid arguments
	tests := []struct {
		req *types.QueryClassRequest
	}{
		{&types.QueryClassRequest{BrandId: "", ClassId: classID}},
		{&types.QueryClassRequest{BrandId: brandID, ClassId: ""}},
		{&types.QueryClassRequest{BrandId: brandID, ClassId: ".classid"}},
		{&types.QueryClassRequest{BrandId: ".brandID", ClassId: classID}},
		//not found
		{&types.QueryClassRequest{BrandId: brandID, ClassId: classID}},
	}

	for _, test := range tests {
		res, err := queryClient.Class(wrapCtx, test.req)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}

	classes := []struct {
		class types.Class
	}{
		{types.NewClass(brandID, classID, 10_000, types.NewClassDescription("", "", "", ""))},
		{types.NewClass(brandID, classID2, 1000, types.NewClassDescription("", "", "", ""))},
		{types.NewClass(brandID2, classID, 100, types.NewClassDescription("", "", "", ""))},
		{types.NewClass(brandID2, classID2, 100, types.NewClassDescription("", "", "", ""))},
		{types.NewClass(brandID, classID3, 10, types.NewClassDescription("", "", "", ""))},
		{types.NewClass(brandID, classID4, 1, types.NewClassDescription("", "", "", ""))},
		{types.NewClass(brandID3, classID, 9999, types.NewClassDescription("", "", "", ""))},
	}

	for _, class := range classes {
		suite.Require().NoError(keeper.SaveClass(ctx, class.class))

		req := &types.QueryClassRequest{BrandId: class.class.BrandId, ClassId: class.class.Id}
		res, err := queryClient.Class(wrapCtx, req)
		suite.Require().NoError(err)
		suite.Require().NotNil(res)
		suite.Require().Equal(res.Class, class.class)

		//check for not found after data stored
		req = &types.QueryClassRequest{BrandId: class.class.BrandId, ClassId: "randomclassID"}
		res, err = queryClient.Class(wrapCtx, req)
		fmt.Println((err))
		suite.Require().Error(err)
		suite.Require().Nil(res)

		req = &types.QueryClassRequest{BrandId: "randombrandID", ClassId: class.class.Id}
		res, err = queryClient.Class(wrapCtx, req)
		fmt.Println((err))
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}
}

func (suite *KeeperTestSuite) TestNFTs() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB, brandIDC := "brandIDA", "brandIDB", "brandIDC"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB, ownerC := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb"), sdk.AccAddress("ownerc")

	pageReq := &query.PageRequest{}

	//invalid arguments
	tests := []struct {
		expectPass bool
		req        *types.QueryNFTsRequest
	}{
		{true, &types.QueryNFTsRequest{BrandId: "", ClassId: "", Owner: "", Pagination: pageReq}},
		{true, &types.QueryNFTsRequest{BrandId: "", ClassId: "", Owner: ""}},

		{false, &types.QueryNFTsRequest{BrandId: "invalid.brandID", ClassId: classIDA, Owner: ownerA.String(), Pagination: pageReq}},
		{false, &types.QueryNFTsRequest{BrandId: "invalid.brandID", ClassId: classIDA, Owner: ownerA.String(), Pagination: pageReq}},
		{false, &types.QueryNFTsRequest{BrandId: brandIDA, ClassId: "invalud.classID", Owner: ownerA.String(), Pagination: pageReq}},
		{false, &types.QueryNFTsRequest{BrandId: brandIDA, ClassId: classIDA, Owner: "invalidowner", Pagination: pageReq}},
		// brandID and classID requires each other
		{false, &types.QueryNFTsRequest{BrandId: brandIDA, ClassId: "", Owner: ownerA.String(), Pagination: pageReq}},
		{false, &types.QueryNFTsRequest{BrandId: "", ClassId: classIDA, Owner: ownerA.String(), Pagination: pageReq}},
	}

	for _, test := range tests {
		res, err := queryClient.NFTs(wrapCtx, test.req)
		if test.expectPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
			suite.Require().Nil(res.Pagination.NextKey)
			suite.Require().Len(res.Nfts, 0)
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(res)
		}
	}

	//ignore set brand first
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDB, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDC, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDB, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDC, classIDC, 0, types.NewClassDescription("", "", "", "")))

	nftData := []struct {
		num int
		b   string
		c   string
		o   sdk.AccAddress
	}{
		{3, brandIDA, classIDA, ownerA}, {10, brandIDA, classIDA, ownerB}, {2, brandIDA, classIDB, ownerA}, {2, brandIDA, classIDC, ownerA},
		{5, brandIDB, classIDA, ownerB}, {3, brandIDB, classIDB, ownerB}, {10, brandIDB, classIDB, ownerC},
		{10, brandIDC, classIDC, ownerC},
	}

	ownerMap := map[string]types.NFTs{}
	classMap := map[string]types.NFTs{}
	totalNftLen := 0

	//mint nfts
	for _, d := range nftData {
		if _, ok := ownerMap[d.o.String()]; !ok {
			ownerMap[d.o.String()] = types.NFTs{}
		}
		if _, ok := classMap[types.GetClassUniqueID(d.b, d.c)]; !ok {
			classMap[types.GetClassUniqueID(d.b, d.c)] = types.NFTs{}
		}

		for i := 0; i < d.num; i++ {
			nft, err := app.NFTKeeper.GenNFT(ctx, d.b, d.c, "ipfs://nft", "")
			suite.Require().NoError(err)
			suite.Require().NoError(app.NFTKeeper.MintNFT(ctx, nft, d.o))

			ownerMap[d.o.String()] = append(ownerMap[d.o.String()], nft)
			classMap[types.GetClassUniqueID(d.b, d.c)] = append(classMap[types.GetClassUniqueID(d.b, d.c)], nft)
			totalNftLen++
		}
	}

	req := &types.QueryNFTsRequest{}
	res, err := queryClient.NFTs(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().Nil(res.Pagination.NextKey)
	suite.Require().Len(res.Nfts, totalNftLen)
	suite.T().Logf("request nfts resultLen: %d", totalNftLen)

	//by owner
	for addr, nfts := range ownerMap {
		req = &types.QueryNFTsRequest{BrandId: "", ClassId: "", Owner: addr, Pagination: pageReq}
		res, err = queryClient.NFTs(wrapCtx, req)
		suite.Require().NoError(err)
		suite.Require().Len(res.Nfts, len(nfts))
		suite.T().Logf("request by only owner for %s, resultLen: %d", addr, len(nfts))
	}

	//by brand and class
	for uid, nfts := range classMap {
		brandID, classID := types.ParseClassUniqueID(uid)
		req = &types.QueryNFTsRequest{BrandId: brandID, ClassId: classID, Owner: "", Pagination: pageReq}
		res, err = queryClient.NFTs(wrapCtx, req)
		suite.Require().NoError(err)
		suite.Require().Len(res.Nfts, len(nfts))
		suite.T().Logf("request by brand and class for ID %s, %s, resultLen: %d", req.BrandId, req.ClassId, len(nfts))

		// and owner
		for owner, _ := range ownerMap {
			var filteredNfts types.NFTs
			for _, nft := range nfts {
				o := app.NFTKeeper.GetOwner(ctx, nft.BrandId, nft.ClassId, nft.Id)
				suite.Require().NotNil(o)
				if strings.EqualFold(owner, o.String()) {
					filteredNfts = append(filteredNfts, nft)
				}
			}
			req = &types.QueryNFTsRequest{BrandId: brandID, ClassId: classID, Owner: owner, Pagination: pageReq}
			res, err = queryClient.NFTs(wrapCtx, req)
			suite.Require().NoError(err)
			suite.Require().Len(res.Nfts, len(filteredNfts))
			suite.T().Logf("request by brand and class and owner for ID %s, %s, and address: %s resultLen: %d", req.BrandId, req.ClassId, owner, len(nfts))
		}
	}

	//check for pagination
	pageReq = &query.PageRequest{
		Key:        nil,
		Limit:      5,
		CountTotal: false,
	}
	req = &types.QueryNFTsRequest{BrandId: brandIDC, ClassId: classIDC, Owner: ownerC.String(), Pagination: pageReq}
	res, err = queryClient.NFTs(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res.Pagination.NextKey)
	suite.Require().Len(res.Nfts, 5)

	//next page
	pageReq = &query.PageRequest{
		Key:        res.Pagination.NextKey,
		Limit:      5,
		CountTotal: false,
	}
	req = &types.QueryNFTsRequest{BrandId: brandIDC, ClassId: classIDC, Owner: ownerC.String(), Pagination: pageReq}
	res, err = queryClient.NFTs(wrapCtx, req)
	suite.Require().NoError(err)
	suite.Require().Nil(res.Pagination.NextKey)
	suite.Require().Len(res.Nfts, 5)
}

func (suite *KeeperTestSuite) TestNFT() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB, brandIDC := "brandIDA", "brandIDB", "brandIDC"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB, ownerC := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb"), sdk.AccAddress("ownerc")

	//invalid arguments
	req := &types.QueryNFTRequest{BrandId: brandIDA, ClassId: classIDA, Id: 0}
	res, err := queryClient.NFT(wrapCtx, req)
	suite.Require().Error(err)
	req = &types.QueryNFTRequest{BrandId: brandIDA, ClassId: "", Id: 1}
	res, err = queryClient.NFT(wrapCtx, req)
	suite.Require().Error(err)
	req = &types.QueryNFTRequest{BrandId: "", ClassId: classIDA, Id: 1}
	res, err = queryClient.NFT(wrapCtx, req)
	suite.Require().Error(err)

	//ignore set brand first
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDB, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDC, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDB, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDC, classIDC, 0, types.NewClassDescription("", "", "", "")))

	nftData := []struct {
		num int
		b   string
		c   string
		o   sdk.AccAddress
	}{
		{3, brandIDA, classIDA, ownerA}, {10, brandIDA, classIDA, ownerB}, {2, brandIDA, classIDB, ownerA}, {2, brandIDA, classIDC, ownerA},
		{5, brandIDB, classIDA, ownerB}, {3, brandIDB, classIDB, ownerB}, {10, brandIDB, classIDB, ownerC},
		{10, brandIDC, classIDC, ownerC},
	}

	//[nftUniqueID]
	nftMap := map[string]string{}

	suite.Require().Zero(len(nftMap))
	//mint nfts
	for _, d := range nftData {
		for i := 0; i < d.num; i++ {
			nft, err := app.NFTKeeper.GenNFT(ctx, d.b, d.c, "ipfs://nft", "")
			suite.Require().NoError(err)
			suite.Require().NoError(app.NFTKeeper.MintNFT(ctx, nft, d.o))

			nftMap[types.GetNFTUniqueID(nft.BrandId, nft.ClassId, nft.Id)] = d.o.String()
		}
	}
	suite.Require().NotZero(len(nftMap))

	for uid, owner := range nftMap {
		brandID, classID, id, err := types.ParseNFTUniqueID(uid)
		suite.Require().NoError(err)

		req = &types.QueryNFTRequest{brandID, classID, id}

		res, err = queryClient.NFT(wrapCtx, req)
		suite.Require().NoError(err)

		suite.Require().Equal(brandID, res.Nft.BrandId)
		suite.Require().Equal(classID, res.Nft.ClassId)
		suite.Require().Equal(id, res.Nft.Id)
		acc, _ := sdk.AccAddressFromBech32(owner)
		suite.Require().Equal(app.NFTKeeper.GetOwner(ctx, brandID, classID, id), acc)
	}

	//not found
	req = &types.QueryNFTRequest{"random", "random", 1}
	res, err = queryClient.NFT(wrapCtx, req)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestOwner() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB, brandIDC := "brandIDA", "brandIDB", "brandIDC"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB, ownerC := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb"), sdk.AccAddress("ownerc")

	//invalid arguments
	req := &types.QueryOwnerRequest{BrandId: brandIDA, ClassId: classIDA, Id: 0}
	res, err := queryClient.Owner(wrapCtx, req)
	suite.Require().Error(err)
	req = &types.QueryOwnerRequest{BrandId: brandIDA, ClassId: "", Id: 1}
	res, err = queryClient.Owner(wrapCtx, req)
	suite.Require().Error(err)
	req = &types.QueryOwnerRequest{BrandId: "", ClassId: classIDA, Id: 1}
	res, err = queryClient.Owner(wrapCtx, req)
	suite.Require().Error(err)

	//ignore set brand first
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDB, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDC, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDB, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDC, classIDC, 0, types.NewClassDescription("", "", "", "")))

	nftData := []struct {
		num int
		b   string
		c   string
		o   sdk.AccAddress
	}{
		{3, brandIDA, classIDA, ownerA}, {10, brandIDA, classIDA, ownerB}, {2, brandIDA, classIDB, ownerA}, {2, brandIDA, classIDC, ownerA},
		{5, brandIDB, classIDA, ownerB}, {3, brandIDB, classIDB, ownerB}, {10, brandIDB, classIDB, ownerC},
		{10, brandIDC, classIDC, ownerC},
	}

	//[nftUniqueID]
	nftMap := map[string]string{}

	suite.Require().Zero(len(nftMap))
	//mint nfts
	for _, d := range nftData {
		for i := 0; i < d.num; i++ {
			nft, err := app.NFTKeeper.GenNFT(ctx, d.b, d.c, "ipfs://nft", "")
			suite.Require().NoError(err)
			suite.Require().NoError(app.NFTKeeper.MintNFT(ctx, nft, d.o))

			nftMap[types.GetNFTUniqueID(nft.BrandId, nft.ClassId, nft.Id)] = d.o.String()
		}
	}
	suite.Require().NotZero(len(nftMap))

	ownerD := sdk.AccAddress("ownerd").String()
	for uid, owner := range nftMap {
		brandID, classID, id, err := types.ParseNFTUniqueID(uid)
		suite.Require().NoError(err)

		req = &types.QueryOwnerRequest{brandID, classID, id}

		res, err = queryClient.Owner(wrapCtx, req)
		suite.Require().NoError(err)
		suite.Require().NotEmpty(res.Owner)
		suite.Require().NotEqual(owner, ownerD)
		suite.Require().Equal(owner, res.Owner)
		suite.Require().Equal(app.NFTKeeper.GetOwner(ctx, brandID, classID, id).String(), res.Owner)
	}
}

func (suite *KeeperTestSuite) TestSupply() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	wrapCtx := sdk.WrapSDKContext(ctx)

	brandIDA, brandIDB, brandIDC := "brandIDA", "brandIDB", "brandIDC"
	classIDA, classIDB, classIDC := "classIDA", "classIDB", "classIDC"
	ownerA, ownerB, ownerC := sdk.AccAddress("ownera"), sdk.AccAddress("ownerb"), sdk.AccAddress("ownerc")

	//not found
	req := &types.QuerySupplyRequest{BrandId: brandIDA, ClassId: classIDA}
	res, err := queryClient.Supply(wrapCtx, req)
	suite.Require().Error(err)
	//invalid arguments
	req = &types.QuerySupplyRequest{BrandId: brandIDA, ClassId: ""}
	res, err = queryClient.Supply(wrapCtx, req)
	suite.Require().Error(err)
	req = &types.QuerySupplyRequest{BrandId: "", ClassId: classIDA}
	res, err = queryClient.Supply(wrapCtx, req)
	suite.Require().Error(err)

	//ignore set brand first
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDB, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDA, classIDC, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDA, 0, types.NewClassDescription("", "", "", "")))
	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDB, classIDB, 0, types.NewClassDescription("", "", "", "")))

	app.NFTKeeper.SaveClass(ctx, types.NewClass(brandIDC, classIDC, 0, types.NewClassDescription("", "", "", "")))

	nftData := []struct {
		num int
		b   string
		c   string
		o   sdk.AccAddress
	}{
		{3, brandIDA, classIDA, ownerA}, {10, brandIDA, classIDA, ownerB}, {2, brandIDA, classIDB, ownerA}, {2, brandIDA, classIDC, ownerA},
		{5, brandIDB, classIDA, ownerB}, {3, brandIDB, classIDB, ownerB}, {10, brandIDB, classIDB, ownerC},
		{10, brandIDC, classIDC, ownerC},
	}

	supplyMap := map[string]*types.Supply{}
	nftMap := map[string]bool{}

	//mint nfts
	for _, d := range nftData {
		req = &types.QuerySupplyRequest{d.b, d.c}

		if _, ok := supplyMap[types.GetClassUniqueID(d.b, d.c)]; !ok {
			supply := types.DefaultSupply()
			supplyMap[types.GetClassUniqueID(d.b, d.c)] = &supply

			res, err = queryClient.Supply(wrapCtx, req)
			suite.Require().NoError(err)

			suite.Require().Equal(res.Amount, supply.TotalSupply)
			suite.Require().Equal(res.NextSequence, supply.Sequence)
		}
		supply := supplyMap[types.GetClassUniqueID(d.b, d.c)]
		prevSupply := *supply

		for i := 0; i < d.num; i++ {
			nft, err := app.NFTKeeper.GenNFT(ctx, d.b, d.c, "ipfs://nft", "")
			suite.Require().NoError(err)
			suite.Require().NoError(app.NFTKeeper.MintNFT(ctx, nft, d.o))

			supply.IncreaseSupply()
			suite.Require().NotEqual(prevSupply.TotalSupply, supply.TotalSupply)
			suite.Require().NotEqual(prevSupply.Sequence, supply.Sequence)

			res, err := queryClient.Supply(wrapCtx, req)
			suite.Require().NoError(err)
			suite.Require().Equal(res.Amount, supply.TotalSupply)
			suite.Require().Equal(res.NextSequence, supply.Sequence)

			supplyMap[types.GetClassUniqueID(d.b, d.c)] = supply
			nftMap[types.GetNFTUniqueID(nft.BrandId, nft.ClassId, nft.Id)] = true
		}
	}

	moduleTotalSupply := uint64(0)
	for _, supply := range supplyMap {
		moduleTotalSupply = moduleTotalSupply + supply.TotalSupply
	}

	suite.Require().Equal(moduleTotalSupply, uint64(len(nftMap)))

	//burrning
	for uid, _ := range nftMap {
		brandID, classID, id, err := types.ParseNFTUniqueID(uid)
		suite.Require().NoError(err)

		supply := supplyMap[types.GetClassUniqueID(brandID, classID)]

		suite.Require().NoError(app.NFTKeeper.BurnNFT(ctx, brandID, classID, id))

		req = &types.QuerySupplyRequest{brandID, classID}
		res, err := queryClient.Supply(wrapCtx, req)
		suite.Require().NoError(err)
		suite.Require().Equal(res.NextSequence, supply.Sequence)
		suite.Require().Equal(res.Amount, supply.TotalSupply-1)
		supply.DecreaseSupply()
	}

	for _, supply := range supplyMap {
		suite.Require().Zero(supply.TotalSupply)
	}

}
