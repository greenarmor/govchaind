package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrDisbursementNotFound          = sdkerrors.Register(ModuleName, 1, "disbursement not found")
	ErrInvalidDisbursementTransition = sdkerrors.Register(ModuleName, 2, "invalid disbursement status transition")
	ErrPerformedByAddress            = sdkerrors.Register(ModuleName, 3, "invalid performed_by address")
	ErrProcurementMissing            = sdkerrors.Register(ModuleName, 4, "procurement not found for disbursement")
	ErrDisbursementAmountExceeded    = sdkerrors.Register(ModuleName, 5, "disbursement amount exceeds procurement total")
)
