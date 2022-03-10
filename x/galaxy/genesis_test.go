package galaxy_test

import (
	"testing"

	keepertest "github.com/galaxies-labs/galaxy/testutil/keeper"
	"github.com/galaxies-labs/galaxy/testutil/nullify"
	"github.com/galaxies-labs/galaxy/x/galaxy"
	"github.com/galaxies-labs/galaxy/x/galaxy/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.GalaxyKeeper(t)
	galaxy.InitGenesis(ctx, *k, genesisState)
	got := galaxy.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
