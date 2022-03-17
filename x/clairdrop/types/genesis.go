package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		ModuleAccountBalance: sdk.NewCoin(DefaultClaimDenom, sdk.ZeroInt()),
		Params:               DefaultParams(),
		ClaimRecords:         []ClaimRecord{},
	}
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	if data.ModuleAccountBalance.Denom != DefaultClaimDenom {
		return fmt.Errorf("denom for module and claim does not match")
	}

	totalClaimable := sdk.Coins{}

	for index, claimRecord := range data.ClaimRecords {
		if claimRecord.InitalClaimableAmount.GetDenomByIndex(0) != DefaultClaimDenom {
			return fmt.Errorf("denom for module and claim records does not match index : %d", index)
		}
		totalClaimable = totalClaimable.Add(claimRecord.InitalClaimableAmount...)
	}

	if !totalClaimable.IsEqual(sdk.NewCoins(data.ModuleAccountBalance)) {
		return fmt.Errorf("claim module account balance != sum of all claim record InitialClaimableAmounts")
	}

	return nil
}
