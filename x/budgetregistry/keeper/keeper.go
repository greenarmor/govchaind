package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"govchain/x/budgetregistry/types"
)

// Keeper wires together module dependencies to manage budget registry state.
type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec

	Schema    collections.Schema
	BudgetSeq collections.Sequence
	Budgets   collections.Map[uint64, types.Budget]
}

// NewKeeper creates a new budget registry keeper with the required dependencies.
func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		Budgets:      collections.NewMap(sb, types.BudgetKey, "budget", collections.Uint64Key, collections.NewJSONValueCodec[types.Budget]()),
		BudgetSeq:    collections.NewSequence(sb, types.BudgetCountKey, "budget_seq"),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build budget registry schema: %w", err))
	}
	k.Schema = schema

	return k
}

// RegisterBudget stores a new budget entry and returns the persisted value.
func (k Keeper) RegisterBudget(ctx context.Context, budget types.Budget) (types.Budget, error) {
	if _, err := k.addressCodec.StringToBytes(budget.CreatedBy); err != nil {
		return types.Budget{}, fmt.Errorf("%w: %v", types.ErrBudgetCreatorInvalidAddr, err)
	}

	if budget.Status == "" {
		budget.Status = types.DefaultBudgetStatus()
	}

	id, err := k.BudgetSeq.Next(ctx)
	if err != nil {
		return types.Budget{}, err
	}
	budget.Id = id + 1

	if err := budget.ValidateBasic(); err != nil {
		return types.Budget{}, err
	}

	if err := k.Budgets.Set(ctx, budget.Id, budget); err != nil {
		return types.Budget{}, err
	}
	return budget, nil
}

// UpdateBudgetStatus changes the lifecycle status of a budget entry.
func (k Keeper) UpdateBudgetStatus(ctx context.Context, id uint64, next types.BudgetStatus) (types.Budget, error) {
	budget, err := k.Budgets.Get(ctx, id)
	if err != nil {
		return types.Budget{}, fmt.Errorf("%w: %v", types.ErrBudgetNotFound, err)
	}

	if budget.Status == next {
		return budget, nil
	}

	if !budget.Status.CanTransitionTo(next) {
		return types.Budget{}, fmt.Errorf("%w: %s -> %s", types.ErrInvalidBudgetTransition, budget.Status, next)
	}

	budget.Status = next
	if err := k.Budgets.Set(ctx, budget.Id, budget); err != nil {
		return types.Budget{}, err
	}
	return budget, nil
}

// GetBudget returns the stored budget for an identifier.
func (k Keeper) GetBudget(ctx context.Context, id uint64) (types.Budget, error) {
	budget, err := k.Budgets.Get(ctx, id)
	if err != nil {
		return types.Budget{}, fmt.Errorf("%w: %v", types.ErrBudgetNotFound, err)
	}
	return budget, nil
}

// HasBudget reports whether a budget id exists.
func (k Keeper) HasBudget(ctx context.Context, id uint64) (bool, error) {
	_, err := k.Budgets.Get(ctx, id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// WalkBudgets iterates all stored budgets.
func (k Keeper) WalkBudgets(ctx context.Context, cb func(types.Budget) (bool, error)) error {
	return k.Budgets.Walk(ctx, nil, func(_ uint64, budget types.Budget) (bool, error) {
		return cb(budget)
	})
}
