package accountabilityscores

import "govchain/x/accountabilityscores/types"

// AppModule wires the accountability scoring logic into the Cosmos application once implemented.
type AppModule struct{}

// NewAppModule returns a placeholder app module implementation.
func NewAppModule() AppModule {
	return AppModule{}
}

// Name returns the module's name for routing and registration.
func (AppModule) Name() string {
	return types.ModuleName
}
