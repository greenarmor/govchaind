package types

const (
	ModuleName  = "accountabilityscores"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_accountabilityscores"
)

// KeyPrefix standardises conversion of string keys to byte prefixes.
func KeyPrefix(key string) []byte {
	return []byte(key)
}
