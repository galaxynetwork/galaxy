package types

import (
	"bytes"
	fmt "fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// reBrandIDString can be 1 ~ 30 characters long and support letters
	reBrandIDString = `[a-zA-Z0-9]{1,30}`

	// Constants pertaining to a Description object
	MaxNameLength    int = 100
	MaxDetailsLength int = 10000
)

var (
	reBrandID = regexp.MustCompile(fmt.Sprintf(`^%s$`, reBrandIDString))
)

// NewBrandDescription returns a brand.
func NewBrand(id string, owner sdk.AccAddress, description BrandDescription) Brand {
	brandAddress := NewBrandAddress(id)
	return Brand{
		Id:           id,
		Owner:        owner.String(),
		BrandAddress: brandAddress.String(),
		Description:  description,
	}
}

// Validate defines a method basic validating and trim all spaces in id.
func (brand *Brand) Validate() error {
	brand.Id = strings.TrimSpace(brand.Id)

	if !reBrandID.MatchString(brand.Id) {
		return sdkerrors.Wrapf(ErrInvalidBrandID, "invalid brand id: %s", brand.Id)
	}

	brandAcc, err := sdk.AccAddressFromBech32(brand.BrandAddress)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidBrandAddress, "invalid brand address (%s)", err)
	}

	if bytes.Compare(brandAcc, NewBrandAddress(brand.Id)) != 0 {
		return sdkerrors.Wrap(ErrInvalidBrandAddress, "brand address is not generated from a id")
	}

	_, err = sdk.AccAddressFromBech32(brand.Owner)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidBrandOwnerAddress, "invalid brand owner address (%s)", err)
	}

	if err = brand.Description.Validate(); err != nil {
		return err
	}

	return nil
}

type Brands = []Brand

// NewBrandDescription returns a brand description.
func NewBrandDescription(name, details, brandImageUri string) BrandDescription {
	return BrandDescription{
		Name:          name,
		Details:       details,
		BrandImageUri: brandImageUri,
	}
}

// Validate defines a method trim all space and basic validation.
func (bd *BrandDescription) Validate() error {
	bd.Name = strings.TrimSpace(bd.Name)
	bd.Details = strings.TrimSpace(bd.Details)
	bd.BrandImageUri = strings.TrimSpace(bd.BrandImageUri)

	if len(bd.Name) == 0 {
		return sdkerrors.Wrap(ErrInvalidBrandName, "brand name cannot be blank")
	}

	if len(bd.Name) > MaxNameLength {
		return sdkerrors.Wrapf(ErrInvalidBrandName, "brand name is longer than max length of %d", MaxNameLength)
	}

	if len(bd.Details) > MaxDetailsLength {
		return sdkerrors.Wrapf(ErrInvalidBrandDetails, "brand details is longer than max length of %d", MaxDetailsLength)
	}

	return nil
}
