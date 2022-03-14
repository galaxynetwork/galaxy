package types

func NewGenesisState(minter Minter, params Params) GenesisState {
	return GenesisState{
		Minter: minter,
		Params: params,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Minter: DefaultInitialMinter(),
	}
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	return ValidateMinter(data.Minter)
}
