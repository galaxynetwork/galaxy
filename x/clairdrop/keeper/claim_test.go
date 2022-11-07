package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/galaxynetwork/galaxy/x/clairdrop/types"
)

func (suite *KeeperTestSuite) TestHookOfUnclaimableAccount() {

	require := suite.Require()

	pubKey1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pubKey1.Address())

	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, pubKey1, 0, 0))

	record, err := suite.app.ClairdropKeeper.GetClaimRecord(suite.ctx, addr1)
	require.NoError(err)

	require.Equal(
		record,
		types.ClaimRecord{},
	)

	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, addr1)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)

	require.Equal(
		sdk.NewCoins(),
		balances,
	)
}

func (suite *KeeperTestSuite) TestHookOfClaimableAccount() {

	require := suite.Require()

	pubKey1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pubKey1.Address())

	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, pubKey1, 0, 0))

	initalCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(1_000)))
	//4 actions
	claimRecords := []types.ClaimRecord{
		{
			Address:               (addr1).String(),
			InitalClaimableAmount: initalCoins,
			ActionCompleted:       []bool{false, false, false, false},
		},
	}

	err := suite.app.ClairdropKeeper.SetClaimRecords(suite.ctx, claimRecords)
	require.NoError(err)

	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, addr1)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)

	require.Equal(
		initalCoins.Sub(sdk.NewCoins(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(250*3)))),
		balances,
	)

	suite.app.ClairdropKeeper.AfterDelegationModified(suite.ctx, addr1, sdk.ValAddress(addr1))

	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)

	require.Equal(
		initalCoins.Sub(sdk.NewCoins(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(250*2)))),
		balances,
	)
}

func (suite *KeeperTestSuite) TestDuplicatedHook() {

	require := suite.Require()

	pubKey1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pubKey1.Address())

	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr1, pubKey1, 0, 0))

	initalCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(1_000)))
	//4 actions
	claimRecords := []types.ClaimRecord{
		{
			Address:               (addr1).String(),
			InitalClaimableAmount: initalCoins,
			ActionCompleted:       []bool{false, false, false, false},
		},
	}

	err := suite.app.ClairdropKeeper.SetClaimRecords(suite.ctx, claimRecords)
	require.NoError(err)

	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, addr1)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)

	require.Equal(
		initalCoins.Sub(sdk.NewCoins(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(250*3)))),
		balances,
	)

	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, addr1)

	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)

	require.Equal(
		initalCoins.Sub(sdk.NewCoins(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(250*3)))),
		balances,
	)
}

func (suite *KeeperTestSuite) TestClaimAction() {
	require := suite.Require()

	p1 := secp256k1.GenPrivKey().PubKey()
	acc1 := sdk.AccAddress(p1.Address())

	initalClaimableAmount := sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100))
	actionCompleted := []bool{false, false, false, false}

	suite.app.ClairdropKeeper.SetClaimRecords(suite.ctx, []types.ClaimRecord{
		{
			Address:               acc1.String(),
			InitalClaimableAmount: initalClaimableAmount,
			ActionCompleted:       actionCompleted,
		},
	})

	record, err := suite.app.ClairdropKeeper.GetClaimRecord(suite.ctx, acc1)
	require.NoError(err)
	require.Equal(record.ActionCompleted, actionCompleted)

	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, acc1)

	record, err = suite.app.ClairdropKeeper.GetClaimRecord(suite.ctx, acc1)
	require.NoError(err)

	//duplicated check
	suite.app.ClairdropKeeper.AfterProposalVote(suite.ctx, acc1)
	record, err = suite.app.ClairdropKeeper.GetClaimRecord(suite.ctx, acc1)
	require.NoError(err)

	actionCompleted2 := actionCompleted
	actionCompleted2[types.ClaimAction_value["Vote"]] = true
	require.Equal(record.ActionCompleted, actionCompleted2)

	suite.app.ClairdropKeeper.AfterDelegationModified(suite.ctx, acc1, sdk.ValAddress(acc1))

	record, err = suite.app.ClairdropKeeper.GetClaimRecord(suite.ctx, acc1)
	require.NoError(err)

	actionCompleted3 := actionCompleted2
	actionCompleted3[types.ClaimAction_value["Delegate"]] = true
	require.Equal(record.ActionCompleted, actionCompleted3)

}

//test action
