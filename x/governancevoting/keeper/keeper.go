package keeper

import "govchain/x/governancevoting/types"

// Keeper will manage proposal metadata, delegation records, and cross-tier quorum enforcement.
type Keeper struct{}

// NewKeeper returns an empty keeper placeholder.
func NewKeeper() Keeper {
	return Keeper{}
}

// DefaultGenesisState exposes the default governance voting genesis definition.
func (k Keeper) DefaultGenesisState() *types.GenesisState {
	return types.DefaultGenesis()
}
