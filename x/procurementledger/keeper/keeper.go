package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"govchain/x/procurementledger/types"
)

// Keeper maintains procurement and contract lifecycle state.
type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	budgets      types.BudgetKeeper

	Schema         collections.Schema
	ProcurementSeq collections.Sequence
	Procurements   collections.Map[uint64, types.Procurement]
}

// NewKeeper constructs a keeper instance with its dependencies satisfied.
func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	budgetKeeper types.BudgetKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:   storeService,
		cdc:            cdc,
		addressCodec:   addressCodec,
		budgets:        budgetKeeper,
		Procurements:   collections.NewMap(sb, types.ProcurementKey, "procurement", collections.Uint64Key, collections.NewJSONValueCodec[types.Procurement]()),
		ProcurementSeq: collections.NewSequence(sb, types.ProcurementCountKey, "procurement_seq"),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build procurement ledger schema: %w", err))
	}
	k.Schema = schema
	return k
}

// RegisterProcurement records a new procurement entry after validating the referenced budget.
func (k Keeper) RegisterProcurement(ctx context.Context, procurement types.Procurement) (types.Procurement, error) {
	if _, err := k.addressCodec.StringToBytes(procurement.Officer); err != nil {
		return types.Procurement{}, fmt.Errorf("%w: %v", types.ErrProcurementOfficerAddr, err)
	}
	if procurement.Status == "" {
		procurement.Status = types.DefaultProcurementStatus()
	}
	if _, err := k.budgets.GetBudget(ctx, procurement.BudgetId); err != nil {
		return types.Procurement{}, fmt.Errorf("%w: %v", types.ErrBudgetMissingForProcurement, err)
	}

	id, err := k.ProcurementSeq.Next(ctx)
	if err != nil {
		return types.Procurement{}, err
	}
	procurement.Id = id + 1

	if err := procurement.ValidateBasic(); err != nil {
		return types.Procurement{}, err
	}

	if err := k.Procurements.Set(ctx, procurement.Id, procurement); err != nil {
		return types.Procurement{}, err
	}
	return procurement, nil
}

// UpdateProcurementStatus adjusts the lifecycle stage of a procurement.
func (k Keeper) UpdateProcurementStatus(ctx context.Context, id uint64, next types.ProcurementStatus) (types.Procurement, error) {
	procurement, err := k.Procurements.Get(ctx, id)
	if err != nil {
		return types.Procurement{}, fmt.Errorf("%w: %v", types.ErrProcurementNotFound, err)
	}

	if procurement.Status == next {
		return procurement, nil
	}
	if !procurement.Status.CanTransitionTo(next) {
		return types.Procurement{}, fmt.Errorf("%w: %s -> %s", types.ErrInvalidProcurementTransition, procurement.Status, next)
	}

	procurement.Status = next
	if err := k.Procurements.Set(ctx, procurement.Id, procurement); err != nil {
		return types.Procurement{}, err
	}
	return procurement, nil
}

// GetProcurement returns the procurement metadata for the supplied identifier.
func (k Keeper) GetProcurement(ctx context.Context, id uint64) (types.Procurement, error) {
	procurement, err := k.Procurements.Get(ctx, id)
	if err != nil {
		return types.Procurement{}, fmt.Errorf("%w: %v", types.ErrProcurementNotFound, err)
	}
	return procurement, nil
}

// WalkProcurements iterates all procurements.
func (k Keeper) WalkProcurements(ctx context.Context, cb func(types.Procurement) (bool, error)) error {
	return k.Procurements.Walk(ctx, nil, func(_ uint64, procurement types.Procurement) (bool, error) {
		return cb(procurement)
	})
}
