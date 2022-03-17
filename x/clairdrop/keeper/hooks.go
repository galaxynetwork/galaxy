package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/galaxies-labs/galaxy/x/clairdrop/types"
)

func (k Keeper) AfterProposalVote(ctx sdk.Context, voterAddr sdk.AccAddress) {
	_, err := k.ClaimForAction(ctx, voterAddr, types.Vote)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	_, err := k.ClaimForAction(ctx, delAddr, types.Delegate)
	if err != nil {
		panic(err)
	}
}

type Hooks struct {
	k Keeper
}

var _ govtypes.GovHooks = Hooks{}
var _ stakingtypes.StakingHooks = Hooks{}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

//staking hooks
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.k.AfterDelegationModified(ctx, delAddr, valAddr)
}

func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
}

//vote hooks
func (h Hooks) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	h.k.AfterProposalVote(ctx, voterAddr)
}
func (h Hooks) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {
}

func (h Hooks) AfterProposalDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) {
}

func (h Hooks) AfterProposalFailedMinDeposit(ctx sdk.Context, proposalID uint64) {
}
func (h Hooks) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {
}
