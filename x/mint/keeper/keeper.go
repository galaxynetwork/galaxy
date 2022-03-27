package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/galaxies-labs/galaxy/x/mint/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramStore paramtypes.Subspace

	ak types.AccountKeeper
	bk types.BankKeeper
	dk types.DistributionKeeper

	feeCollectorName string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	dk types.DistributionKeeper,
	feeCollectorName string,

) Keeper {

	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		paramStore:       ps,
		ak:               ak,
		bk:               bk,
		dk:               dk,
		feeCollectorName: feeCollectorName,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetDeveloperAddress(ctx sdk.Context) []string {
	params := k.GetParams(ctx)

	address := []string{}

	for _, d := range params.WeightedDeveloperRewardsReceivers {
		address = append(address, d.Address)
	}
	return address
}

func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		return nil
	}
	return k.bk.MintCoins(ctx, types.ModuleName, newCoins)
}

func (k Keeper) ModuleBalance(ctx sdk.Context) sdk.Coin {
	params := k.GetParams(ctx)
	return k.bk.GetBalance(
		ctx,
		k.ak.GetModuleAddress(types.ModuleName),
		params.MintDenom,
	)
}

func (k Keeper) TokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bk.GetSupply(ctx, denom).Amount
}

func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.bk.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}

func (k Keeper) FundToCommuinityPool(ctx sdk.Context, coins sdk.Coins) error {
	err := k.dk.FundCommunityPool(
		ctx,
		coins,
		k.ak.GetModuleAddress(types.ModuleName),
	)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetProportions(ctx sdk.Context, mintedCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(mintedCoin.Denom, mintedCoin.Amount.ToDec().Mul(ratio).TruncateInt())
}

func (k Keeper) DistributeMintedCoin(ctx sdk.Context, mintedCoin sdk.Coin) error {
	params := k.GetParams(ctx)
	proportions := params.DistributionProportions

	stakingIncentives := sdk.NewCoins(k.GetProportions(ctx, mintedCoin, proportions.Staking))
	err := k.bk.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, stakingIncentives)
	if err != nil {
		return err
	}

	//fund to community pool before develop
	ecosystemIncentivesCoins := sdk.NewCoins(k.GetProportions(ctx, mintedCoin, proportions.EcosystemIncentives))
	err = k.dk.FundCommunityPool(ctx, ecosystemIncentivesCoins, k.ak.GetModuleAddress(types.ModuleName))

	if err != nil {
		return err
	}

	developerRewards := k.GetProportions(ctx, mintedCoin, proportions.DeveloperRewards)
	developerRewardsCoins := sdk.NewCoins(developerRewards)
	developerReceivers := params.WeightedDeveloperRewardsReceivers
	//if dev is none fund to communiy
	if len(developerReceivers) == 0 {
		developerRewardsCoins := sdk.NewCoins(developerRewards)
		err = k.dk.FundCommunityPool(ctx, developerRewardsCoins, k.ak.GetModuleAddress(types.ModuleName))
		if err != nil {
			return err
		}
	} else {
		for _, weightedReceiver := range developerReceivers {
			devProportions := k.GetProportions(ctx, developerRewards, weightedReceiver.Weight)
			developerRewardProPortion := sdk.NewCoins(devProportions)
			if weightedReceiver.Address == "" {
				err = k.dk.FundCommunityPool(ctx, developerRewardProPortion,
					k.ak.GetModuleAddress(types.ModuleName))
				if err != nil {
					return err
				}
			} else {
				address, err := sdk.AccAddressFromBech32(weightedReceiver.Address)
				if err != nil {
					return err
				}
				err = k.bk.SendCoinsFromModuleToAccount(
					ctx, types.ModuleName, address, developerRewardProPortion)
				if err != nil {
					return err
				}
			}
		}
	}

	communityPool := sdk.NewCoins(mintedCoin).Sub(stakingIncentives).Sub(ecosystemIncentivesCoins).Sub(developerRewardsCoins)
	err = k.dk.FundCommunityPool(ctx, communityPool, k.ak.GetModuleAddress(types.ModuleName))

	if err != nil {
		return err
	}

	return nil
}
