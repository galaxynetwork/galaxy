package keeper_test

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/galaxies-labs/galaxy/x/clairdrop/types"
	minttypes "github.com/galaxies-labs/galaxy/x/mint/types"
)

func (suite *KeeperTestSuite) TestEndAirdropCall() {
	require := suite.Require()

	params := suite.app.ClairdropKeeper.GetParams(suite.ctx)

	tests := []struct {
		time   time.Time
		expect bool
	}{
		{time: params.ClairdropStartTime, expect: false},
		{time: params.ClairdropEndTime, expect: false},
		{time: params.ClairdropEndTime.Add(time.Hour), expect: true},
		{time: params.ClairdropEndTime.Add(time.Hour * 2), expect: false},
	}

	suite.T().Log(
		suite.app.ClairdropKeeper.GetModuleAccountBalance(suite.ctx).String(),
	)

	for _, test := range tests {
		if test.time.After(params.ClairdropEndTime) && suite.app.ClairdropKeeper.GetModuleAccountBalance(suite.ctx).IsPositive() {
			require.True(test.expect)
			err := suite.app.ClairdropKeeper.EndAirdrop(suite.ctx)
			require.NoError(err)
		} else {
			require.True(!test.expect)
		}
	}
}

func (suite *KeeperTestSuite) TestClawbackAirdrop() {
	require := suite.Require()

	p1 := secp256k1.GenPrivKey().PubKey()
	p2 := secp256k1.GenPrivKey().PubKey()
	p3 := secp256k1.GenPrivKey().PubKey()
	p4 := secp256k1.GenPrivKey().PubKey()
	p5 := secp256k1.GenPrivKey().PubKey()
	p6 := secp256k1.GenPrivKey().PubKey()
	p7 := secp256k1.GenPrivKey().PubKey()
	p8 := secp256k1.GenPrivKey().PubKey()
	p9 := secp256k1.GenPrivKey().PubKey()
	p10 := secp256k1.GenPrivKey().PubKey()

	tests := []struct {
		name    string
		expect  sdk.Coin
		address string
	}{{
		name:    "set record | active | hooks | genesis",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 125),
		address: sdk.AccAddress(p1.Address()).String(),
	}, {
		name:    "set record | active | no hooks | genesis",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		address: sdk.AccAddress(p2.Address()).String(),
	}, {
		name:    "set record | inactive | hooks | genesis",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		address: sdk.AccAddress(p3.Address()).String(),
	}, {
		name:    "set record | inactive | no hooks | genesis",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		address: sdk.AccAddress(p4.Address()).String(),
	}, {
		name:    "no record | genesis",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		address: sdk.AccAddress(p5.Address()).String(),
	}, {
		name:    "no record",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		address: sdk.AccAddress(p6.Address()).String(),
	}, {
		name:    "dev | active | hooks",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 125),
		address: sdk.AccAddress(p7.Address()).String(),
	}, {
		name:    "dev | no active | hooks",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 125),
		address: sdk.AccAddress(p8.Address()).String(),
	}, {
		name:    "dev | active | no hooks",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		address: sdk.AccAddress(p9.Address()).String(),
	}, {
		name:    "dev | no active | no hooks",
		expect:  sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		address: sdk.AccAddress(p10.Address()).String(),
	}}

	//set dev
	mintParams := suite.app.MintKeeper.GetParams(suite.ctx)
	for i, test := range tests {
		if i >= 6 {
			mintParams.WeightedDeveloperRewardsReceivers = append(mintParams.WeightedDeveloperRewardsReceivers, minttypes.DevloperWeightedAddress{
				Address: test.address,
				Weight:  sdk.NewDecWithPrec(25, 2),
			})

		}
	}
	suite.app.MintKeeper.SetParams(suite.ctx, mintParams)

	claimRecords := []types.ClaimRecord{}

	for i, test := range tests {
		if i == 4 || i == 5 {
			continue
		}
		suite.T().Log(
			test.address,
		)
		claimRecords = append(claimRecords, types.ClaimRecord{
			Address:               test.address,
			InitalClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100)),
			ActionCompleted:       []bool{false, false, false, false},
		})
	}

	suite.app.ClairdropKeeper.SetClaimRecords(suite.ctx, claimRecords)

	for i, test := range tests {
		if i == 5 {
			continue
		}
		acc, err := sdk.AccAddressFromBech32(test.address)
		require.NoError(err)
		simapp.FundAccount(
			suite.app.BankKeeper,
			suite.ctx,
			acc,
			sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100)),
		)
	}

	//hooks
	acc, err := sdk.AccAddressFromBech32(tests[0].address)
	require.NoError(err)
	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, acc)

	acc3, err3 := sdk.AccAddressFromBech32(tests[2].address)
	require.NoError(err3)
	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, acc3)

	//dev hooks
	acc6, err6 := sdk.AccAddressFromBech32(tests[6].address)
	require.NoError(err6)
	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, acc6)

	acc7, err7 := sdk.AccAddressFromBech32(tests[7].address)
	require.NoError(err7)
	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, acc7)

	//active
	acc2, err2 := sdk.AccAddressFromBech32(tests[1].address)
	require.NoError(err2)

	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(
		acc,
		p1, 0, 1,
	))

	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(
		acc2,
		p2, 0, 1,
	))

	//dev actinve
	acc8, err8 := sdk.AccAddressFromBech32(tests[8].address)
	require.NoError(err8)

	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(
		acc6,
		p7, 0, 1,
	))

	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(
		acc8,
		p9, 0, 1,
	))

	err = suite.app.ClairdropKeeper.EndAirdrop(suite.ctx)
	require.NoError(err)

	for _, test := range tests {
		acc, err := sdk.AccAddressFromBech32(test.address)
		require.NoError(err)

		suite.T().Log(test.name, test.expect,
			suite.app.BankKeeper.GetBalance(suite.ctx, acc, types.DefaultClaimDenom))
		require.Equal(
			test.expect,
			suite.app.BankKeeper.GetBalance(suite.ctx, acc, types.DefaultClaimDenom),
		)
	}
}
