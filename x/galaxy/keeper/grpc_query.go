package keeper

import (
	"github.com/galaxies-labs/galaxy/x/galaxy/types"
)

var _ types.QueryServer = Keeper{}
