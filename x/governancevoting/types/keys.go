package types

const (
	ModuleName  = "governancevoting"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_governancevoting"
)

// KeyPrefix provides a helper for generating byte prefixes from string identifiers.
func KeyPrefix(key string) []byte {
	return []byte(key)
}
