package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestPhaseInflation(t *testing.T) {
	minter := InitialMinter(sdk.NewDecWithPrec(50, 2))
	params := DefaultParams()

	tests := []struct {
		Phase     int
		Inflation sdk.Dec
	}{
		{
			Phase:     1,
			Inflation: sdk.NewDecWithPrec(50, 2),
		},
		{
			Phase:     2,
			Inflation: sdk.NewDecWithPrec(100, 2),
		}, {
			Phase:     3,
			Inflation: sdk.NewDecWithPrec(10, 2),
		}, {
			Phase:     4,
			Inflation: sdk.NewDecWithPrec(9, 2),
		}, {
			Phase:     5,
			Inflation: sdk.NewDecWithPrec(8, 2),
		}, {
			Phase:     6,
			Inflation: sdk.NewDecWithPrec(7, 2),
		}, {
			Phase:     7,
			Inflation: sdk.NewDecWithPrec(6, 2),
		}, {
			Phase:     8,
			Inflation: sdk.NewDecWithPrec(5, 2),
		}, {
			Phase:     9,
			Inflation: sdk.NewDecWithPrec(4, 2),
		}, {
			Phase:     10,
			Inflation: sdk.NewDecWithPrec(3, 2),
		}, {
			Phase:     11,
			Inflation: sdk.NewDecWithPrec(2, 2),
		}, {
			Phase:     12,
			Inflation: sdk.NewDecWithPrec(1, 2),
		},
	}

	for i, test := range tests {

		require.True(
			t,
			test.Inflation.Equal(
				minter.PhaseInflationRate(uint64(i+1), params),
			),
		)
	}
}

func TestAnnualProvisions(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()

	supply := sdk.NewInt(1_000_000_000)
	tests := []struct {
		Inflation sdk.Dec
		Supply    sdk.Int
	}{
		{
			Inflation: sdk.NewDecWithPrec(50, 2),
			Supply:    sdk.NewInt(500_000_000),
		},
		{
			Inflation: sdk.NewDecWithPrec(100, 2),
			Supply:    sdk.NewInt(1_000_000_000),
		},
		{
			Inflation: sdk.NewDecWithPrec(30, 2),
			Supply:    sdk.NewInt(300_000_000),
		},
		{
			Inflation: sdk.NewDecWithPrec(40, 2),
			Supply:    sdk.NewInt(400_000_000),
		},
		{
			Inflation: sdk.NewDecWithPrec(50, 2),
			Supply:    sdk.NewInt(500_000_000),
		},
		{
			Inflation: sdk.NewDecWithPrec(60, 2),
			Supply:    sdk.NewInt(600_000_000),
		},
	}
	for _, test := range tests {
		minter.Inflation = test.Inflation
		currentSupply := minter.NextAnnualProvisions(params, supply)
		t.Log(currentSupply)
		require.True(
			t,
			currentSupply.Equal(test.Supply.ToDec()),
		)
	}
}

func TestAnnualProvisionsEachBlock(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()
	genesisSupply := sdk.NewInt(1_000_000_000_000_000)
	supply := genesisSupply

	prev := sdk.NewDec(0)
	for i := 1; i <= int(params.BlocksPerYear)*int(params.StopInflationPhase); i = i + 100 {
		if minter.Inflation.Equal(sdk.ZeroDec()) {
			continue
		}
		currentBlock := uint64(i)

		currentPhase := minter.CurrentPhase(params, int64(currentBlock))
		if minter.Phase != uint64(currentPhase) {
			minter.Phase = uint64(minter.CurrentPhase(params, int64(currentBlock)))
			minter.Inflation = minter.PhaseInflationRate(uint64(minter.Phase), params)
			minter.AnnualProvisions = minter.NextAnnualProvisions(params, supply)
			require.NotEqual(
				t,
				minter.AnnualProvisions,
				prev,
			)
			t.Logf("changed : %s, current : %s", minter.AnnualProvisions.String(), prev.String())
		} else {
			require.Equal(
				t,
				minter.AnnualProvisions,
				prev,
			)
		}
		prev = minter.AnnualProvisions

		if minter.Inflation.Equal(sdk.ZeroDec()) {
			continue
		}

		mintedCoin := minter.BlockProvision(params)
		supply = supply.Add(sdk.NewInt(mintedCoin.Amount.Int64() * 100))
	}
}

func TestCurrentPhase(t *testing.T) {
	params := DefaultParams()
	minter := DefaultInitialMinter()
	tests := []struct {
		CurrentBlock int64
		Phase        int64
	}{
		{CurrentBlock: 1, Phase: 1},
		{CurrentBlock: 5, Phase: 1},
		{CurrentBlock: 100, Phase: 1},
		{CurrentBlock: int64(params.BlocksPerYear) - 1, Phase: 1},
		{CurrentBlock: int64(params.BlocksPerYear), Phase: 1},

		{CurrentBlock: int64(params.BlocksPerYear) + 1, Phase: 2},
		{CurrentBlock: int64(params.BlocksPerYear)*2 - 1, Phase: 2},
		{CurrentBlock: int64(params.BlocksPerYear) * 2, Phase: 2},
		{CurrentBlock: int64(params.BlocksPerYear)*3 - int64(params.BlocksPerYear), Phase: 2},

		{CurrentBlock: int64(params.BlocksPerYear)*3 - int64(params.BlocksPerYear) + 1, Phase: 3},
		{CurrentBlock: int64(params.BlocksPerYear)*2 + int64(params.BlocksPerYear) - 1, Phase: 3},
		{CurrentBlock: int64(params.BlocksPerYear) * 3, Phase: 3},

		{CurrentBlock: int64(params.BlocksPerYear)*3 + 1, Phase: 4},
		{CurrentBlock: int64(params.BlocksPerYear)*3 + int64(params.BlocksPerYear) - 1, Phase: 4},
		{CurrentBlock: int64(params.BlocksPerYear)*3 + int64(params.BlocksPerYear), Phase: 4},
		{CurrentBlock: int64(params.BlocksPerYear)*3 + int64(params.BlocksPerYear) + 1, Phase: 5},

		{CurrentBlock: int64(params.BlocksPerYear)*3 + 1, Phase: 4},
		{CurrentBlock: int64(params.BlocksPerYear)*3 + int64(params.BlocksPerYear), Phase: 4},
		{CurrentBlock: int64(params.BlocksPerYear) * 4, Phase: 4},

		{CurrentBlock: int64(params.BlocksPerYear)*10 - 1, Phase: 10},
		{CurrentBlock: int64(params.BlocksPerYear) * 10, Phase: 10},
		{CurrentBlock: int64(params.BlocksPerYear)*10 + 1, Phase: 11},
	}
	for _, test := range tests {
		currentPhase := minter.CurrentPhase(params, (test.CurrentBlock))

		t.Logf("current block : %d", test.CurrentBlock)
		require.True(
			t,
			currentPhase == (test.Phase),
		)

	}
}

func TestBlockProvision(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()

	secondsPerYear := int64(60 * 60 * 8766)

	tests := []struct {
		annualProvisions int64
		expProvisions    int64
	}{
		{secondsPerYear / 5, 1},
		{secondsPerYear/5 + 1, 1},
		{(secondsPerYear / 5) * 2, 2},
		{(secondsPerYear / 5) / 2, 0},
	}
	for _, tc := range tests {
		minter.AnnualProvisions = sdk.NewDec(tc.annualProvisions)
		provisions := minter.BlockProvision(params)

		expProvisions := sdk.NewCoin(params.MintDenom,
			sdk.NewInt(tc.expProvisions))

		require.True(t, expProvisions.IsEqual(provisions))
	}

}
