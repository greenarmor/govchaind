package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"govchain/x/datasets/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority            []byte
	accountabilityKeeper types.AccountabilityKeeper

	Schema   collections.Schema
	Params   collections.Item[types.Params]
	EntrySeq collections.Sequence
	Entry    collections.Map[uint64, types.Entry]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

	accountabilityKeeper types.AccountabilityKeeper,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}
	if accountabilityKeeper == nil {
		panic("accountability keeper must be provided")
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:         storeService,
		cdc:                  cdc,
		addressCodec:         addressCodec,
		authority:            authority,
		accountabilityKeeper: accountabilityKeeper,

		Params:   collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Entry:    collections.NewMap(sb, types.EntryKey, "entry", collections.Uint64Key, codec.CollValue[types.Entry](cdc)),
		EntrySeq: collections.NewSequence(sb, types.EntryCountKey, "entrySequence"),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
