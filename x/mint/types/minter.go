package types

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMinter(inflation, annualProvisions sdk.Dec, phase uint64) Minter {
	return Minter{
		AnnualProvisions: annualProvisions,
		Phase:            phase,
		Inflation:        inflation,
	}
}

func InitialMinter(inflation sdk.Dec) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(0),
		uint64(0),
	)
}

func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDecWithPrec(50, 2),
	)
}

func ValidateMinter(minter Minter) error {
	if minter.Inflation.IsNegative() {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s",
			minter.Inflation.String())
	}
	return nil
}

func (m Minter) PhaseInflationRate(phase uint64, param Params) sdk.Dec {
	if param.ThresholdPhase >= phase {
		return sdk.NewDecWithPrec(50, 2).MulInt64(int64(phase))
	} else {
		return sdk.NewDecWithPrec(int64(param.StopInflationPhase-phase), 2)
	}
}

func (m Minter) CurrentPhase(params Params, currentBlock int64) int64 {
	v := int64(math.Ceil(float64(currentBlock) / float64(params.BlocksPerYear)))
	if v == 0 {
		return 1
	}
	return v
}

func (m Minter) NextAnnualProvisions(_ Params, totalSupply sdk.Int) sdk.Dec {
	return m.Inflation.MulInt(totalSupply)
}

func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(sdk.NewInt(int64(params.BlocksPerYear)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
