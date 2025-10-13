package types

// GenesisState defines the baseline governance voting configuration, ready to be expanded with
// delegated quorum and accountability metadata.
type GenesisState struct{}

// DefaultGenesis returns the placeholder genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs module-specific genesis validation.
func (gs GenesisState) Validate() error {
	return nil
}
