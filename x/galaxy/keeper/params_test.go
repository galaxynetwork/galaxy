package keeper_test

import (
	"testing"

	testkeeper "github.com/galaxies-labs/galaxy/testutil/keeper"
	"github.com/galaxies-labs/galaxy/x/galaxy/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.GalaxyKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
