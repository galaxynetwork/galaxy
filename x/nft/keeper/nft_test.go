package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

func (suite *KeeperTestSuite) TestStoreNFT() {
	keeper, ctx := suite.app.NFTKeeper, suite.ctx

	tests := []struct {
		nft    types.NFT
		owner  sdk.AccAddress
		stored bool
	}{
		{types.NewNFT(0, "brandid", "classid", "ipfs://nft", "ipfs://nft"), sdk.AccAddress("addr1"), true},
		{types.NewNFT(0, "brandid", "classid", "ipfs://nft", "ipfs://nft"), sdk.AccAddress("addr1"), true},
		{types.NewNFT(0, "brandid", "classid", "ipfs://nft", "ipfs://nft"), sdk.AccAddress("addr2"), true},
		{types.NewNFT(0, "brandid", "classid2", "ipfs://nft", "ipfs://nft"), sdk.AccAddress("addr2"), true},
		{types.NewNFT(0, "brandid2", "classid", "ipfs://nft", "ipfs://nft"), sdk.AccAddress("addr3"), true},
		{types.NewNFT(0, "brandid2", "classid", "ipfs://nft", "ipfs://nft"), sdk.AccAddress("addr4"), false},
	}

	var classes = map[string]uint64{}
	var ownerMap = map[string]int{}
	for i, test := range tests {
		nft, owner := test.nft, test.owner

		_, ok := ownerMap[owner.String()]
		if !ok {
			ownerMap[owner.String()] = 0
		}
		if test.stored {
			ownerMap[owner.String()]++
		}

		suite.Require().Error(nft.Validate())

		var genNFT types.NFT
		var err error

		if !keeper.HasClass(ctx, nft.BrandId, nft.ClassId) {
			genNFT, err = keeper.GenNFT(ctx, nft.BrandId, nft.ClassId, nft.Uri, nft.VarUri)
			suite.Require().Error(err)

			suite.Require().NoError(
				keeper.SaveClass(ctx, types.NewClass(nft.BrandId, nft.ClassId, 0, types.NewClassDescription("", "", "", ""))),
			)
		}

		genNFT, err = keeper.GenNFT(ctx, nft.BrandId, nft.ClassId, nft.Uri, nft.VarUri)
		suite.Require().NoError(err)

		nft = genNFT
		tests[i].nft = nft

		//mint
		if test.stored {
			suite.Require().NoError(
				keeper.MintNFT(ctx, nft, owner),
			)
			classes[types.GetClassUniqueID(nft.BrandId, nft.ClassId)]++
		}
	}

	// has/get
	for _, test := range tests {
		nft := test.nft

		hasExist := keeper.HasNFT(ctx, nft.BrandId, nft.ClassId, nft.Id)
		gnft, exist := keeper.GetNFT(ctx, nft.BrandId, nft.ClassId, nft.Id)
		if test.stored {
			suite.Require().True(hasExist)
			suite.Require().True(exist)
			suite.Require().Equal(nft, gnft)
		} else {
			suite.Require().False(hasExist)
			suite.Require().False(exist)
			suite.Require().Empty(gnft)
		}
	}

	// check nfts len by owner
	for k, v := range ownerMap {
		acc, _ := sdk.AccAddressFromBech32(k)
		nfts := keeper.GetNFTsByOwner(ctx, acc)
		suite.Require().Len(nfts, v, "check nfts len by owner for: %s", k)
	}

	//update
	newVarUri := "ifps://netnft"
	for _, test := range tests {
		nft := test.nft

		if test.stored {
			gnft, _ := keeper.GetNFT(ctx, nft.BrandId, nft.ClassId, nft.Id)
			suite.Require().NotEqual(gnft.VarUri, newVarUri)

			suite.Require().NoError(keeper.UpdateNFT(ctx, nft.BrandId, nft.ClassId, nft.Id, newVarUri))

			gnft, _ = keeper.GetNFT(ctx, nft.BrandId, nft.ClassId, nft.Id)
			suite.Require().Equal(gnft.VarUri, newVarUri)
		} else {
			suite.Require().Error(keeper.UpdateNFT(ctx, nft.BrandId, nft.ClassId, nft.Id, newVarUri))

		}
	}

	//transfer
	for i, test := range tests {
		nft := test.nft
		owner := keeper.GetOwner(ctx, nft.BrandId, nft.ClassId, nft.Id)
		newOwner := sdk.AccAddress(fmt.Sprintf("newaddr%d", i))

		//for check owner store
		if _, ok := ownerMap[newOwner.String()]; !ok {
			ownerMap[newOwner.String()] = 0
		}

		if test.stored {
			ownerMap[owner.String()]--
			suite.Require().NotNil(owner)
			suite.Require().NotEqual(owner, newOwner)
			suite.Require().NoError(
				keeper.TransferNFT(ctx, nft.BrandId, nft.ClassId, nft.Id, owner, newOwner),
			)
			owner = keeper.GetOwner(ctx, nft.BrandId, nft.ClassId, nft.Id)
			suite.Require().Equal(owner, newOwner)
			ownerMap[newOwner.String()]++
		} else {
			suite.Require().Nil(owner)
			suite.Require().Error(
				keeper.TransferNFT(ctx, nft.BrandId, nft.ClassId, nft.Id, owner, newOwner),
			)
		}
	}

	// check nfts len by owner after transfer
	for k, v := range ownerMap {
		acc, _ := sdk.AccAddressFromBech32(k)
		nfts := keeper.GetNFTsByOwner(ctx, acc)
		suite.Require().Len(nfts, v, "check nfts len by owner after transferring for: %s", k)
	}

	for key, count := range classes {
		brandID, classID := types.ParseClassUniqueID(key)
		supply, err := keeper.GetTotalSupplyOfClass(ctx, brandID, classID)
		suite.Require().NoError(err)
		suite.Require().Equal(supply, count)
	}

	//burn
	for _, test := range tests {
		nft := test.nft

		burned := keeper.BurnedNFT(ctx, nft.BrandId, nft.ClassId, nft.Id)
		suite.Require().False(burned)

		if test.stored {
			suite.Require().NoError(
				keeper.BurnNFT(ctx, nft.BrandId, nft.ClassId, nft.Id),
			)

			burned = keeper.BurnedNFT(ctx, nft.BrandId, nft.ClassId, nft.Id)
			suite.Require().True(burned)
		} else {
			suite.Require().Error(
				keeper.BurnNFT(ctx, nft.BrandId, nft.ClassId, nft.Id),
			)
			burned = keeper.BurnedNFT(ctx, nft.BrandId, nft.ClassId, nft.Id)
			suite.Require().False(burned)
		}
	}

	// check nfts len by owner after burning
	for k, _ := range ownerMap {
		acc, _ := sdk.AccAddressFromBech32(k)
		nfts := keeper.GetNFTsByOwner(ctx, acc)
		suite.Require().Len(nfts, 0, "check nfts len by owner after burning for: %s", k)
	}

	// check totalSupply after burning
	for key, _ := range classes {
		brandID, classID := types.ParseClassUniqueID(key)
		supply, err := keeper.GetTotalSupplyOfClass(ctx, brandID, classID)
		suite.Require().NoError(err)
		suite.Require().Equal(uint64(0), supply)
	}
}
