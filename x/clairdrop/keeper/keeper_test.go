package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/app"
	"github.com/galaxies-labs/galaxy/x/clairdrop/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.App
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = app.Setup(false)
	suite.ctx = suite.app.GetBaseApp().NewContext(false, tmproto.Header{Height: 1, ChainID: "galaxy-1", Time: time.Now().UTC()})

	airdropStartTime := time.Now()
	airdropEndTime := time.Now().Add(time.Hour * 2)

	suite.app.ClairdropKeeper.CreateModuleAccount(
		suite.ctx,
		sdk.NewCoin("uglx", sdk.NewInt(1_000_000)),
	)

	suite.app.ClairdropKeeper.SetParams(
		suite.ctx,
		types.Params{
			ClairdropStartTime: airdropStartTime,
			ClairdropEndTime:   airdropEndTime,
		},
	)

	suite.ctx = suite.ctx.WithBlockTime(airdropStartTime)

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
