package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyClairdropStartTime = []byte("ClairdropStartTime")
	KeyClairdropEndTime   = []byte("ClairdropEndTime")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	clairdropStartTime time.Time,
	clairdropEndTime time.Time,
) Params {
	return Params{
		ClairdropStartTime: clairdropStartTime,
		ClairdropEndTime:   clairdropEndTime,
	}
}

func DefaultParams() Params {
	return NewParams(
		time.Time{},
		time.Time{}.Add(time.Hour*24*150),
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyClairdropStartTime, &p.ClairdropStartTime, validateClairdropTime),
		paramtypes.NewParamSetPair(KeyClairdropEndTime, &p.ClairdropEndTime, validateClairdropTime),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateClairdropTime(p.ClairdropStartTime); err != nil {
		return err
	}
	if err := validateClairdropTime(p.ClairdropEndTime); err != nil {
		return err
	}
	if p.ClairdropEndTime.Before(p.ClairdropStartTime) {
		return fmt.Errorf("clairdrop end time must be late than clairdrop start time")
	}
	return nil
}

func validateClairdropTime(i interface{}) error {
	_, ok := i.(time.Time)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
