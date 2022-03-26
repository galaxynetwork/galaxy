package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/galaxies-labs/galaxy/app"
	"github.com/galaxies-labs/galaxy/x/mint/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.App
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC(), Height: 1})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestMintCoins() {
	mintKeeper := suite.app.MintKeeper
	require := suite.Require()

	mintCoin := sdk.NewInt(100000000000)
	err := mintKeeper.MintCoins(suite.ctx, sdk.NewCoins(sdk.NewCoin("uglx", mintCoin)))
	require.NoError(err)

	amount := mintKeeper.TokenSupply(suite.ctx, "uglx")

	require.Equal(amount, mintCoin)
}

func (suite *KeeperTestSuite) TestMintCoinsToFeeCollectorAndGetProportions() {
	mintKeeper := suite.app.MintKeeper
	require := suite.Require()

	fee := sdk.NewCoin("uglx", sdk.NewInt(0))
	coin := mintKeeper.GetProportions(suite.ctx, fee, sdk.NewDecWithPrec(2, 1))
	require.Equal("0uglx", coin.String())

	fee = sdk.NewCoin("uglx", sdk.NewInt(100000))
	fees := sdk.NewCoins(fee)

	err := simapp.FundModuleAccount(
		suite.app.BankKeeper,
		suite.ctx,
		authtypes.FeeCollectorName,
		fees,
	)
	require.NoError(err)

	coin = mintKeeper.GetProportions(suite.ctx, fee, sdk.NewDecWithPrec(2, 1))
	require.Equal(fees[0].Amount.Quo(sdk.NewInt(5)), coin.Amount)
}

func (suite *KeeperTestSuite) TestGetProportions() {
	mintKeeper := suite.app.MintKeeper
	require := suite.Require()
	params := mintKeeper.GetParams(suite.ctx)

	mintCoin := sdk.NewCoin("uglx", sdk.NewInt(100000))

	require.Equal(
		mintKeeper.GetProportions(suite.ctx, mintCoin, params.DistributionProportions.DeveloperRewards).String(),
		"20000uglx",
	)
	require.Equal(
		mintKeeper.GetProportions(suite.ctx, mintCoin, params.DistributionProportions.CommunityPool).String(),
		"10000uglx",
	)
	require.Equal(
		mintKeeper.GetProportions(suite.ctx, mintCoin, params.DistributionProportions.Staking).String(),
		"20000uglx",
	)
	require.Equal(
		mintKeeper.GetProportions(suite.ctx, mintCoin, params.DistributionProportions.EcosystemIncentives).String(),
		"50000uglx",
	)
}

func (suite *KeeperTestSuite) TestDistrAssetToCommunityPoolBeforeDevelopModuleWhenDevEmpty() {
	mintKeeper := suite.app.MintKeeper
	distrKeeper := suite.app.DistrKeeper
	bankKeeper := suite.app.BankKeeper
	authKeeper := suite.app.AccountKeeper
	require := suite.Require()
	params := mintKeeper.GetParams(suite.ctx)

	mintCoin := sdk.NewCoin("uglx", sdk.NewInt(100000))

	suite.T().Log(mintCoin.String())

	mintKeeper.MintCoins(suite.ctx, sdk.NewCoins(mintCoin))

	err := mintKeeper.DistributeMintedCoin(suite.ctx, mintCoin)

	require.NoError(err)

	communityCoins := distrKeeper.GetFeePoolCommunityCoins(suite.ctx)

	require.Equal(
		bankKeeper.GetBalance(suite.ctx, authKeeper.GetModuleAddress(authtypes.FeeCollectorName), "uglx").Amount.ToDec().String(),
		mintCoin.Amount.ToDec().Mul(params.DistributionProportions.Staking).String(),
	)

	require.Equal(
		communityCoins.AmountOf("uglx").String(),
		mintCoin.Amount.ToDec().Mul(params.DistributionProportions.EcosystemIncentives).Add(
			mintCoin.Amount.ToDec().Mul(params.DistributionProportions.CommunityPool),
		).Add(
			mintCoin.Amount.ToDec().Mul(params.DistributionProportions.DeveloperRewards),
		).String(),
	)
}

func (suite *KeeperTestSuite) TestDistrAssetToCommunityPoolBeforeDevelopModuleWhenDevNotEmpty() {
	mintKeeper := suite.app.MintKeeper
	distrKeeper := suite.app.DistrKeeper
	bankKeeper := suite.app.BankKeeper
	authKeeper := suite.app.AccountKeeper
	require := suite.Require()

	params := mintKeeper.GetParams(suite.ctx)
	params.WeightedDeveloperRewardsReceivers = []types.DevloperWeightedAddress{
		{
			Address: sdk.AccAddress([]byte("addr2---")).String(),
			Weight:  sdk.NewDecWithPrec(5, 1),
		}, {
			Address: sdk.AccAddress([]byte("addr2---")).String(),
			Weight:  sdk.NewDecWithPrec(5, 1),
		},
	}

	err := params.Validate()
	require.NoError(err)

	mintKeeper.SetParams(suite.ctx, params)

	mintCoin := sdk.NewCoin("uglx", sdk.NewInt(100000))

	mintKeeper.MintCoins(suite.ctx, sdk.NewCoins(mintCoin))

	err = mintKeeper.DistributeMintedCoin(suite.ctx, mintCoin)

	require.NoError(err)

	communityCoins := distrKeeper.GetFeePoolCommunityCoins(suite.ctx)

	require.Equal(
		bankKeeper.GetBalance(suite.ctx, authKeeper.GetModuleAddress(authtypes.FeeCollectorName), "uglx").Amount.ToDec().String(),
		mintCoin.Amount.ToDec().Mul(params.DistributionProportions.Staking).String(),
	)

	require.Equal(
		communityCoins.AmountOf("uglx").String(),
		mintCoin.Amount.ToDec().Mul(params.DistributionProportions.EcosystemIncentives).Add(
			mintCoin.Amount.ToDec().Mul(params.DistributionProportions.CommunityPool),
		).String(),
	)
}

func (suite *KeeperTestSuite) TestMintInflationThenDistr() {
	require := suite.Require()
	mintKeeper := suite.app.MintKeeper
	bankKeeper := suite.app.BankKeeper
	authKeeper := suite.app.AccountKeeper
	params := mintKeeper.GetParams(suite.ctx)

	var devRemaningCoin sdk.Coin

	params.WeightedDeveloperRewardsReceivers = []types.DevloperWeightedAddress{
		{
			Address: sdk.AccAddress([]byte("addr1")).String(),
			Weight:  sdk.NewDecWithPrec(4, 1),
		}, {
			Address: sdk.AccAddress([]byte("addr2")).String(),
			Weight:  sdk.NewDecWithPrec(3, 1),
		}, {
			Address: sdk.AccAddress([]byte("addr3")).String(),
			Weight:  sdk.NewDecWithPrec(2, 1),
		}, {
			Address: sdk.AccAddress([]byte("addr4")).String(),
			Weight:  sdk.NewDecWithPrec(1, 1),
		},
	}

	mintKeeper.SetParams(suite.ctx, params)
	genesisSupply := sdk.NewInt(1_000_000_000_000_000)
	mintKeeper.MintCoins(suite.ctx, sdk.NewCoins(sdk.NewCoin(params.MintDenom, genesisSupply)))
	suite.T().Log("genesis total supply : ", mintKeeper.TokenSupply(suite.ctx, params.MintDenom).String())

	tests := []struct {
		Phase uint64
	}{}

	for i := uint64(0); i <= params.StopInflationPhase; i++ {
		tests = append(tests, struct {
			Phase uint64
		}{Phase: i})
	}

	for _, test := range tests {
		phase := test.Phase
		suite.T().Log(phase)

		defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

		minter := mintKeeper.GetMinter(suite.ctx)

		newInflation := minter.PhaseInflationRate(phase, params)
		totalSupply := mintKeeper.TokenSupply(suite.ctx, params.MintDenom)
		minter.Inflation = newInflation
		minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalSupply)

		// infalation end
		if params.StopInflationPhase <= phase || phase == 0 {
			require.Equal(
				minter.Inflation,
				sdk.ZeroDec(),
			)
			require.Equal(
				minter.AnnualProvisions,
				sdk.ZeroDec(),
			)
		}

		suite.T().Log("minter inflation ", minter.Inflation)
		suite.T().Log("minter annual provisions ", minter.AnnualProvisions)
		mintKeeper.SetMinter(suite.ctx, minter)

		if minter.Inflation.Equal(sdk.ZeroDec()) {
			suite.T().Log("return because inflation zero")
			if minter.Inflation.Equal(sdk.ZeroDec()) {
				//if still has ramaning amount  it will be fund to community pool
				devRemaningCoin = mintKeeper.ModuleBalance(suite.ctx)
				if devRemaningCoin.IsPositive() {
					suite.T().Log("remaning fund to community pool")
					err := mintKeeper.FundToCommuinityPool(suite.ctx, sdk.NewCoins(devRemaningCoin))
					require.NoError(err)
				}
			}
			continue
		}

		blockProvision := minter.BlockProvision(params)
		mintCoin := sdk.NewCoin(blockProvision.Denom, blockProvision.Amount.Mul(sdk.NewIntFromUint64(params.BlocksPerYear)))
		mintCoins := sdk.NewCoins(mintCoin)
		suite.T().Log("minted coin", mintCoin.String())

		err := mintKeeper.MintCoins(suite.ctx, mintCoins)
		require.NoError(err)

		err = mintKeeper.DistributeMintedCoin(suite.ctx, mintCoin)
		require.NoError(err)

		if mintCoin.Amount.IsInt64() {
			defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintCoin.Amount.Int64()), "minted_tokens")
		}

		suite.ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMint,
				sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
				sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, mintCoin.Amount.String()),
			),
		)
	}

	require.Equal(
		bankKeeper.GetBalance(
			suite.ctx,
			authKeeper.GetModuleAddress(
				types.ModuleName,
			),
			params.MintDenom,
		).Amount.ToDec(),
		sdk.ZeroDec(),
	)

	totalSupply := mintKeeper.TokenSupply(suite.ctx, params.MintDenom)
	feeCollectorAmount := bankKeeper.GetBalance(suite.ctx, authKeeper.GetModuleAddress(authtypes.FeeCollectorName), params.MintDenom)
	distributionAmount := bankKeeper.GetBalance(suite.ctx, authKeeper.GetModuleAddress(distrtypes.ModuleName), params.MintDenom)

	require.Equal(
		totalSupply.ToDec(),
		feeCollectorAmount.Add(
			distributionAmount,
		).Add(
			mintKeeper.GetProportions(
				suite.ctx,
				sdk.NewCoin(params.MintDenom, totalSupply.Sub(genesisSupply)),
				params.DistributionProportions.DeveloperRewards,
			),
		).Sub(
			devRemaningCoin,
		).Amount.ToDec(),
	)

}
