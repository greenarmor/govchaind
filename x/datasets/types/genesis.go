package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:    DefaultParams(),
		EntryList: []Entry{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	entryIdMap := make(map[uint64]bool)
	entryCount := gs.GetEntryCount()
	for _, elem := range gs.EntryList {
		if _, ok := entryIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for entry")
		}
		if elem.Id >= entryCount {
			return fmt.Errorf("entry id should be lower or equal than the last id")
		}
		entryIdMap[elem.Id] = true
	}

	return gs.Params.Validate()
}
