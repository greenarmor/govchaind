package keeper

import "govchain/x/procurementledger/types"

// Keeper will maintain procurement and contract lifecycle state as the module is implemented.
type Keeper struct{}

// NewKeeper constructs a placeholder keeper instance.
func NewKeeper() Keeper {
	return Keeper{}
}

// DefaultGenesisState exposes the module default genesis to simplify future dependency wiring.
func (k Keeper) DefaultGenesisState() *types.GenesisState {
	return types.DefaultGenesis()
}
