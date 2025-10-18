package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"

	"govchain/x/disbursementtracker/types"
)

// Keeper orchestrates state transitions for the disbursement tracker module.
type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	procurements types.ProcurementKeeper

	Schema             collections.Schema
	DisbursementSeq    collections.Sequence
	Disbursements      collections.Map[uint64, types.Disbursement]
	DisbursementTotals collections.Map[uint64, sdkmath.Int]
}

// NewKeeper constructs the keeper with its dependencies.
func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	procurementKeeper types.ProcurementKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:       storeService,
		cdc:                cdc,
		addressCodec:       addressCodec,
		procurements:       procurementKeeper,
		Disbursements:      collections.NewMap(sb, types.DisbursementKey, "disbursement", collections.Uint64Key, collections.NewJSONValueCodec[types.Disbursement]()),
		DisbursementSeq:    collections.NewSequence(sb, types.DisbursementCountKey, "disbursement_seq"),
		DisbursementTotals: collections.NewMap(sb, types.DisbursementTotalKey, "disbursement_totals", collections.Uint64Key, collections.NewJSONValueCodec[sdkmath.Int]()),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build disbursement tracker schema: %w", err))
	}
	k.Schema = schema
	return k
}

// RegisterDisbursement records a new disbursement tied to a procurement.
func (k Keeper) RegisterDisbursement(ctx context.Context, disbursement types.Disbursement) (types.Disbursement, error) {
	if _, err := k.addressCodec.StringToBytes(disbursement.PerformedBy); err != nil {
		return types.Disbursement{}, fmt.Errorf("%w: %v", types.ErrPerformedByAddress, err)
	}
	if disbursement.Status == "" {
		disbursement.Status = types.DefaultDisbursementStatus()
	}

	procurement, err := k.procurements.GetProcurement(ctx, disbursement.ProcurementId)
	if err != nil {
		return types.Disbursement{}, fmt.Errorf("%w: %v", types.ErrProcurementMissing, err)
	}

	id, err := k.DisbursementSeq.Next(ctx)
	if err != nil {
		return types.Disbursement{}, err
	}
	disbursement.Id = id + 1

	if err := disbursement.ValidateBasic(); err != nil {
		return types.Disbursement{}, err
	}

	total, err := k.currentTotal(ctx, disbursement.ProcurementId)
	if err != nil {
		return types.Disbursement{}, err
	}
	newTotal := total.Add(disbursement.Amount)
	if newTotal.GT(procurement.Amount) {
		return types.Disbursement{}, types.ErrDisbursementAmountExceeded
	}

	if err := k.Disbursements.Set(ctx, disbursement.Id, disbursement); err != nil {
		return types.Disbursement{}, err
	}
	if err := k.DisbursementTotals.Set(ctx, disbursement.ProcurementId, newTotal); err != nil {
		return types.Disbursement{}, err
	}
	return disbursement, nil
}

func (k Keeper) currentTotal(ctx context.Context, procurementID uint64) (sdkmath.Int, error) {
	total, err := k.DisbursementTotals.Get(ctx, procurementID)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return sdkmath.ZeroInt(), nil
		}
		return sdkmath.Int{}, err
	}
	return total, nil
}

// UpdateDisbursementStatus updates the status of a disbursement entry.
func (k Keeper) UpdateDisbursementStatus(ctx context.Context, id uint64, next types.DisbursementStatus) (types.Disbursement, error) {
	disbursement, err := k.Disbursements.Get(ctx, id)
	if err != nil {
		return types.Disbursement{}, fmt.Errorf("%w: %v", types.ErrDisbursementNotFound, err)
	}
	if disbursement.Status == next {
		return disbursement, nil
	}
	if !disbursement.Status.CanTransitionTo(next) {
		return types.Disbursement{}, fmt.Errorf("%w: %s -> %s", types.ErrInvalidDisbursementTransition, disbursement.Status, next)
	}
	disbursement.Status = next
	if err := k.Disbursements.Set(ctx, disbursement.Id, disbursement); err != nil {
		return types.Disbursement{}, err
	}
	return disbursement, nil
}

// GetDisbursement retrieves a disbursement by id.
func (k Keeper) GetDisbursement(ctx context.Context, id uint64) (types.Disbursement, error) {
	disbursement, err := k.Disbursements.Get(ctx, id)
	if err != nil {
		return types.Disbursement{}, fmt.Errorf("%w: %v", types.ErrDisbursementNotFound, err)
	}
	return disbursement, nil
}

// WalkDisbursements iterates all stored disbursements.
func (k Keeper) WalkDisbursements(ctx context.Context, cb func(types.Disbursement) (bool, error)) error {
	return k.Disbursements.Walk(ctx, nil, func(_ uint64, disbursement types.Disbursement) (bool, error) {
		return cb(disbursement)
	})
}

// GetDisbursedTotal returns the cumulative disbursed amount for a procurement.
func (k Keeper) GetDisbursedTotal(ctx context.Context, procurementID uint64) (sdkmath.Int, error) {
	return k.currentTotal(ctx, procurementID)
}
