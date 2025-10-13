package types

// GenesisState captures the initial state for accountability scoring metadata.
type GenesisState struct{}

// DefaultGenesis returns a placeholder genesis state for the accountability scores module.
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs no validation until the module introduces structured data.
func (gs GenesisState) Validate() error {
	return nil
}
