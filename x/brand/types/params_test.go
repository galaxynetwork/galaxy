package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	tests := []struct {
		params     Params
		expectPass bool
	}{
		{
			params:     NewParams(sdk.NewCoin("ustake", sdk.NewInt(0))),
			expectPass: false,
		},
		{
			params:     NewParams(sdk.Coin{}),
			expectPass: false,
		},
		{
			params:     NewParams(sdk.NewCoin(DefaultBrandCreationFeeDenom, sdk.NewInt(0))),
			expectPass: true,
		},
		{
			params:     NewParams(sdk.NewCoin(DefaultBrandCreationFeeDenom, sdk.NewInt(5_000_000))),
			expectPass: true,
		},
		{
			params:     DefaultParams(),
			expectPass: true,
		},
	}

	for i, test := range tests {
		if test.expectPass {
			require.NoError(t, test.params.Validate(), "test params index: %d, params: %s", i, test.params.String())
		} else {
			require.Error(t, test.params.Validate(), "test params index: %d, params: %s", i, test.params.String())
		}
	}
}
