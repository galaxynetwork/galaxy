package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisState(t *testing.T) {
	tests := []struct {
		genState   *GenesisState
		expectPass bool
	}{
		{
			genState:   DefaultGenesisState(),
			expectPass: true,
		},
		{
			genState:   NewGenesisState(Brands{{Id: "brand1"}}, DefaultParams()),
			expectPass: true,
		},
		{
			genState:   NewGenesisState(Brands{{Id: "brand1"}, {Id: "brand2"}}, DefaultParams()),
			expectPass: true,
		},
		{
			genState:   NewGenesisState(Brands{{Id: ""}}, DefaultParams()),
			expectPass: false,
		},
		{
			genState:   NewGenesisState(Brands{{Id: "brand1"}, {Id: ""}}, DefaultParams()),
			expectPass: false,
		},
		{
			genState:   NewGenesisState(Brands{{Id: "brand1"}, {Id: "brand1"}}, DefaultParams()),
			expectPass: false,
		},
		{
			genState:   NewGenesisState(Brands{{Id: "brand1_"}, {Id: "brand2"}}, DefaultParams()),
			expectPass: false,
		},
	}

	for i, test := range tests {
		if test.expectPass {
			require.NoError(t, ValidateGenesis(*test.genState), "test genState index: %d, genState: %s", i, test.genState.String())
		} else {
			require.Error(t, ValidateGenesis(*test.genState), "test genState index: %d, genState: %s", i, test.genState.String())
		}
	}
}
