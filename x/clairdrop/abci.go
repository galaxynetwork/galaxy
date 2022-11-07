package clairdrop

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxynetwork/galaxy/x/clairdrop/keeper"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

	params := k.GetParams(ctx)

	if ctx.BlockTime().After(params.ClairdropEndTime) && k.GetModuleAccountBalance(ctx).IsPositive() {
		err := k.EndAirdrop(ctx)
		if err != nil {
			panic(err)
		}
	}
}
