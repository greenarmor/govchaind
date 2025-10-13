package keeper

import "govchain/x/budgetregistry/types"

// Keeper wires together module dependencies to manage budget registry state. The struct is empty
// while the module's store layout and message handlers are defined.
type Keeper struct{}

// NewKeeper creates a new budget registry keeper placeholder.
func NewKeeper() Keeper {
	return Keeper{}
}

// DefaultGenesisState exposes the module's default genesis helper for integration with the app
// module wiring once the keeper gains concrete behaviour.
func (k Keeper) DefaultGenesisState() *types.GenesisState {
	return types.DefaultGenesis()
}
