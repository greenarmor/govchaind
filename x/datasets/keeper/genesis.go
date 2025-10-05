package keeper

import (
	"context"

	"govchain/x/datasets/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.EntryList {
		if err := k.Entry.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.EntrySeq.Set(ctx, genState.EntryCount); err != nil {
		return err
	}
	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = k.Entry.Walk(ctx, nil, func(key uint64, elem types.Entry) (bool, error) {
		genesis.EntryList = append(genesis.EntryList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.EntryCount, err = k.EntrySeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
