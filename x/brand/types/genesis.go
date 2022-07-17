package types

import "fmt"

// NewGenesisState returns a brand genesis state.
func NewGenesisState(brands Brands, params Params) *GenesisState {
	return &GenesisState{
		Brands: brands,
		Params: params,
	}
}

// DefaultGenesisState returns a brand default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Brands: Brands{},
		Params: DefaultParams(),
	}
}

// ValidateGenesis check the given genesis state has no integrity issues
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	seenBrands := make(map[string]bool)

	for _, brand := range data.Brands {
		if seenBrands[brand.Id] {
			return fmt.Errorf("duplicate brand for id %s", brand.Id)
		}

		if err := ValidateBrandID(brand.Id); err != nil {
			return err
		}

		seenBrands[brand.Id] = true
	}

	return nil
}
