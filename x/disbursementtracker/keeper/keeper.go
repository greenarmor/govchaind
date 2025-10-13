package keeper

import "govchain/x/disbursementtracker/types"

// Keeper orchestrates state transitions for the disbursement tracker module as the detailed
// implementation comes together.
type Keeper struct{}

// NewKeeper constructs an empty keeper placeholder.
func NewKeeper() Keeper {
	return Keeper{}
}

// DefaultGenesisState exposes the default genesis configuration for module wiring.
func (k Keeper) DefaultGenesisState() *types.GenesisState {
	return types.DefaultGenesis()
}
