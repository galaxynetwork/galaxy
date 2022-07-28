package types

import (
	"testing"

	"github.com/galaxies-labs/galaxy/internal/util"
	"github.com/stretchr/testify/require"
)

func TestClass(t *testing.T) {
	tests := []struct {
		class      Class
		expectPass bool
	}{
		{NewClass("brandid", "classid", 10_000, NewClassDescription("name", "details", "externalurl", "imageuri")), true},
		{NewClass("brandid", "classid", 10_000, NewClassDescription("", "", "", "")), true},
		{NewClass("brandid", "class-id", 10_000, NewClassDescription("", "", "", "")), true},
		{NewClass("brandid", "class-id123", 10_000, NewClassDescription("", "", "", "")), true},
		{NewClass("brandid", "classid", 1, NewClassDescription("", "", "", "")), true},
		{NewClass("brandid", "classid", 10_000, NewClassDescription("name", "details", "externalurl", "imageuri")), true},
		{NewClass("brandid", "classid", 10_000, NewClassDescription(
			util.GenStringWithLength(MaxClassNameLength), util.GenStringWithLength(MaxClassDetailsLength),
			util.GenStringWithLength(MaxUriLength), util.GenStringWithLength(MaxUriLength))), true},
		// empty brand id
		{NewClass("", "classid", 10_000, NewClassDescription("", "", "", "")), false},
		// invalid class id length
		{NewClass("brandid", "cl", 10_000, NewClassDescription("", "", "", "")), false},
		// fee basis points over max value
		{NewClass("brandid", "classid", 10_002, NewClassDescription("", "", "", "")), false},
		// starts with hypen
		{NewClass("brandid", "-classid", 10_000, NewClassDescription("", "", "", "")), false},
		// invalid length description
		{NewClass("brandid", "classid", 10_000, NewClassDescription(util.GenStringWithLength(MaxClassNameLength+1), "details", "externalurl", "imageuri")), false},
		{NewClass("brandid", "classid", 10_000, NewClassDescription("name", util.GenStringWithLength(MaxClassDetailsLength+1), "externalurl", "imageuri")), false},
		{NewClass("brandid", "classid", 10_000, NewClassDescription("name", "details", util.GenStringWithLength(MaxUriLength+1), "imageuri")), false},
		{NewClass("brandid", "classid", 10_000, NewClassDescription("name", "details", "externalurl", util.GenStringWithLength(MaxUriLength+1))), false},
	}

	for index, test := range tests {
		err := test.class.Validate()
		if test.expectPass {
			require.NoErrorf(t, err, "index: %d", index)
		} else {
			require.Errorf(t, err, "index: %d", index)
		}
	}
}

func TestNFT(t *testing.T) {
	tests := []struct {
		expectPass bool
		nft        NFT
	}{
		{true, NewNFT(1, "id", "classId", "ipfs://hash", "")},
		{true, NewNFT(2, "id", "classId", util.GenStringWithLength(MaxUriLength), util.GenStringWithLength(MaxUriLength))},
		//zero nft id
		{false, NewNFT(0, "id", "classId", "ipfs://hash", "")},
		//invalid brand id
		{false, NewNFT(1, "", "classId", "ipfs://hash", "")},
		//invalid class id
		{false, NewNFT(1, "id", "", "ipfs://hash", "")},
		//invalid uri
		{false, NewNFT(1, "id", "classid", "", "")},
		//invalid uri length
		{false, NewNFT(1, "id", "classid", util.GenStringWithLength(MaxUriLength+1), util.GenStringWithLength(MaxUriLength+1))},
	}

	for index, test := range tests {
		err := test.nft.Validate()
		if test.expectPass {
			require.NoErrorf(t, err, "index: %d", index)
		} else {
			require.Errorf(t, err, "index: %d", index)
		}
	}
}

func TestSupply(t *testing.T) {
	supply := DefaultSupply()

	require.Equal(t, supply.Sequence, uint64(1))
	require.Equal(t, supply.TotalSupply, uint64(0))

	supply.DecreaseSupply()

	require.Equal(t, supply.Sequence, uint64(1))
	require.Equal(t, supply.TotalSupply, uint64(0))

	var lastSequence uint64
	var i uint64
	for i = 1; i <= 100; i++ {
		// save nft
		require.Equal(t, supply.Sequence, i)
		require.Equal(t, supply.TotalSupply, i-1)

		supply.IncreaseSupply()

		require.Equal(t, supply.TotalSupply, i)
		lastSequence = supply.Sequence
	}

	for supply.TotalSupply != 0 {
		currentSupply := supply.TotalSupply

		// burn nft
		supply.DecreaseSupply()

		require.Equal(t, supply.Sequence, lastSequence)
		require.Equal(t, supply.TotalSupply, currentSupply-1)

	}

	require.Equal(t, supply.Sequence, lastSequence)
	require.Equal(t, supply.TotalSupply, uint64(0))

}
