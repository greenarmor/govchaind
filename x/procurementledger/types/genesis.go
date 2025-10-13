package types

// GenesisState will capture the initial procurement ledger data, including historical tenders
// and award baselines once they are modelled.
type GenesisState struct{}

// DefaultGenesis returns the placeholder genesis definition.
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs module-specific genesis validation. Validation is deferred until the module
// introduces real fields.
func (gs GenesisState) Validate() error {
	return nil
}
