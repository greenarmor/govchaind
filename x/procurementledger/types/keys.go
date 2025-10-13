package types

const (
	ModuleName  = "procurementledger"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_procurementledger"
)

// KeyPrefix converts a storage key name into a []byte prefix.
func KeyPrefix(key string) []byte {
	return []byte(key)
}
