package types

import (
	"fmt"
	"regexp"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	brandtypes "github.com/galaxies-labs/galaxy/x/brand/types"
)

const (
	DoNotModifyDesc = "[do-not-modify]"

	reClassIDString = `[a-zA-Z0-9][a-zA-Z0-9-]{2,50}`

	MaxFeeBasisPoints     uint32 = 10_000
	MaxClassNameLength    int    = 100
	MaxClassDetailsLength int    = 1024
	MaxUriLength          int    = 2048
)

var (
	reClassID = regexp.MustCompile(fmt.Sprintf(`^%s$`, reClassIDString))
)

func NewClass(brandID, id string, feeBasisPoints uint32, description ClassDescription) Class {
	return Class{
		BrandId:        brandID,
		Id:             id,
		FeeBasisPoints: feeBasisPoints,
		Description:    description,
	}
}

func (class *Class) Validate() error {
	if err := brandtypes.ValidateBrandID(class.BrandId); err != nil {
		return err
	}

	if err := ValidateClassId(class.Id); err != nil {
		return err
	}

	if err := ValidateFeeBasisPoints(class.FeeBasisPoints); err != nil {
		return err
	}

	if err := class.Description.Validate(); err != nil {
		return err
	}

	return nil
}

func NewClassDescription(name, details, externalUrl, imageUri string) ClassDescription {
	return ClassDescription{
		Name:        name,
		Details:     details,
		ExternalUrl: externalUrl,
		ImageUri:    imageUri,
	}
}

func (desc *ClassDescription) UpdateDescription(desc2 ClassDescription) ClassDescription {
	if desc2.Name == DoNotModifyDesc {
		desc2.Name = desc.Name
	}
	if desc2.Details == DoNotModifyDesc {
		desc2.Details = desc.Details
	}
	if desc2.ExternalUrl == DoNotModifyDesc {
		desc2.ExternalUrl = desc.ExternalUrl
	}
	if desc2.ImageUri == DoNotModifyDesc {
		desc2.ImageUri = desc.ImageUri
	}

	return desc2
}

func (desc *ClassDescription) TrimSpace() ClassDescription {
	desc.Name = strings.TrimSpace(desc.Name)
	desc.Details = strings.TrimSpace(desc.Details)
	desc.ImageUri = strings.TrimSpace(desc.ImageUri)
	desc.ExternalUrl = strings.TrimSpace(desc.ExternalUrl)

	return NewClassDescription(desc.Name, desc.Details, desc.ExternalUrl, desc.ImageUri)
}

func (desc *ClassDescription) Validate() error {
	if len(desc.Name) > MaxClassNameLength {
		return sdkerrors.Wrapf(ErrInvalidClassDescription, "invalid name length; got: %d, max: %d", len(desc.Name), MaxClassNameLength)
	}

	if len(desc.Details) > MaxClassDetailsLength {
		return sdkerrors.Wrapf(ErrInvalidClassDescription, "invalid details length; got: %d, max: %d", len(desc.Details), MaxClassDetailsLength)
	}

	if len(desc.ExternalUrl) > MaxUriLength {
		return sdkerrors.Wrapf(ErrInvalidClassDescription, "invalid external_url length; got: %d, max: %d", len(desc.ExternalUrl), MaxUriLength)
	}

	if len(desc.ImageUri) > MaxUriLength {
		return sdkerrors.Wrapf(ErrInvalidClassDescription, "invalid image_uri length; got: %d, max: %d", len(desc.ImageUri), MaxUriLength)
	}

	return nil
}

func NewNFT(id uint64, brandID, classID, uri, varUri string) NFT {
	return NFT{
		BrandId: brandID,
		Id:      id,
		ClassId: classID,
		Uri:     uri,
		VarUri:  varUri,
	}
}

func (desc *NFT) TrimSpace() NFT {
	desc.BrandId = strings.TrimSpace(desc.BrandId)
	desc.ClassId = strings.TrimSpace(desc.ClassId)
	desc.Uri = strings.TrimSpace(desc.Uri)
	desc.VarUri = strings.TrimSpace(desc.VarUri)

	return NewNFT(desc.Id, desc.BrandId, desc.ClassId, desc.Uri, desc.VarUri)
}

// need to validate url
func (nft *NFT) Validate() error {
	if nft.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidNFTID, "nft id must be greater 0")
	}

	if err := brandtypes.ValidateBrandID(nft.BrandId); err != nil {
		return err
	}

	if err := ValidateClassId(nft.ClassId); err != nil {
		return err
	}

	if len(nft.Uri) == 0 {
		return sdkerrors.Wrap(ErrInvalidNFTUri, "uri can not be blank")
	}

	if len(nft.Uri) > MaxUriLength {
		return sdkerrors.Wrapf(ErrInvalidNFTUri, "invalid uri length; got: %d, max: %d", len(nft.Uri), MaxUriLength)
	}

	if len(nft.VarUri) > MaxUriLength {
		return sdkerrors.Wrapf(ErrInvalidNFTVarUri, "invalid var_uri length; got: %d, max: %d", len(nft.VarUri), MaxUriLength)
	}

	return nil
}

func NewSupply(sequence, totalSupply uint64) Supply {
	return Supply{
		Sequence:    sequence,
		TotalSupply: totalSupply,
	}
}

func DefaultSupply() Supply {
	return NewSupply(1, 0)
}

func (supply *Supply) IncreaseSupply() {
	supply.Sequence++
	supply.TotalSupply++
}

func (supply *Supply) DecreaseSupply() {
	if supply.TotalSupply == 0 {
		return
	}
	supply.TotalSupply--
}

func ValidateClassId(id string) error {
	if !reClassID.MatchString(id) {
		return ErrInvalidClassID
	}

	return nil
}

func ValidateFeeBasisPoints(feeBasisPoints uint32) error {
	if feeBasisPoints > MaxFeeBasisPoints {
		return sdkerrors.Wrapf(ErrInvalidFeeBasisPoints, "invalid fee basis_points; got: %d, max: %d", feeBasisPoints, MaxFeeBasisPoints)
	}

	return nil
}
