package mint

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxynetwork/galaxy/x/mint/keeper"
	"github.com/galaxynetwork/galaxy/x/mint/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	minter := k.GetMinter(ctx)

	// inflation end
	if minter.Inflation.Equal(sdk.ZeroDec()) {
		return
	}

	params := k.GetParams(ctx)
	totalSupply := k.TokenSupply(ctx, params.MintDenom)
	currentBlock := uint64(ctx.BlockHeight())

	currentPhase := minter.CurrentPhase(params, int64(currentBlock))
	if minter.Phase != uint64(currentPhase) {
		minter.Phase = uint64(minter.CurrentPhase(params, int64(currentBlock)))
		minter.Inflation = minter.PhaseInflationRate(uint64(minter.Phase), params)
		minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalSupply)
	}

	k.SetMinter(ctx, minter)

	// inflation end
	if minter.Inflation.Equal(sdk.ZeroDec()) {
		//if still has ramaning amount  it will be fund to community pool
		coin := k.ModuleBalance(ctx)
		if coin.IsPositive() {
			k.FundToCommuinityPool(ctx, sdk.NewCoins(coin))
		}
		return
	}

	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)

	if err != nil {
		panic(err)
	}

	err = k.DistributeMintedCoin(ctx, mintedCoin)

	if err != nil {
		panic(err)
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)

}
