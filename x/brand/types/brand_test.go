package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestBrandValidate(t *testing.T) {
	tests := []struct {
		brand      Brand
		expectPass bool
	}{
		{
			brand: NewBrand(
				"brand",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "details", "https://image.png"),
			),
			expectPass: true,
		},
		{
			brand: NewBrand(
				"brand1",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "", ""),
			),
			expectPass: true,
		},
		{
			brand: NewBrand(
				"1brand",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "", ""),
			),
			expectPass: true,
		},
		{
			brand: NewBrand(
				"123",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "", ""),
			),
			expectPass: true,
		},
		{
			brand: NewBrand(
				"1",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "", ""),
			),
			expectPass: true,
		},
		{
			brand: NewBrand(
				"123456789123456789123456789123",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "", ""),
			),
			expectPass: true,
		},
		{
			brand: NewBrand(
				"-brand",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "", ""),
			),
			expectPass: false,
		},
		{
			brand: NewBrand(
				"",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "", ""),
			),
			expectPass: false,
		},
		{
			brand: NewBrand(
				"br-and",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "", ""),
			),
			expectPass: false,
		},
		{
			brand: NewBrand(
				"brand12345678912345678912345678",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("name", "", ""),
			),
			expectPass: false,
		},
		{
			brand: NewBrand(
				"brand",
				sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()),
				NewBrandDescription("", "", ""),
			),
			expectPass: false,
		},
		//invalid owner address
		{
			brand: NewBrand(
				"brand",
				[]byte(""),
				NewBrandDescription("name", "", ""),
			),
			expectPass: false,
		},
		//invalid brand address
		{
			brand: Brand{
				Id:           "brand",
				Owner:        sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()).String(),
				Description:  NewBrandDescription("brand", "", ""),
				BrandAddress: "",
			},
			expectPass: false,
		},
		{
			brand: Brand{
				Id:           "brand",
				Owner:        sdk.AccAddress(ed25519.GenPrivKey().PubKey().String()).String(),
				Description:  NewBrandDescription("brand", "", ""),
				BrandAddress: sdk.AccAddress(NewBrandAddress("diffrentbrandid")).String(),
			},
			expectPass: false,
		},
	}

	for i, test := range tests {
		if test.expectPass {
			require.NoError(t, test.brand.Validate(), "test brand index: %d, brand: %s", i, test.brand.String())
		} else {
			require.Error(t, test.brand.Validate(), "test brand index: %d, brand: %s", i, test.brand.String())
		}
	}
}
