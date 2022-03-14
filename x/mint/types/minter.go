package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMinter(inflation, annualProvisions sdk.Dec, phase uint64, startPhaseBlock uint64) Minter {
	return Minter{
		AnnualProvisions: annualProvisions,
		Phase:            phase,
		Inflation:        inflation,
		StartPhaseBlock:  startPhaseBlock,
	}
}

func InitialMinter(inflation sdk.Dec) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(0),
		uint64(0),
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

func (m Minter) NextPhase(params Params, currentBlock uint64) uint64 {
	nonePhase := m.Phase == 0
	if nonePhase {
		return 1
	}
	blockNewPhase := m.StartPhaseBlock + params.BlocksPerYear
	if blockNewPhase > currentBlock {
		return m.Phase
	}
	return m.Phase + 1
}

func (m Minter) NextAnnualProvisions(_ Params, totalSupply sdk.Int) sdk.Dec {
	return m.Inflation.MulInt(totalSupply)
}

func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(sdk.NewInt(int64(params.BlocksPerYear)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
