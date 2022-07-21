package types

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Classes = []Class
type Entries = []Entry

func NewGenesisState(classes Classes, entries Entries) *GenesisState {
	return &GenesisState{
		Classes: classes,
		Entries: entries,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(Classes{}, Entries{})
}

func (state *GenesisState) Validate() error {
	seenClassWithinBrand := map[string]bool{}
	for _, class := range state.Classes {
		url := strings.Join([]string{class.BrandId, class.Id}, "/")
		if seenClassWithinBrand[url] {
			return fmt.Errorf("duplicate class for id %s within the brandID %s", class.Id, class.BrandId)
		}

		if err := class.Validate(); err != nil {
			return err
		}

		seenClassWithinBrand[url] = true
	}

	seenNFT := map[string]bool{}
	for _, entry := range state.Entries {
		if _, err := sdk.AccAddressFromBech32(entry.Owner); err != nil {
			return err
		}

		for _, nft := range entry.Nfts {
			url := strings.Join([]string{nft.BrandId, nft.ClassId, strconv.FormatUint(nft.Id, 10)}, "/")
			if seenNFT[url] {
				return fmt.Errorf("duplicate nft for url(brandID/classID/id) %s", url)
			}

			if err := nft.Validate(); err != nil {
				return err
			}

			seenNFT[url] = true
		}
	}
	return nil
}
