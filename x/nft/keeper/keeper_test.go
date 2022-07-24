package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/app"
	"github.com/galaxies-labs/galaxy/x/nft/keeper"
	"github.com/galaxies-labs/galaxy/x/nft/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.App
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = app.Setup(false)
	suite.ctx = suite.app.GetBaseApp().NewContext(false, tmproto.Header{Height: 1, ChainID: "galaxy-1", Time: time.Now().UTC()})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(suite.app.NFTKeeper))
	suite.queryClient = types.NewQueryClient(queryHelper)

	suite.msgServer = keeper.NewMsgServerImpl(suite.app.NFTKeeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
