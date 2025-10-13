package types

// GenesisState captures the initial state for the disbursement tracker module. Fields will be
// added to represent payment schedules and reconciliation checkpoints.
type GenesisState struct{}

// DefaultGenesis returns the placeholder genesis structure.
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs module-specific genesis validation and is a no-op for now.
func (gs GenesisState) Validate() error {
	return nil
}
