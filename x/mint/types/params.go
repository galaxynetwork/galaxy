package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyMintDenom                        = []byte("MintDenom")
	KeyThresholdPhase                   = []byte("ThresholdPhase")
	KeyStopInflationPhase               = []byte("StopInflationPhase")
	KeyDistributionProportions          = []byte("DistributionProportions")
	KeyWeightedDeveloperRewardsReceiver = []byte("WeightedDeveloperRewardsReceiver")
	KeyBlocksPerYear                    = []byte("BlocksPerYear")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string,
	thresholdPhase uint64,
	stopInflationPhase uint64,
	distrProportions DistributionProportions,
	weightedDevRewardsReceivers []DevloperWeightedAddress,
	blocksPerYear uint64,
) Params {
	return Params{
		MintDenom:                         mintDenom,
		ThresholdPhase:                    thresholdPhase,
		StopInflationPhase:                stopInflationPhase,
		DistributionProportions:           distrProportions,
		WeightedDeveloperRewardsReceivers: weightedDevRewardsReceivers,
		BlocksPerYear:                     blocksPerYear,
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultMintDenom,
		uint64(2),
		uint64(13),
		DistributionProportions{
			Staking:             sdk.NewDecWithPrec(2, 1),
			EcosystemIncentives: sdk.NewDecWithPrec(5, 1),
			DeveloperRewards:    sdk.NewDecWithPrec(2, 1),
			CommunityPool:       sdk.NewDecWithPrec(1, 1),
		},
		[]DevloperWeightedAddress{},
		uint64(60*60*8766/5),
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyThresholdPhase, &p.ThresholdPhase, validateThresholdPhase),
		paramtypes.NewParamSetPair(KeyStopInflationPhase, &p.StopInflationPhase, validateStopInflationPhase),
		paramtypes.NewParamSetPair(KeyDistributionProportions, &p.DistributionProportions, validateDistributionProportions),
		paramtypes.NewParamSetPair(KeyWeightedDeveloperRewardsReceiver, &p.WeightedDeveloperRewardsReceivers, validateWeightedDeveloperRewardsReceivers),
		paramtypes.NewParamSetPair(KeyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
	}
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateThresholdPhase(p.ThresholdPhase); err != nil {
		return err
	}
	if err := validateStopInflationPhase(p.StopInflationPhase); err != nil {
		return err
	}
	if err := validateDistributionProportions(p.DistributionProportions); err != nil {
		return err
	}
	if err := validateWeightedDeveloperRewardsReceivers(p.WeightedDeveloperRewardsReceivers); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}
	if p.ThresholdPhase >= p.StopInflationPhase {
		return fmt.Errorf("threshold phase must be smaller than stop inflation phase")
	}
	return nil
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}

	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateThresholdPhase(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateStopInflationPhase(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDistributionProportions(i interface{}) error {
	v, ok := i.(DistributionProportions)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.Staking.IsNegative() {
		return errors.New("staking distribution ratio should not be negative")
	}
	if v.EcosystemIncentives.IsNegative() {
		return errors.New("ecosystem incentives distribution ratio should not be negative")
	}
	if v.DeveloperRewards.IsNegative() {
		return errors.New("developer rewards distribution ratio should not be negative")
	}
	if v.CommunityPool.IsNegative() {
		return errors.New("community pool distribution ratio should not be negative")
	}

	totalProportions := v.Staking.Add(v.EcosystemIncentives).Add(v.DeveloperRewards).Add(v.CommunityPool)

	if !totalProportions.Equal(sdk.NewDec(1)) {
		return errors.New("total distributions ratio should be 1")
	}

	return nil
}

func validateWeightedDeveloperRewardsReceivers(i interface{}) error {
	v, ok := i.([]DevloperWeightedAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// fund community pool when rewards address is empty
	if len(v) == 0 {
		return nil
	}

	weightSum := sdk.NewDec(0)
	for i, w := range v {
		if w.Address != "" {
			_, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return fmt.Errorf("invalid address at %dth", i)
			}
		}
		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(sdk.NewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}
