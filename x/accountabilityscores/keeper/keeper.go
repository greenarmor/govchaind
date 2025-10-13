package keeper

import "govchain/x/accountabilityscores/types"

// Keeper will ultimately generate, persist, and query accountability scoring data across
// delegated stakeholders.
type Keeper struct{}

// NewKeeper creates an empty keeper instance ready for dependency injection in the app wiring.
func NewKeeper() Keeper {
	return Keeper{}
}

// DefaultGenesisState exposes the default genesis configuration for the accountability module.
func (k Keeper) DefaultGenesisState() *types.GenesisState {
	return types.DefaultGenesis()
}
