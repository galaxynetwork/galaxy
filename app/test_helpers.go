package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// Setup initializes a new Galaxy App
func Setup(isCheckTx bool) *App {
	db := dbm.NewMemDB()

	enc := MakeEncodingConfig(ModuleBasics)

	app := New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,

		5,
		enc,
		simapp.EmptyAppOptions{},
	)
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(enc.Marshaler)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}
