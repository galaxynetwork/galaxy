package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/app"
	"github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	app         *app.App
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) TestCreateBrand() {
}

func (suite *KeeperTestSuite) TestEditBrand() {
}

func (suite *KeeperTestSuite) TestTransferOwnershipBrand() {
}
