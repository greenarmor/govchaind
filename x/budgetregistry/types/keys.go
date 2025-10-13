package types

const (
	ModuleName  = "budgetregistry"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_budgetregistry"
)

// KeyPrefix turns a key string into a byte prefix for persistent storage maps.
func KeyPrefix(key string) []byte {
	return []byte(key)
}
