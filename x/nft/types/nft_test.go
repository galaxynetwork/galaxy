package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func genStringWithLength(length int) string {
	bz := []byte{}
	for i := 0; i < length; i++ {
		bz = append(bz, byte(i))
	}
	return string(bz[:])
}

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
			genStringWithLength(MaxClassNameLength), genStringWithLength(MaxClassDetailsLength),
			genStringWithLength(MaxUriLength), genStringWithLength(MaxUriLength))), true},
		// empty brand id
		{NewClass("", "classid", 10_000, NewClassDescription("", "", "", "")), false},
		// invalid class id length
		{NewClass("brandid", "cl", 10_000, NewClassDescription("", "", "", "")), false},
		// fee basis points over max value
		{NewClass("brandid", "classid", 10_001, NewClassDescription("", "", "", "")), false},
		// starts with hypen
		{NewClass("brandid", "-classid", 10_000, NewClassDescription("", "", "", "")), false},
		// invalid length description
		{NewClass("brandid", "classid", 10_000, NewClassDescription(genStringWithLength(MaxClassNameLength+1), "details", "externalurl", "imageuri")), false},
		{NewClass("brandid", "classid", 10_000, NewClassDescription("name", genStringWithLength(MaxClassDetailsLength+1), "externalurl", "imageuri")), false},
		{NewClass("brandid", "classid", 10_000, NewClassDescription("name", "details", genStringWithLength(MaxUriLength+1), "imageuri")), false},
		{NewClass("brandid", "classid", 10_000, NewClassDescription("name", "details", "externalurl", genStringWithLength(MaxUriLength+1))), false},
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
		{true, NewNFT(1, "id", "classId", "", "")},
		{true, NewNFT(1, "id", "classId", genStringWithLength(MaxUriLength), genStringWithLength(MaxUriLength))},
		//zero nft id
		{false, NewNFT(0, "id", "classId", "", "")},
		//invalid brand id
		{false, NewNFT(1, "", "classId", "", "")},
		//invalid class id
		{false, NewNFT(1, "id", "", "", "")},
		//invalid uri length
		{false, NewNFT(1, "id", "classid", genStringWithLength(MaxUriLength+1), genStringWithLength(MaxUriLength+1))},
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
