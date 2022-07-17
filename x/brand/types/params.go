package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyBrandCreationFee = []byte("BrandCreationFee")
)

// ParamKeyTable ParamTable for brand module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams returns a brand params.
func NewParams(brandCreationFee sdk.Coin) Params {
	return Params{
		BrandCreationFee: brandCreationFee,
	}
}

// DefaultParams returns a brand default params.
func DefaultParams() Params {
	return NewParams(sdk.NewCoin(DefaultBrandCreationFeeDenom, sdk.NewInt(100_000_000)))
}

//
func (params Params) Validate() error {
	if err := validateBrandCreationFee(params.BrandCreationFee); err != nil {
		return err
	}
	return nil
}

// ParamSetPairs Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBrandCreationFee, &p.BrandCreationFee, validateBrandCreationFee),
	}
}

// validateBrandCreationFee defines a method validating sdk coin type
func validateBrandCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Validate() != nil {
		return fmt.Errorf("invalid brand creation fee: %+v", i)
	}

	if !strings.EqualFold(v.Denom, DefaultBrandCreationFeeDenom) {
		return fmt.Errorf("invalid brand creation fee denom: %s, it must be: %s", v.Denom, DefaultBrandCreationFeeDenom)
	}

	return nil
}
