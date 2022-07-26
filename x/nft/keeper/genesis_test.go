package keeper_test

import (
	brandtypes "github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestInitGenesis() {
	brandOwnerA := sdk.AccAddress("brandOwner1")
	brandOwnerB := sdk.AccAddress("brandOwner2")

	ownerA := sdk.AccAddress("owner1")
	ownerB := sdk.AccAddress("owner2")

	brandA := brandtypes.NewBrand("brandid", brandOwnerA, brandtypes.NewBrandDescription("name", "", ""))
	brandB := brandtypes.NewBrand("brandid2", brandOwnerB, brandtypes.NewBrandDescription("name", "", ""))

	classA := types.NewClass(brandA.Id, "classid", 0, types.NewClassDescription("", "", "", ""))
	classA2 := types.NewClass(brandA.Id, "classid2", 0, types.NewClassDescription("", "", "", ""))
	classB := types.NewClass(brandB.Id, "classid", 0, types.NewClassDescription("", "", "", ""))
	classB2 := types.NewClass(brandB.Id, "classid2", 0, types.NewClassDescription("", "", "", ""))

	nftA := types.NewNFT(1, classA.BrandId, classA.Id, "ipfs://nft", "")
	nftA2 := types.NewNFT(2, classA.BrandId, classA.Id, "ipfs://nft", "")
	nftA3 := types.NewNFT(10, classA.BrandId, classA.Id, "ipfs://nft", "")

	nftB := types.NewNFT(1, classB.BrandId, classB.Id, "ipfs://nft", "")

	tests := []struct {
		expectPass bool
		storeBrand bool
		genesis    *types.GenesisState
	}{
		{true, true, types.DefaultGenesisState()},
		// no brand
		{false, false, types.NewGenesisState(
			types.ClassEntries{
				types.ClassEntry{Class: classA, NextSequence: 1},
				types.ClassEntry{Class: classB, NextSequence: 1},
			},
			types.Entries{},
		)},
		{true, true, types.NewGenesisState(
			types.ClassEntries{
				types.ClassEntry{Class: classA, NextSequence: 1},
				types.ClassEntry{Class: classB, NextSequence: 1},
			},
			types.Entries{},
		)},
		// same class
		{false, true, types.NewGenesisState(
			types.ClassEntries{
				types.ClassEntry{Class: classA, NextSequence: 1},
				types.ClassEntry{Class: classB, NextSequence: 1},
				types.ClassEntry{Class: classA, NextSequence: 1},
			},
			types.Entries{},
		)},
		{true, true, types.NewGenesisState(
			types.ClassEntries{
				types.ClassEntry{Class: classA, NextSequence: 3},
				types.ClassEntry{Class: classB, NextSequence: 2},
				types.ClassEntry{Class: classA2, NextSequence: 1},
				types.ClassEntry{Class: classB2, NextSequence: 1},
			},
			types.Entries{
				types.Entry{Owner: ownerA.String(), Nfts: types.NFTs{nftA, nftA2}},
				types.Entry{Owner: ownerB.String(), Nfts: types.NFTs{nftB}},
			},
		)},
		// out of sequence
		{false, true, types.NewGenesisState(
			types.ClassEntries{
				types.ClassEntry{Class: classA, NextSequence: 3},
				types.ClassEntry{Class: classB, NextSequence: 2},
				types.ClassEntry{Class: classA2, NextSequence: 1},
				types.ClassEntry{Class: classB2, NextSequence: 1},
			},
			types.Entries{
				types.Entry{Owner: ownerA.String(), Nfts: types.NFTs{nftA, nftA2, nftA3}},
				types.Entry{Owner: ownerB.String(), Nfts: types.NFTs{nftB}},
			},
		)},
		{false, true, types.NewGenesisState(
			types.ClassEntries{
				types.ClassEntry{Class: classA, NextSequence: 10},
				types.ClassEntry{Class: classB, NextSequence: 2},
				types.ClassEntry{Class: classA2, NextSequence: 1},
				types.ClassEntry{Class: classB2, NextSequence: 1},
			},
			types.Entries{
				types.Entry{Owner: ownerA.String(), Nfts: types.NFTs{nftA, nftA2, nftA3}},
				types.Entry{Owner: ownerB.String(), Nfts: types.NFTs{nftB}},
			},
		)},
		//same nft
		{false, true, types.NewGenesisState(
			types.ClassEntries{
				types.ClassEntry{Class: classA, NextSequence: 11},
				types.ClassEntry{Class: classB, NextSequence: 2},
				types.ClassEntry{Class: classA2, NextSequence: 1},
				types.ClassEntry{Class: classB2, NextSequence: 1},
			},
			types.Entries{
				types.Entry{Owner: ownerA.String(), Nfts: types.NFTs{nftA, nftA2, nftA3}},
				types.Entry{Owner: ownerB.String(), Nfts: types.NFTs{nftB, nftB}},
			},
		)},

		// check for supply and owner length
		{true, true, types.NewGenesisState(
			types.ClassEntries{
				types.ClassEntry{Class: classA, NextSequence: 11},
				types.ClassEntry{Class: classB, NextSequence: 2},
				types.ClassEntry{Class: classA2, NextSequence: 1},
				types.ClassEntry{Class: classB2, NextSequence: 1},
			},
			types.Entries{
				types.Entry{Owner: ownerA.String(), Nfts: types.NFTs{nftA, nftA2, nftA3}},
				types.Entry{Owner: ownerB.String(), Nfts: types.NFTs{nftB}},
			},
		)},
	}

	for i, test := range tests {
		if test.storeBrand {
			suite.Require().NoError(
				suite.app.BrandKeeper.SetBrand(suite.ctx, brandA),
			)
			suite.Require().NoError(
				suite.app.BrandKeeper.SetBrand(suite.ctx, brandB),
			)
		}

		err := suite.app.NFTKeeper.InitGenesis(suite.ctx, test.genesis)
		if test.expectPass {
			suite.Require().NoErrorf(err, "test for index: %d", i)
		} else {
			suite.Require().Errorf(err, "test for index: %d", i)
		}

		if len(tests) == i+1 {
			//classs length check
			//nft length check
			//all nft owner check
			classes := suite.app.NFTKeeper.GetClasses(suite.ctx)
			suite.Require().Len(
				classes,
				len(test.genesis.ClassEntries),
			)

			for _, class := range classes {
				nfts := suite.app.NFTKeeper.GetNFTsOfClass(suite.ctx, class.BrandId, class.Id)
				totalSupply, err := suite.app.NFTKeeper.GetTotalSupplyOfClass(suite.ctx, class.BrandId, class.Id)
				suite.Require().NoError(err)
				suite.Require().Equal(uint64(len(nfts)), totalSupply)
			}
		}

		suite.SetupTest()
	}
}

func (suite *KeeperTestSuite) TestExportGenesis() {
	brandOwnerA := sdk.AccAddress("brandOwner1")
	brandOwnerB := sdk.AccAddress("brandOwner2")

	ownerA := sdk.AccAddress("owner1")
	ownerB := sdk.AccAddress("owner2")

	brandA := brandtypes.NewBrand("brandid", brandOwnerA, brandtypes.NewBrandDescription("name", "", ""))
	brandB := brandtypes.NewBrand("brandid2", brandOwnerB, brandtypes.NewBrandDescription("name", "", ""))

	classA := types.NewClass(brandA.Id, "classid", 0, types.NewClassDescription("", "", "", ""))
	classA2 := types.NewClass(brandA.Id, "classid2", 0, types.NewClassDescription("", "", "", ""))
	classB := types.NewClass(brandB.Id, "classid", 0, types.NewClassDescription("", "", "", ""))

	nftA := types.NewNFT(1, classA.BrandId, classA.Id, "ipfs://nft", "")
	nftA2 := types.NewNFT(2, classA.BrandId, classA.Id, "ipfs://nft", "")
	nftA3 := types.NewNFT(3, classA.BrandId, classA.Id, "ipfs://nft", "")
	nftA4 := types.NewNFT(4, classA.BrandId, classA.Id, "ipfs://nft", "")
	nftB := types.NewNFT(1, classB.BrandId, classB.Id, "ipfs://nft", "")

	//save brand
	suite.Require().NoError(suite.app.BrandKeeper.SetBrand(suite.ctx, brandA))
	suite.Require().NoError(suite.app.BrandKeeper.SetBrand(suite.ctx, brandB))

	//save class
	suite.Require().NoError(suite.app.NFTKeeper.SaveClass(suite.ctx, classA))
	suite.Require().NoError(suite.app.NFTKeeper.SaveClass(suite.ctx, classA2))
	suite.Require().NoError(suite.app.NFTKeeper.SaveClass(suite.ctx, classB))

	//save nft
	suite.Require().NoError(suite.app.NFTKeeper.MintNFT(suite.ctx, nftA, ownerA))
	suite.Require().NoError(suite.app.NFTKeeper.MintNFT(suite.ctx, nftA2, ownerA))
	suite.Require().NoError(suite.app.NFTKeeper.MintNFT(suite.ctx, nftA3, ownerA))
	suite.Require().NoError(suite.app.NFTKeeper.MintNFT(suite.ctx, nftA4, ownerA))
	suite.Require().NoError(suite.app.NFTKeeper.MintNFT(suite.ctx, nftB, ownerB))

	state := suite.app.NFTKeeper.ExportGenesis(suite.ctx)

	for _, ce := range state.ClassEntries {
		uniqueID := types.GetClassUniqueID(ce.Class.BrandId, ce.Class.Id)
		totalSupply, err := suite.app.NFTKeeper.GetTotalSupplyOfClass(suite.ctx, ce.Class.BrandId, ce.Class.Id)
		suite.Require().NoError(err)

		switch uniqueID {
		case types.GetClassUniqueID(classA.BrandId, classA.Id):
			suite.Require().Equal(ce.NextSequence, uint64(5))
			suite.Require().Equal(totalSupply, uint64(4))
			suite.Require().Equal(ce.Class, classA)
		case types.GetClassUniqueID(classA2.BrandId, classA2.Id):
			suite.Require().Equal(ce.NextSequence, uint64(1))
			suite.Require().Equal(totalSupply, uint64(0))

			suite.Require().Equal(ce.Class, classA2)
		case types.GetClassUniqueID(classB.BrandId, classB.Id):
			suite.Require().Equal(ce.NextSequence, uint64(2))
			suite.Require().Equal(totalSupply, uint64(1))
			suite.Require().Equal(ce.Class, classB)
		default:
			suite.FailNow("exported invalid class for uniqueID: %s", uniqueID)
		}
	}

	var nfts types.NFTs
	for _, e := range state.Entries {
		owner, err := sdk.AccAddressFromBech32(e.Owner)
		suite.Require().NoError(err)

		//check get nfts by owner
		switch {
		case owner.Equals(ownerA):
		case owner.Equals(ownerB):
		default:
			suite.FailNow("exported owner for address: %s", e.Owner)
		}
		for _, nft := range e.Nfts {
			nfts = append(nfts, nft)
		}
	}

	suite.Require().Len(nfts, 5)
	for _, nft := range nfts {
		uniqueID := types.GetNFTUniqueID(nft.BrandId, nft.ClassId, nft.Id)
		switch uniqueID {
		case types.GetNFTUniqueID(nftA.BrandId, nftA.ClassId, nftA.Id):
			suite.Require().Equal(nft, nftA)
		case types.GetNFTUniqueID(nftA2.BrandId, nftA2.ClassId, nftA2.Id):
			suite.Require().Equal(nft, nftA2)
		case types.GetNFTUniqueID(nftA3.BrandId, nftA3.ClassId, nftA3.Id):
			suite.Require().Equal(nft, nftA3)
		case types.GetNFTUniqueID(nftA4.BrandId, nftA4.ClassId, nftA4.Id):
			suite.Require().Equal(nft, nftA4)
		case types.GetNFTUniqueID(nftB.BrandId, nftB.ClassId, nftB.Id):
			suite.Require().Equal(nft, nftB)
		default:
			suite.Failf("failed", "exported invalid nft for uniqueID: %s", uniqueID)
		}
	}

	//burn
	suite.Require().NoError(
		suite.app.NFTKeeper.BurnNFT(suite.ctx, nftA2.BrandId, nftA2.ClassId, nftA2.Id),
	)

	state = suite.app.NFTKeeper.ExportGenesis(suite.ctx)

	nfts = types.NFTs{}
	for _, e := range state.Entries {
		for _, nft := range e.Nfts {
			nfts = append(nfts, nft)
		}
	}

	suite.Require().Len(nfts, 4)
	for _, nft := range nfts {
		uniqueID := types.GetNFTUniqueID(nft.BrandId, nft.ClassId, nft.Id)
		switch uniqueID {
		case types.GetNFTUniqueID(nftA.BrandId, nftA.ClassId, nftA.Id):
			suite.Require().Equal(nft, nftA)
		case types.GetNFTUniqueID(nftA2.BrandId, nftA2.ClassId, nftA2.Id):
			suite.Failf("failed", "exported burned nft for uniqueID: %s", uniqueID)
		case types.GetNFTUniqueID(nftA3.BrandId, nftA3.ClassId, nftA3.Id):
			suite.Require().Equal(nft, nftA3)
		case types.GetNFTUniqueID(nftA4.BrandId, nftA4.ClassId, nftA4.Id):
			suite.Require().Equal(nft, nftA4)
		case types.GetNFTUniqueID(nftB.BrandId, nftB.ClassId, nftB.Id):
			suite.Require().Equal(nft, nftB)
		default:
			suite.Failf("failed", "exported invalid nft for uniqueID: %s", uniqueID)
		}
	}

	for _, ce := range state.ClassEntries {
		totalSupply, err := suite.app.NFTKeeper.GetTotalSupplyOfClass(suite.ctx, ce.Class.BrandId, ce.Class.Id)

		switch types.GetClassUniqueID(ce.Class.BrandId, ce.Class.Id) {
		case types.GetClassUniqueID(classA.BrandId, classA.Id):
			suite.Require().NoError(err)
			//check for same sequence after burnning
			suite.Require().Equal(ce.NextSequence, uint64(5))
			//check for decrease totalSupply after burnning
			suite.Require().Equal(totalSupply, uint64(3))
		}
	}

	//mint new
	nftC, err := suite.app.NFTKeeper.GenNFT(suite.ctx, nftA.BrandId, nftA.ClassId, "ipfs://nft", "")
	suite.Require().NoError(err)
	suite.Require().Equal(nftC.Id, uint64(5))

	suite.Require().NoError(suite.app.NFTKeeper.MintNFT(suite.ctx, nftC, ownerA))

	state = suite.app.NFTKeeper.ExportGenesis(suite.ctx)

	for _, ce := range state.ClassEntries {
		totalSupply, err := suite.app.NFTKeeper.GetTotalSupplyOfClass(suite.ctx, ce.Class.BrandId, ce.Class.Id)

		switch types.GetClassUniqueID(ce.Class.BrandId, ce.Class.Id) {
		case types.GetClassUniqueID(classA.BrandId, classA.Id):
			suite.Require().NoError(err)
			// check for mint after burnning
			suite.Require().Equal(ce.NextSequence, uint64(6))
			suite.Require().Equal(totalSupply, uint64(4))
		}
	}

	nfts = types.NFTs{}
	for _, e := range state.Entries {
		for _, nft := range e.Nfts {
			nfts = append(nfts, nft)
		}
	}

	suite.Require().Len(nfts, 5)
}
