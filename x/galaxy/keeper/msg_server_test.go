package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/galaxies-labs/galaxy/testutil/keeper"
	"github.com/galaxies-labs/galaxy/x/galaxy/keeper"
	"github.com/galaxies-labs/galaxy/x/galaxy/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.GalaxyKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
