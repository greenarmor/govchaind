package types

const (
	ModuleName  = "disbursementtracker"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_disbursementtracker"
)

// KeyPrefix prepares module map prefixes for use with the Cosmos SDK KVStore APIs.
func KeyPrefix(key string) []byte {
	return []byte(key)
}
