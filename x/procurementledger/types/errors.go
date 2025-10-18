package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrProcurementNotFound          = sdkerrors.Register(ModuleName, 1, "procurement not found")
	ErrInvalidProcurementTransition = sdkerrors.Register(ModuleName, 2, "invalid procurement status transition")
	ErrProcurementOfficerAddr       = sdkerrors.Register(ModuleName, 3, "invalid procurement officer address")
	ErrBudgetMissingForProcurement  = sdkerrors.Register(ModuleName, 4, "budget not found for procurement")
)
