package types

// GenesisState defines the initial module state for the budget registry module. The structure is
// intentionally minimal while the storage model for budget line items is designed.
type GenesisState struct{}

// DefaultGenesis returns a default genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation checks.
func (gs GenesisState) Validate() error {
	return nil
}
