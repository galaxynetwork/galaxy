package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

//only check for genesis validate
//detail type checking of Class and NFT on nft_test.go
func TestGenesis(t *testing.T) {
	classA := NewClass("brandid", "classid", 10_000, NewClassDescription("name", "", "", ""))
	classA2 := NewClass("brandid", "classid", 10_000, NewClassDescription("name", "", "", ""))
	classB := NewClass("brandid2", "classid", 10_000, NewClassDescription("name", "", "", ""))
	invalidClassA := NewClass("brandidc", "", 10_000, NewClassDescription("name", "", "", ""))
	invalidClassB := NewClass("", "classida", 10_000, NewClassDescription("name", "", "", ""))

	validNFTA := NewNFT(1, classA.BrandId, classA.Id, "", "")
	validNFTA2 := validNFTA
	validNFTB := NewNFT(1, classB.BrandId, classB.Id, "", "")
	validNFTC := NewNFT(2, classA.BrandId, classA.Id, "", "")
	invalidNFTA := NewNFT(1, invalidClassA.BrandId, classA.Id, "", "")
	invalidNFTB := NewNFT(1, classA.BrandId, invalidClassB.Id, "", "")

	ownerA := sdk.AccAddress("ownera...")
	ownerB := sdk.AccAddress("ownerb...")

	tests := []struct {
		expectPass   bool
		genesisState *GenesisState
	}{
		{true, DefaultGenesisState()},
		{true, NewGenesisState(Classes{classA, classB}, Entries{{ownerA.String(), NFTs{validNFTA, validNFTB}}, {ownerB.String(), NFTs{validNFTC}}})},
		// duplicate class
		{false, NewGenesisState(Classes{classA, classA2, classB}, Entries{})},
		// invalid brandID format
		{false, NewGenesisState(Classes{classA, classB, invalidClassB}, Entries{})},
		// invalid classID format
		{false, NewGenesisState(Classes{classA, classB, invalidClassA}, Entries{})},
		// nft not within brandID
		{false, NewGenesisState(Classes{classA, classB}, Entries{{ownerA.String(), NFTs{validNFTA, invalidNFTA}}})},
		// nft not within brandID/classID
		{false, NewGenesisState(Classes{classA, classB}, Entries{{ownerA.String(), NFTs{validNFTA, invalidNFTB}}})},
		// duplicate nft within brandID/classID
		{false, NewGenesisState(Classes{classA, classB}, Entries{{ownerA.String(), NFTs{validNFTB}}, {ownerB.String(), NFTs{validNFTA, validNFTA2}}})},
		// invalid owner address
		{false, NewGenesisState(Classes{classA, classB}, Entries{{ownerA.String(), NFTs{validNFTA, validNFTB}}, {"", NFTs{validNFTC}}})},
	}

	for index, test := range tests {
		err := test.genesisState.Validate()
		if test.expectPass {
			require.NoErrorf(t, err, "test for index: %d", index)
		} else {
			require.Errorf(t, err, "test for index: %d", index)
		}
	}

}
